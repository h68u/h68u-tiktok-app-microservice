package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"h68u-tiktok-app-microservice/service/http/internal/config"
	"h68u-tiktok-app-microservice/service/http/internal/middleware"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
)

type ServiceContext struct {
	Config   config.Config
	Auth     rest.Middleware
	VideoRpc videoclient.Video
	UserRpc  userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Auth:     middleware.NewAuthMiddleware(c).Handle,
		VideoRpc: videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:  userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}