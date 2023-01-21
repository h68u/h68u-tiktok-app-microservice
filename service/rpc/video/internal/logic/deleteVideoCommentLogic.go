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

type DeleteVideoCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVideoCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoCommentLogic {
	return &DeleteVideoCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteVideoCommentLogic) DeleteVideoComment(in *video.DeleteVideoCommentRequest) (*video.Empty, error) {
	// 删除评论,这里添加了用户id是为了防止普通用户删除别人的评论,实现水平权限管理
	// 也可以选择查询出comment的用户id,进行判断后删除,权限不足可以返回error
	if err := l.svcCtx.DBList.Mysql.
		Where("id = ? And user_id = ?", in.CommentId, in.UserId).
		Delete(&model.Comment{}).Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	// todo 如果支持视频作者,管理员删除视频,这里需要额外判断,删除comment_id的评论

	return &video.Empty{}, nil
}
