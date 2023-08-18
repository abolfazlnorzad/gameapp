package backofficehandler

import (
	"gameapp/service/aclservice"
	"gameapp/service/authservice"
)

type Handler struct {
	aclSvc  aclservice.Service
	authSvc authservice.Service
	authCfg authservice.Config
}

func New(a aclservice.Service, authSvc authservice.Service, authCfg authservice.Config) Handler {
	return Handler{
		aclSvc:  a,
		authSvc: authSvc,
		authCfg: authCfg,
	}
}
