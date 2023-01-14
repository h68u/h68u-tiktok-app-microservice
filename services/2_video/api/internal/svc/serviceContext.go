package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/config"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/middleware"
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
