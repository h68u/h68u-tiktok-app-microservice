package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FansListLogic) FansList(req *types.FansListRequest) (resp *types.FansListReply, err error) {
	logx.WithContext(l.ctx).Infof("获取粉丝列表: %v", req)

	//拿到粉丝列表的信息
	GetFansListReply, err := l.svcCtx.UserRpc.GetFansList(l.ctx, &user.GetFansListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("FansListLogic.FansList GetFansList err: %+v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}
	var fansList []types.User
	for _, fans := range GetFansListReply.FansList {
		//先判断你是否关注你的粉丝
		isFollowReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       req.UserId,
			FollowUserId: fans.Id,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("FansListLogic.FansList IsFollow err: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}
		fansList = append(fansList, types.User{
			Id:            fans.Id,
			Name:          fans.Name,
			FollowCount:   fans.FollowCount,
			FollowerCount: fans.FansCount,
			IsFollow:      isFollowReply.IsFollow,
		})
	}

	return &types.FansListReply{
		BasicReply: types.BasicReply(apiErr.Success),
		Users:      fansList,
	}, nil
}
