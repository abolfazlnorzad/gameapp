package main

import (
	"fmt"
	"gameapp/adapter/redisadapter"
	"gameapp/config"
	"gameapp/delivery/grpcserver/presenceserver"
	"gameapp/repository/mysql"
	"gameapp/repository/redis/redispresence"
	"gameapp/service/authservice"
	"gameapp/service/presenceservice"
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
	redisAdp := redisadapter.New(cfg.Redis)
	pr := redispresence.New(redisAdp)
	presenceSvc := presenceservice.New(pr, cfg.PresenceService)
	fmt.Println("befor")
	fmt.Println("presence svc", presenceSvc.GetRepo())
	server := presenceserver.New(presenceSvc)
	server.Start()
}
