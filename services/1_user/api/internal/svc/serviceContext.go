package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/config"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/middleware"
	"h68u-tiktok-app-microservice/services/1_user/rpc/userclient"
)

type ServiceContext struct {
	Config  config.Config
	Auth    rest.Middleware
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Auth:    middleware.NewAuthMiddleware().Handle,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
