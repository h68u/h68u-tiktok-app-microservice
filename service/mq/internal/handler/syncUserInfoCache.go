package handler

import (
	"context"
	"github.com/hibiken/asynq"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"
)

func (l *AsynqServer) syncUserInfoCacheHandler(ctx context.Context, t *asynq.Task) error {

	res, err := l.svcCtx.Redis.LRange(ctx, utils.GenPopUserListCacheKey(), 0, -1).Result()
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	for _, v := range res {
		userId := utils.Str2Int64(v)
		// 读取缓存
		userInfo, err := l.svcCtx.Redis.HGetAll(ctx, utils.GenUserInfoCacheKey(userId)).Result()
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
		// 更新用户信息
		_, err = l.svcCtx.UserRpc.UpdateUser(ctx, &user.UpdateUserRequest{
			Id:          userId,
			Name:        userInfo["Name"],
			Password:    userInfo["Password"],
			FollowCount: utils.Str2Int64(userInfo["FollowCount"]),
			FanCount:    utils.Str2Int64(userInfo["FanCount"]),
		})
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
