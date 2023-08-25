package main

import (
	"context"
	"fmt"
	"gameapp/adapter/natsmatchinguser"
	"gameapp/adapter/presence"
	"gameapp/adapter/redisadapter"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/delivery/httpserver/backofficehandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"gameapp/repository/mysql"
	"gameapp/repository/mysql/mysqlacl"
	"gameapp/repository/mysql/mysqluser"
	"gameapp/repository/redis/redismatching"
	"gameapp/repository/redis/redispresence"
	"gameapp/scheduler"
	"gameapp/service/aclservice"
	"gameapp/service/authservice"
	"gameapp/service/matchingservice"
	"gameapp/service/presenceservice"
	"gameapp/service/userservice"
	"gameapp/validator/matchingvalidator"
	"gameapp/validator/uservalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	JwtSignKey                 = "secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	//cfg := config.Load()

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 7777},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "root",
			Password: "password",
			Port:     3306,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
		PresenceService: presenceservice.Config{
			ExpirationTime: time.Duration(time.Hour * 1),
			Prefix:         "presence",
		},
		Redis: redisadapter.Config{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
	}
	fmt.Println("cfg", cfg)
	cl, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer cl.Close()
	userSvc, authSvc, userV, aclSvc, matchSvc, matchV, presenceSvc := setupServices(cl, cfg)

	uh := userhttpserverhandler.New(authSvc, userSvc, userV, cfg.Auth, presenceSvc)
	bh := backofficehandler.New(aclSvc, authSvc, cfg.Auth)
	mh := matchinghandler.New(matchSvc, authSvc, cfg.Auth, matchV)
	server := httpserver.New(cfg, uh, bh, mh)

	go func() {
		server.Serve()
	}()
	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(matchSvc)
		wg.Add(1)
		sch.Start(done, &wg)
	}()

	ex := make(chan os.Signal)
	signal.Notify(ex, os.Interrupt)
	ddd := <-ex
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	err = server.Router.Shutdown(ctxWithTimeout)
	if err != nil {
		fmt.Println("http server shutdown error", err)
	}
	fmt.Printf("gracefully shout down called. %+v \n", ddd)
	done <- true
	time.Sleep(5 * time.Second)
	wg.Wait()
	<-ctxWithTimeout.Done()

}

func setupServices(conn *grpc.ClientConn, cfg config.Config) (userservice.Service, authservice.Service, uservalidator.Validator, aclservice.Service, matchingservice.Service, matchingvalidator.Validator, presenceservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	r := mysql.New(cfg.Mysql)
	userRepo := mysqluser.New(r)
	userSvc := userservice.NewUserSvc(userRepo, authSvc)
	userV := uservalidator.New(userRepo)
	aclRepo := mysqlacl.New(r)
	aclSvc := aclservice.New(aclRepo)
	redisAdp := redisadapter.New(cfg.Redis)
	pr := redispresence.New(redisAdp)
	presenceSvc := presenceservice.New(pr, cfg.PresenceService)
	presenceClient := presence.New()
	matchRepo := redismatching.New(redisAdp)
	natsBroker := natsmatchinguser.New("nats://127.0.0.1:4222")
	matchSvc := matchingservice.New(cfg.Matching, matchRepo, presenceClient, natsBroker)
	matchV := matchingvalidator.New()

	return userSvc, authSvc, userV, aclSvc, matchSvc, matchV, presenceSvc

}
