package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
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
	// 查询缓存是否存在
	if l.svcCtx.DBList.Redis.
		Exists(l.ctx, utils.GenFavoriteVideoCacheKey(in.UserId, in.VideoId)).
		Val() == 1 {
		return &video.IsFavoriteVideoResponse{
			IsFavorite: true,
		}, nil
	}

	// 查询记录是否存在
	err := l.svcCtx.DBList.Mysql.
		Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).
		First(&model.Favorite{}).Error

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

	// 记录存在，设置缓存
	err = l.svcCtx.DBList.Redis.
		Set(l.ctx, utils.GenFavoriteVideoCacheKey(in.UserId, in.VideoId), 1, 0).Err()
	if err != nil {
		return nil, status.Error(rpcErr.CacheError.Code, err.Error())
	}

	return &video.IsFavoriteVideoResponse{
		IsFavorite: true,
	}, nil
}
