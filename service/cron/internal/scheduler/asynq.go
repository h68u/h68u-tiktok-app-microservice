package scheduler

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/cron"
	"h68u-tiktok-app-microservice/service/cron/internal/svc"
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

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     l.svcCtx.Config.Redis.Address,
			Password: l.svcCtx.Config.Redis.Password},
		nil,
	)

	syncUserInfoCacheTask := cron.NewSyncUserInfoCacheTask()
	entryID, err := scheduler.Register("@every 1h", syncUserInfoCacheTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	syncVideoInfoCacheTask := cron.NewSyncVideoInfoCacheTask()
	entryID, err = scheduler.Register("@every 301s", syncVideoInfoCacheTask)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}

func (l *AsynqServer) Stop() {
	fmt.Println("AsynqTask stop")
}
