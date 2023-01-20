package video

import (
	"context"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishedListLogic {
	return &PublishedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishedListLogic) PublishedList(req *types.PublishedListRequest) (resp *types.PublishedListReply, err error) {
	// todo: add your logic here and delete this line

	return
}
