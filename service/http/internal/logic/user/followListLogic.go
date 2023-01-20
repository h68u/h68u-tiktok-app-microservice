package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.FollowListRequest) (resp *types.FollowListReply, err error) {
	//拿到关注列表的数据
	GetFollowListReply, err := l.svcCtx.UserRpc.GetFollowList(l.ctx, &user.GetFollowListRequest{
		UserId: int32(req.UserId),
	})
	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	var followList []types.User
	for _, follow := range GetFollowListReply.FollowList {
		//判断关注者是否关注了你
		isFollowReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       follow.Id,
			FollowUserId: int32(req.UserId),
		})
		if err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}
		followList = append(followList, types.User{
			Id:            int(follow.Id),
			Name:          follow.Name,
			FollowCount:   int(follow.FollowCount),
			FollowerCount: int(follow.FansCount),
			IsFollow:      isFollowReply.IsFollow,
		})
	}
	return &types.FollowListReply{
		Code:  apiErr.SuccessCode,
		Msg:   apiErr.Success.Msg,
		Users: followList,
	}, nil
}