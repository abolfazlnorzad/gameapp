package userhttpserverhandler

import (
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authCfg       authservice.Config
}

func New(authSvc authservice.Service, userSvc userservice.Service, validator uservalidator.Validator, authCfg authservice.Config) Handler {
	return Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: validator,
		authCfg:       authCfg,
	}
}
