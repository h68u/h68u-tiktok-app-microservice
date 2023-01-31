package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

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
	// 参数检查
	var Id int64
	if req.UserId == 0 {
		Id, err = utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
		if err != nil {
			return nil, apiErr.InvalidToken
		}
	} else {
		Id = int64(req.UserId)
	}
	if Id == int64(req.ToUserId) {
		return nil, apiErr.InvalidParams.WithDetails("不能关注自己")
	}
	if req.ActionType == 1 {
		//判断是否已经关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       int32(Id),
			FollowUserId: int32(req.ToUserId),
		})
		if isFollowReply.IsFollow {
			return nil, apiErr.AlreadyFollowed
		}
		//关注
		_, err := l.svcCtx.UserRpc.FollowUser(l.ctx, &user.FollowUserRequest{
			UserId:       int32(Id),
			FollowUserId: int32(req.ToUserId),
		})
		if err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}
	} else {
		//判断是否已经关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       int32(Id),
			FollowUserId: int32(req.ToUserId),
		})
		if !isFollowReply.IsFollow {
			return nil, apiErr.NotFollowed
		}
		//取消关注
		_, err := l.svcCtx.UserRpc.UnFollowUser(l.ctx, &user.UnFollowUserRequest{
			UserId:         int32(Id),
			UnFollowUserId: int32(req.ToUserId),
		})
		if err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}

	}
	return &types.FollowReply{
		BasicReply: types.BasicReply(apiErr.Success),
	}, nil
}
