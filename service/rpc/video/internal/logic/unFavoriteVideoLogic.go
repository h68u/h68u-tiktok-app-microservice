package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFavoriteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFavoriteVideoLogic {
	return &UnFavoriteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFavoriteVideoLogic) UnFavoriteVideo(in *video.UnFavoriteVideoRequest) (*video.Empty, error) {
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		// 查询用户喜欢记录是否存在
		f := model.Favorite{}
		err := tx.Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).First(&f).Error

		// 点赞记录不存在
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		// 删除记录
		if err := tx.Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).Delete(&f).Error; err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		// 视频点赞数减一
		if err := tx.Model(&model.Video{}).
			Where("id = ?", in.VideoId).
			//Update("favorite_count", "favorite_count - 1").
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).
			Error; err != nil {

			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		return nil
	})

	return &video.Empty{}, err
}
