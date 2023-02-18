package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"
)

func (l *AsynqServer) tryMakeFriendsHandler(ctx context.Context, t *asynq.Task) error {
	var p mq.TryMakeFriendsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		l.Logger.Errorf("json.Unmarshal failed: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	isAFollowBReply, err := l.svcCtx.UserRpc.IsFollow(ctx, &user.IsFollowRequest{
		UserId:       p.UserAId,
		FollowUserId: p.UserBId,
	})
	if err != nil {
		l.Logger.Errorf("查询是否关注失败: %v", err)
		return err
	}

	isBFollowAReply, err := l.svcCtx.UserRpc.IsFollow(ctx, &user.IsFollowRequest{
		UserId:       p.UserBId,
		FollowUserId: p.UserAId,
	})
	if err != nil {
		l.Logger.Errorf("查询是否关注失败: %v", err)
		return err
	}

	// 如果相互关注了就加好友
	if isAFollowBReply.IsFollow == true && isBFollowAReply.IsFollow == true {
		_, err := l.svcCtx.ContactRpc.MakeFriends(ctx, &contact.MakeFriendsRequest{
			UserAId: p.UserBId,
			UserBId: p.UserAId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("加好友失败: %v", err)
			return err
		}
	}

	return nil
}
