package cron

import (
	"github.com/hibiken/asynq"
)

const TypeSyncVideoInfoCache = "cache:videoInfo:sync"

func NewSyncVideoInfoCacheTask() *asynq.Task {
	return asynq.NewTask(TypeSyncVideoInfoCache, nil)
}
