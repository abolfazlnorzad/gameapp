package main

import (
	"gameapp/adapter/redisadapter"
	"gameapp/config"
	"gameapp/delivery/grpcserver/presenceserver"
	"gameapp/repository/redis/redispresence"
	"gameapp/service/precenseservice"
)

func main() {
	cfg := config.Load()
	redisAdp := redisadapter.New(cfg.Redis)
	pr := redispresence.New(redisAdp)
	presenceSvc := precenseservice.New(pr, cfg.PresenceService)
	server := presenceserver.New(presenceSvc)
	server.Start()
}
