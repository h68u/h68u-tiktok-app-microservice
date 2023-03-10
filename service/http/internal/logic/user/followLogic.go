package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowRequest) (resp *types.FollowReply, err error) {
	logx.WithContext(l.ctx).Infof("关注用户: %v", req)

	// 参数检查
	var Id int64

	Id, err = utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	if Id == req.ToUserId {
		return nil, apiErr.IllegalOperation.WithDetails("不能关注自己")
	}
	if req.ActionType == 1 {

		//判断是否已经关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       Id,
			FollowUserId: req.ToUserId,
		})
		if isFollowReply.IsFollow {
			logx.WithContext(l.ctx).Errorf("已经关注过了")
			return nil, apiErr.AlreadyFollowed
		}

		//关注
		_, err := l.svcCtx.UserRpc.FollowUser(l.ctx, &user.FollowUserRequest{
			UserId:       Id,
			FollowUserId: req.ToUserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("关注失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 发送异步任务：尝试加好友
		task, err := mq.NewTryMakeFriendsTask(Id, req.ToUserId)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
			logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		//
		//// 检查是否互相关注
		//isFollowReply, _ = l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
		//	UserId:       int64(req.ToUserId),
		//	FollowUserId: int64(Id),
		//})
		//
		//// 如果相互关注了就加好友
		//if isFollowReply.IsFollow {
		//	_, err := l.svcCtx.ContactRpc.MakeFriends(l.ctx, &contact.MakeFriendsRequest{
		//		UserAId: int64(Id),
		//		UserBId: int64(req.ToUserId),
		//	})
		//	if err != nil {
		//		logx.WithContext(l.ctx).Errorf("加好友失败: %v", err)
		//		return nil, apiErr.InternalError(l.ctx, err.Error())
		//	}
		//}

	} else {
		//判断是否已经关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       Id,
			FollowUserId: req.ToUserId,
		})
		if !isFollowReply.IsFollow {
			logx.WithContext(l.ctx).Errorf("还没有关注过")
			return nil, apiErr.NotFollowed
		}

		//取消关注
		_, err := l.svcCtx.UserRpc.UnFollowUser(l.ctx, &user.UnFollowUserRequest{
			UserId:         Id,
			UnFollowUserId: req.ToUserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("取消关注失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 发送异步任务：解除好友关系
		task, err := mq.NewLoseFriendsTask(Id, req.ToUserId)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
			logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		//// 解除好友关系
		//_, err = l.svcCtx.ContactRpc.LoseFriends(l.ctx, &contact.LoseFriendsRequest{
		//	UserAId: int64(Id),
		//	UserBId: int64(req.ToUserId),
		//})
		//if err != nil {
		//	logx.WithContext(l.ctx).Errorf("解除好友关系失败: %v", err)
		//	return nil, apiErr.InternalError(l.ctx, err.Error())
		//}
	}
	return &types.FollowReply{
		BasicReply: types.BasicReply(apiErr.Success),
	}, nil
}
