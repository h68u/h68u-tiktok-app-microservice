package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/common/apiErr"
	"net/http"

	"h68u-tiktok-app-microservice/services/2_video/api/internal/config"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/handler"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/video.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *apiErr.ApiError:
			return http.StatusOK, *e
		default:
			return http.StatusInternalServerError, err
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
