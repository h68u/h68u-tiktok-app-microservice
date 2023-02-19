package handler

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/service/mq/internal/svc"
	"log"
)

type AsynqServer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAsynqServer(ctx context.Context, svcCtx *svc.ServiceContext) *AsynqServer {
	return &AsynqServer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AsynqServer) Start() {
	fmt.Println("AsynqTask start")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: l.svcCtx.Config.Redis.Address, Password: l.svcCtx.Config.Redis.Password},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(mq.TypeTryMakeFriends, l.tryMakeFriendsHandler)
	mux.HandleFunc(mq.TypeLoseFriends, l.loseFriendsHandler)
	mux.HandleFunc(mq.TypeDelCache, l.delCacheHandler)
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (l *AsynqServer) Stop() {
	fmt.Println("AsynqTask stop")
}
