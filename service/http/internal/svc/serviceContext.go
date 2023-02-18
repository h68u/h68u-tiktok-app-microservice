package svc

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	oss2 "h68u-tiktok-app-microservice/common/oss"
	"h68u-tiktok-app-microservice/service/http/internal/config"
	"h68u-tiktok-app-microservice/service/http/internal/middleware"
	"h68u-tiktok-app-microservice/service/rpc/contact/contactclient"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
)

type ServiceContext struct {
	Config       config.Config
	Auth         rest.Middleware
	AliyunClient *oss.Client
	AsynqClient  *asynq.Client
	VideoRpc     videoclient.Video
	UserRpc      userclient.User
	ContactRpc   contactclient.Contact
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		Auth:         middleware.NewAuthMiddleware(c).Handle,
		AliyunClient: oss2.AliyunInit(c.OSS),
		AsynqClient:  asynq.NewClient(asynq.RedisClientOpt{Addr: c.Redis.Address, Password: c.Redis.Password}),
		VideoRpc:     videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:      userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ContactRpc:   contactclient.NewContact(zrpc.MustNewClient(c.ContactRpc)),
	}
}
