package svc

import (
	"h68u-tiktok-app-microservice/service/rpc/contact/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
