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
		VideoId: int64(in.VideoId),
		UserId:  int64(in.UserId),
		Content: in.Content,
	}

	if err := l.svcCtx.DBList.Mysql.Create(&comment).Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &video.CommentVideoResponse{
		Id:          comment.VideoId,
		UserId:      comment.UserId,
		Content:     comment.Content,
		CreatedTime: int32(comment.CreatedAt.Unix()),
	}, nil
}
