package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/common/utils"
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
		err := tx.
			Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&f).Error

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
			Clauses(clause.Locking{Strength: "UPDATE"}).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).
			Error; err != nil {

			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	// 异步删除缓存
	task, err := mq.NewDelCacheTask(utils.GenFavoriteVideoCacheKey(in.UserId, in.VideoId))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
		return nil, status.Error(rpcErr.MQError.Code, err.Error())
	}
	if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
		logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
		return nil, status.Error(rpcErr.MQError.Code, err.Error())
	}

	return &video.Empty{}, nil
}
