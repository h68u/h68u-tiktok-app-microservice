package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/services/2_video/model"

	"h68u-tiktok-app-microservice/services/2_video/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteVideoLogic {
	return &IsFavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFavoriteVideoLogic) IsFavoriteVideo(in *video.IsFavoriteVideoRequest) (*video.IsFavoriteVideoResponse, error) {
	// 查询记录是否存在
	err := l.svcCtx.DBList.Mysql.Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).First(&model.Favorite{}).Error

	// 记录不存在
	if err == gorm.ErrRecordNotFound {
		return &video.IsFavoriteVideoResponse{
			IsFavorite: false,
		}, nil
	}

	// 数据库查询错误
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &video.IsFavoriteVideoResponse{
		IsFavorite: true,
	}, nil
}
