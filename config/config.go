package config

import (
	"gameapp/adapter/redisadapter"
	"gameapp/repository/mysql"
	"gameapp/service/authservice"
	"gameapp/service/precenseservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}
type Config struct {
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	PresenceService precenseservice.Config `koanf:"presence_service"`
	Redis           redisadapter.Config    `koanf:"redis"`
}
