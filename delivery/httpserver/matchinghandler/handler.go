package matchinghandler

import (
	"gameapp/service/authservice"
	"gameapp/service/matchingservice"
	"gameapp/validator/matchingvalidator"
)

type Handler struct {
	matchingSvc       matchingservice.Service
	authSvc           authservice.Service
	authCfg           authservice.Config
	matchingValidator matchingvalidator.Validator
}

func New(matchingSvc matchingservice.Service,
	authSvc authservice.Service,
	authCfg authservice.Config,
	matchingValidator matchingvalidator.Validator) Handler {
	return Handler{
		matchingSvc:       matchingSvc,
		authSvc:           authSvc,
		authCfg:           authCfg,
		matchingValidator: matchingValidator,
	}
}
