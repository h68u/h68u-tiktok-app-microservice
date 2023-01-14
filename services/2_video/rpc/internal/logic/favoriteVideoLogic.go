package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/2_video/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteVideoLogic {
	return &FavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteVideoLogic) FavoriteVideo(in *video.FavoriteVideoRequest) (*video.Empty, error) {
	// todo: add your logic here and delete this line

	return &video.Empty{}, nil
}
