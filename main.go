package main

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver"
	"gameapp/delivery/httpserver/userhttpserverhandler"
	"gameapp/repository/mysql"
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

	userSvc, authSvc, userV := setupServices(cfg)

	uh := userhttpserverhandler.New(authSvc, userSvc, userV, cfg.Auth)
	server := httpserver.New(cfg, uh)

	server.Serve()

}

func setupServices(cfg config.Config) (userservice.Service, authservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)
	r := mysql.New(cfg.Mysql)
	userSvc := userservice.NewUserSvc(r, authSvc)
	userV := uservalidator.New(r)
	return userSvc, authSvc, userV
}
