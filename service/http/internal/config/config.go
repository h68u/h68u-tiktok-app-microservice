package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"h68u-tiktok-app-microservice/common/oss"
)

type Config struct {
	// 启动配置
	rest.RestConf

	// RPC 配置
	VideoRpc   zrpc.RpcClientConf
	UserRpc    zrpc.RpcClientConf
	ContactRpc zrpc.RpcClientConf

	// JWT 配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	OSS oss.AliyunCfg
}
