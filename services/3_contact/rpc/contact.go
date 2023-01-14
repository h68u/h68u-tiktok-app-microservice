package main

import (
	"flag"
	"fmt"

	"h68u-tiktok-app-microservice/services/3_contact/rpc/internal/config"
	"h68u-tiktok-app-microservice/services/3_contact/rpc/internal/server"
	"h68u-tiktok-app-microservice/services/3_contact/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/3_contact/rpc/types/contact"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/contact.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		contact.RegisterContactServer(grpcServer, server.NewContactServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
