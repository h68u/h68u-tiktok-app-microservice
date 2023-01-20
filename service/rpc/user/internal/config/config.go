package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// 数据库配置
	DBList DBListConf
}

type DBListConf struct {
	Mysql MysqlConf
	Redis RedisConf
}

type MysqlConf struct {
	Address     string
	Username    string
	Password    string
	DBName      string
	TablePrefix string
}

type RedisConf struct {
	Address  string
	Password string
	//DB       int
}
