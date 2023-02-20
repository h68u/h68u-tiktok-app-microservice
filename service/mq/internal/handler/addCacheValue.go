package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"h68u-tiktok-app-microservice/common/mq"
)

func (l *AsynqServer) addCacheValueHandler(ctx context.Context, t *asynq.Task) error {
	var p mq.AddCacheValuePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		l.Logger.Errorf("json.Unmarshal failed: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	err := l.svcCtx.Redis.HIncrBy(ctx, p.Key, p.Field, p.Add).Err()
	if err != nil {
		l.Logger.Errorf("redis.HIncrBy failed: %v", err)
		return err
	}
	return nil
}
