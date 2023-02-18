package main

import (
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/trace"
	"h68u-tiktok-app-microservice/service/mq/internal/config"
	"h68u-tiktok-app-microservice/service/mq/internal/handler"
	"h68u-tiktok-app-microservice/service/mq/internal/svc"
)

var configFile = flag.String("f", "etc/mq.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)
	prometheus.StartAgent(c.Prometheus)
	trace.StartAgent(c.Telemetry)

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	logx.DisableStat()

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	serviceGroup.Add(handler.NewAsynqServer(ctx, svcContext))
	serviceGroup.Start()
}
