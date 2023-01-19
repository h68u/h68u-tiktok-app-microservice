package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/config"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/middleware"
	"h68u-tiktok-app-microservice/services/2_video/rpc/videoclient"
)

type ServiceContext struct {
	Config   config.Config
	Auth     rest.Middleware
	VideoRpc videoclient.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Auth:     middleware.NewAuthMiddleware(c).Handle,
		VideoRpc: videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
	}
}
