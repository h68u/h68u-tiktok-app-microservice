package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"h68u-tiktok-app-microservice/common/mq"
)

func (l *AsynqServer) delCacheHandler(ctx context.Context, t *asynq.Task) error {
	var p mq.DelCachePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		l.Logger.Errorf("json.Unmarshal failed: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := l.svcCtx.Redis.Del(ctx, p.Key).Err(); err != nil {
		l.Logger.Errorf("redis.Del failed: %v", err)
		return err
	}

	return nil
}
