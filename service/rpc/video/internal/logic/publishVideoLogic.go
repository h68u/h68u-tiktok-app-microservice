package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishVideoLogic) PublishVideo(in *video.PublishVideoRequest) (*video.Empty, error) {
	newVideo := &model.Video{
		AuthorId: int64(in.Video.AuthorId),
		Title:    in.Video.Title,
		PlayUrl:  in.Video.PlayUrl,
		CoverUrl: in.Video.CoverUrl,
	}

	if err := l.svcCtx.DBList.Mysql.Create(newVideo).Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &video.Empty{}, nil
}
