package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

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
	// todo: add your logic here and delete this line

	return &user.GetFollowListReply{}, nil
}
