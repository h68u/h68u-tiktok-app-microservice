package logic

import (
	"context"

	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowV2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowV2Logic {
	return &IsFollowV2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowV2Logic) IsFollowV2(in *user.IsFollowV2Request) (*user.IsFollowV2Reply, error) {
	// todo: add your logic here and delete this line

	return &user.IsFollowV2Reply{}, nil
}
