package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
)

func (l *AsynqServer) loseFriendsHandler(ctx context.Context, t *asynq.Task) error {
	var p mq.LoseFriendsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		l.Logger.Errorf("json.Unmarshal failed: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := l.svcCtx.ContactRpc.LoseFriends(ctx, &contact.LoseFriendsRequest{
		UserAId: p.UserAId,
		UserBId: p.UserBId,
	})
	if err != nil {
		fmt.Printf("解除好友关系失败: %v\n", err)
		l.Logger.Errorf("解除好友关系失败: %v", err)
		return err
	}

	return nil
}
