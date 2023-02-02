package video

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	PublishCommentAction = 1
	DeleteCommentAction  = 2
)

type CommentVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentVideoLogic {
	return &CommentVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentVideoLogic) CommentVideo(req *types.CommentVideoRequest) (resp *types.CommentVideoReply, err error) {
	logx.WithContext(l.ctx).Infof("评论视频: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	switch req.ActionType {
	case PublishCommentAction:
		// 创建评论并获取评论数据
		res, err := l.svcCtx.VideoRpc.CommentVideo(l.ctx, &videoclient.CommentVideoRequest{
			UserId:  int32(UserId),
			VideoId: int32(req.VideoId),
			Content: req.Content,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("rpc调用失败: %s", err.Error())
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 获取评论用户信息
		userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
			Id: int32(res.UserId),
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("rpc调用失败: %s", err.Error())
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 用户是否关注改评论的作者
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
			UserId:       int32(UserId),
			FollowUserId: userInfo.Id,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("rpc调用失败: %s", err.Error())
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 封装返回数据
		return &types.CommentVideoReply{
			BasicReply: types.BasicReply(apiErr.Success),
			Comment: types.Comment{
				Id:         int(res.Id),
				Content:    res.Content,
				CreateTime: int(res.CreatedTime),
				User: types.User{
					Id:            int(userInfo.Id),
					Name:          userInfo.Name,
					FollowCount:   int(userInfo.FollowCount),
					FollowerCount: int(userInfo.FanCount),
					IsFollow:      isFollowRes.IsFollow,
				},
			},
		}, nil

	case DeleteCommentAction:
		// 权限校验，判断用户是否有权限删除此评论，目前仅支持评论作者删除此评论
		// 获取评论信息
		commentInfo, err := l.svcCtx.VideoRpc.GetCommentInfo(l.ctx, &videoclient.GetCommentInfoRequest{
			CommentId: int64(req.CommentId),
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("rpc调用失败: %s", err.Error())
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 权限校验
		if commentInfo.UserId != UserId {
			return nil, apiErr.PermissionDenied.WithDetails("此用户无权删除此评论")
		}

		// 删除评论
		if _, err = l.svcCtx.VideoRpc.DeleteVideoComment(l.ctx, &videoclient.DeleteVideoCommentRequest{
			CommentId: int64(req.CommentId),
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("rpc调用失败: %s", err.Error())
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		return &types.CommentVideoReply{
			BasicReply: types.BasicReply(apiErr.Success),
		}, nil

	default:
		// 未知的评论action
		return nil, apiErr.CommentActionUnknown
	}
}
