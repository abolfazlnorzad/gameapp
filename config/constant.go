package config

import "time"

const (
	JwtSignKey                 = "secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	EnvPrefix                  = "GAMEAPP_"
)
