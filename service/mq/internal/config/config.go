package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	Redis RedisConf

	VideoRpc   zrpc.RpcClientConf
	UserRpc    zrpc.RpcClientConf
	ContactRpc zrpc.RpcClientConf
}

type RedisConf struct {
	Address  string
	Password string
	//DB       int
}
