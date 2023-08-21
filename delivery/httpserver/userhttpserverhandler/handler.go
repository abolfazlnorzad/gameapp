package userhttpserverhandler

import (
	"gameapp/service/authservice"
	"gameapp/service/precenseservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authCfg       authservice.Config
	presenceSvc   precenseservice.Service
}

func New(authSvc authservice.Service, userSvc userservice.Service, validator uservalidator.Validator, authCfg authservice.Config, presenceSvc precenseservice.Service) Handler {
	return Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: validator,
		authCfg:       authCfg,
		presenceSvc:   presenceSvc,
	}
}
