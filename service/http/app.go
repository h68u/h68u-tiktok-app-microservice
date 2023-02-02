package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"net/http"

	"h68u-tiktok-app-microservice/service/http/internal/config"
	"h68u-tiktok-app-microservice/service/http/internal/handler"
	"h68u-tiktok-app-microservice/service/http/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/app.yaml", "the config file")

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
		case apiErr.ApiError:
			return http.StatusOK, e.Response()
		case apiErr.ApiErrorInternal:
			//logx.WithContext(logx.ContextWithFields(context.Background(),
			//			//	logx.LogField{
			//			//		Key:   "TraceID",
			//			//		Value: e.TraceId,
			//			//	}),
			//			//).Error(e.Details)
			//logx.Errorw("ApiErrorInternal",
			//	logx.Field("Code", e.Code),
			//	logx.Field("Msg", e.Msg),
			//	logx.Field("TraceID", e.TraceId),
			//	logx.Field("Details", e.Details),
			//)
			return http.StatusOK, e.Response(c.RestConf)
		default:
			return http.StatusInternalServerError, err
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
