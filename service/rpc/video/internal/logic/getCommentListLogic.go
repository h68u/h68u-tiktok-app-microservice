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

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *video.GetCommentListRequest) (*video.GetCommentListResponse, error) {
	// 按倒序获取视频的评论列表
	var comments []model.Comment
	if err := l.svcCtx.DBList.Mysql.
		Where("video_id = ?", in.VideoId).
		Limit(model.PopularVideoStandard).
		Order("created_at").
		Find(&comments).Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	// 封装评论数据
	commentList := make([]*video.Comment, 0, len(comments))
	for _, v := range comments {
		commentList = append(commentList, &video.Comment{
			Id:         int64(v.ID),
			AuthorId:   v.UserId,
			CreateTime: v.CreatedAt.Unix(),
			Content:    v.Content,
		})
	}

	return &video.GetCommentListResponse{
		CommentList: commentList,
	}, nil
}
