package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	// 启动配置
	rest.RestConf

	// RPC 配置
	VideoRpc zrpc.RpcClientConf
	UserRpc  zrpc.RpcClientConf

	// JWT 配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
