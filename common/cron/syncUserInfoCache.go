package cron

import (
	"github.com/hibiken/asynq"
)

const TypeSyncUserInfoCache = "cache:userInfo:sync"

func NewSyncUserInfoCacheTask() *asynq.Task {
	return asynq.NewTask(TypeSyncUserInfoCache, nil)
}
