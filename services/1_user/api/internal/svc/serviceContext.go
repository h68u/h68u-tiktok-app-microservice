package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/config"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/middleware"
)

type ServiceContext struct {
	Config config.Config
	Auth   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Auth:   middleware.NewAuthMiddleware().Handle,
	}
}
