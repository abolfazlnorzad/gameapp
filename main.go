package main

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/delivery/httpserver/backofficehandler"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"gameapp/repository/mysql"
	"gameapp/repository/mysql/mysqlacl"
	"gameapp/repository/mysql/mysqluser"
	"gameapp/service/aclservice"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
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
	cfg2 := config.Load()

	fmt.Printf("cfg222 %+v \n ", cfg2)

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

	userSvc, authSvc, userV, aclSvc := setupServices(cfg)

	uh := userhttpserverhandler.New(authSvc, userSvc, userV, cfg.Auth)
	bh := backofficehandler.New(aclSvc, authSvc, cfg.Auth)
	server := httpserver.New(cfg, uh, bh)

	server.Serve()

}

func setupServices(cfg config.Config) (userservice.Service, authservice.Service, uservalidator.Validator, aclservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	r := mysql.New(cfg.Mysql)
	userRepo := mysqluser.New(r)
	userSvc := userservice.NewUserSvc(userRepo, authSvc)
	userV := uservalidator.New(userRepo)
	aclRepo := mysqlacl.New(r)
	aclSvc := aclservice.New(aclRepo)
	return userSvc, authSvc, userV, aclSvc
}
