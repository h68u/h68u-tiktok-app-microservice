package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListRequest) (resp *types.CommentListReply, err error) {
	logx.WithContext(l.ctx).Infof("获取评论列表: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	// 获取评论数据
	commentListData, err := l.svcCtx.VideoRpc.GetCommentList(l.ctx, &videoclient.GetCommentListRequest{
		VideoId: int32(req.VideoId),
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取评论列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	// 用于 reduce 时保持原来的顺序
	// commentList 中 commentId 是唯一的, key 选择 commentId, value 是该 comment 再 commentList 中原始的位置
	orderMp := make(map[int]int, len(commentListData.CommentList))

	// mapreduce 并发处理列表请求
	commentList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, v := range commentListData.CommentList {
			source <- v
			orderMp[int(v.Id)] = i
		}

	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		comment := item.(*video.Comment)

		// 获取评论用户信息
		l.Logger.Infof("评论用户id: %d", comment.AuthorId)
		userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
			Id: comment.AuthorId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取评论用户信息失败: %v", err)
			cancel(apiErr.InternalError(l.ctx, err.Error()))
			return
		}

		// 获取用户是否关注该作者
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
			UserId:       int32(UserId),
			FollowUserId: comment.AuthorId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否关注该作者失败: %v", err)
			cancel(apiErr.InternalError(l.ctx, err.Error()))
			return
		}

		// 将包装好的数据返回
		writer.Write(types.Comment{
			Id:         int(comment.Id),
			Content:    comment.Content,
			CreateTime: int(comment.CreateTime),
			User: types.User{
				Id:            int(userInfo.Id),
				Name:          userInfo.Name,
				FollowCount:   int(userInfo.FollowCount),
				FollowerCount: int(userInfo.FanCount),
				IsFollow:      isFollowRes.IsFollow,
			},
		})

	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		list := make([]types.Comment, len(commentListData.CommentList))
		for item := range pipe {
			comment := item.(types.Comment)
			// 从 orderMp 中获取评论在列表中的原始位置，避免 mapreduce 处理后导致分页数据混乱
			i, ok := orderMp[comment.Id]
			if !ok {
				cancel(apiErr.InternalError(l.ctx, "获取评论在列表中的原始位置失败"))
				return
			}

			list[i] = comment
		}

		writer.Write(list)
	})

	if err != nil {
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	return &types.CommentListReply{
		BasicReply:  types.BasicReply(apiErr.Success),
		CommentList: commentList.([]types.Comment),
	}, nil

}
