package config

var defaultCfg = map[string]any{
	"auth.refresh_subject": RefreshTokenSubject,
	"auth.access_subject":  AccessTokenSubject,
}
