package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *user.IsFollowRequest) (*user.IsFollowReply, error) {
	// todo: add your logic here and delete this line

	return &user.IsFollowReply{}, nil
}
