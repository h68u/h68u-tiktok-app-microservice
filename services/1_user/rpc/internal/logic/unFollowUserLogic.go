package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowUserLogic {
	return &UnFollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFollowUserLogic) UnFollowUser(in *user.UnFollowUserRequest) (*user.Empty, error) {
	// todo: add your logic here and delete this line

	return &user.Empty{}, nil
}
