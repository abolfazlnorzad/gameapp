package main

import (
	"context"
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/delivery/httpserver/backofficehandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"gameapp/repository/mysql"
	"gameapp/repository/mysql/mysqlacl"
	"gameapp/repository/mysql/mysqluser"
	"gameapp/scheduler"
	"gameapp/service/aclservice"
	"gameapp/service/authservice"
	"gameapp/service/matchingservice"
	"gameapp/service/userservice"
	"gameapp/validator/matchingvalidator"
	"gameapp/validator/uservalidator"
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
	}
	fmt.Println("cfg", cfg.Mysql)
	userSvc, authSvc, userV, aclSvc, matchSvc, matchV := setupServices(cfg)

	uh := userhttpserverhandler.New(authSvc, userSvc, userV, cfg.Auth)
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
	err := server.Router.Shutdown(ctxWithTimeout)
	if err != nil {
		fmt.Println("http server shutdown error", err)
	}
	fmt.Printf("gracefully shout down called. %+v \n", ddd)
	done <- true
	time.Sleep(5 * time.Second)
	wg.Wait()
	<-ctxWithTimeout.Done()

}

func setupServices(cfg config.Config) (userservice.Service, authservice.Service, uservalidator.Validator, aclservice.Service, matchingservice.Service, matchingvalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)
	r := mysql.New(cfg.Mysql)
	userRepo := mysqluser.New(r)
	userSvc := userservice.NewUserSvc(userRepo, authSvc)
	userV := uservalidator.New(userRepo)
	aclRepo := mysqlacl.New(r)
	aclSvc := aclservice.New(aclRepo)
	matchSvc := matchingservice.New()
	matchV := matchingvalidator.New()
	return userSvc, authSvc, userV, aclSvc, matchSvc, matchV

}
