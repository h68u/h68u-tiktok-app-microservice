package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"h68u-tiktok-app-microservice/service/mq/internal/config"
	"h68u-tiktok-app-microservice/service/rpc/contact/contactclient"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
)

type ServiceContext struct {
	Config     config.Config
	VideoRpc   videoclient.Video
	UserRpc    userclient.User
	ContactRpc contactclient.Contact
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		VideoRpc:   videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ContactRpc: contactclient.NewContact(zrpc.MustNewClient(c.ContactRpc)),
	}
}
