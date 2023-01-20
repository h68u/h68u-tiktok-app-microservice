package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowListLogic) GetFollowList(in *user.GetFollowListRequest) (*user.GetFollowListReply, error) {
	var follows model.User
	err := l.svcCtx.DBList.Mysql.Where("id = ?", in.UserId).Preload("Follows").Find(&follows).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	var followList []*user.UserInfo
	for _, follow := range follows.Follows {
		followList = append(followList, &user.UserInfo{
			Id:          int32(follow.ID),
			Name:        follow.Username,
			FollowCount: int32(follow.FollowCount),
			FansCount:   int32(follow.FanCount),
		})
	}

	return &user.GetFollowListReply{
		FollowList: followList,
	}, nil
}
