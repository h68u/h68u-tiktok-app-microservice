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
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		// 先查询用户是否点赞过该视频
		f := model.Favorite{}
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? And video_id = ?", in.UserId, in.VideoId).
			First(&f).Error

		// 点赞记录已存在
		if err == nil {
			return nil
		}

		// 数据库查询错误
		if err != gorm.ErrRecordNotFound {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		// 未点赞，创建记录
		f.VideoId = in.VideoId
		f.UserId = in.UserId
		if err := tx.Create(&f).Error; err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		// 视频点赞量加一
		// 如果是热门视频（在缓存中），就只更新缓存，交给定时任务更新数据库
		result, err := l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.UserId)).Result()
		if result == 1 {
			task, err := mq.NewAddCacheValueTask(utils.GenVideoInfoCacheKey(in.VideoId), "FavoriteCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return err
			}
		} else {
			if err != nil {
				l.Logger.Error(rpcErr.CacheError.Code, err.Error())
			}
			if err := tx.Model(&model.Video{}).
				Where("id = ?", in.VideoId).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).
				Error; err != nil {
				return status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}
		return nil
	})

	return &video.Empty{}, err
}
