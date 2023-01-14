package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/2_video/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteVideoListLogic {
	return &GetFavoriteVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteVideoListLogic) GetFavoriteVideoList(in *video.GetFavoriteVideoListRequest) (*video.GetFavoriteVideoListResponse, error) {
	// todo: add your logic here and delete this line

	return &video.GetFavoriteVideoListResponse{}, nil
}
