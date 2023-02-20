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

type CommentVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentVideoLogic {
	return &CommentVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentVideoLogic) CommentVideo(in *video.CommentVideoRequest) (*video.CommentVideoResponse, error) {
	// 创建评论记录
	comment := model.Comment{
		VideoId: in.VideoId,
		UserId:  in.UserId,
		Content: in.Content,
	}

	tx := l.svcCtx.DBList.Mysql.Begin()

	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	// 更新视频评论数
	// 如果是热门视频（在缓存中），就只更新缓存，交给定时任务更新数据库
	result, err := l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.UserId)).Result()
	if result == 1 {
		task, err := mq.NewAddCacheValueTask(utils.GenVideoInfoCacheKey(in.VideoId), "CommentCount", 1)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
			tx.Rollback()
			return nil, err
		}
		if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
			logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
			tx.Rollback()
			return nil, err
		}
	} else {
		if err != nil {
			l.Logger.Error(rpcErr.CacheError.Code, err.Error())
		}
		if err := tx.Model(&model.Video{}).
			Where("id = ?", in.VideoId).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
			Error; err != nil {
			return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
	}
	tx.Commit()

	return &video.CommentVideoResponse{
		Id:          int64(comment.ID),
		UserId:      comment.UserId,
		Content:     comment.Content,
		CreatedTime: comment.CreatedAt.Unix(),
	}, nil
}
