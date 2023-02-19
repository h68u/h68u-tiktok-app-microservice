package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/mr"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListRequest) (resp *types.FavoriteListReply, err error) {
	logx.WithContext(l.ctx).Infof("获取用户喜欢视频列表: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	l.Logger.Debugf("获取用户喜欢视频列表, 用户id:%d\n", UserId)

	// 获取用户喜欢视频列表
	likeVideoList, err := l.svcCtx.VideoRpc.GetFavoriteVideoList(l.ctx, &videoclient.GetFavoriteVideoListRequest{
		UserId: req.UserId,
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户喜欢视频列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	l.Logger.Infof("获取到的点赞视频数量为: %d\n", len(likeVideoList.VideoList))

	// 用于 reduce 时保持原来的顺序
	// likeVideoList 中 videoId 是唯一的, key 选择 videoId, value 是该 video 再 likeVideoList 中原始的位置
	orderMp := make(map[int]int, len(likeVideoList.VideoList))

	// mapreduce 并发处理列表请求
	videoList, err := mr.MapReduce(func(source chan<- interface{}) {
		for i, v := range likeVideoList.VideoList {
			source <- v
			orderMp[int(v.Id)] = i
		}

	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		videoInfo := item.(*video.VideoInfo)

		// 获取作者信息
		authorInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
			Id: videoInfo.AuthorId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取作者信息失败: %v", err)
			cancel(apiErr.InternalError(l.ctx, err.Error()))
			return
		}

		// 获取用户是否关注该作者
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
			UserId:       UserId,
			FollowUserId: authorInfo.Id,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否关注该作者失败: %v", err)
			cancel(apiErr.InternalError(l.ctx, err.Error()))
			return
		}

		author := types.User{
			Id:            authorInfo.Id,
			Name:          authorInfo.Name,
			FollowCount:   authorInfo.FollowCount,
			FollowerCount: authorInfo.FanCount,
			IsFollow:      isFollowRes.IsFollow,
		}

		writer.Write(types.Video{
			Id:            videoInfo.Id,
			Title:         videoInfo.Title,
			Author:        author,
			PlayUrl:       videoInfo.PlayUrl,
			CoverUrl:      videoInfo.CoverUrl,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			// 这里查询的是用户喜欢的视频列表,无需获取用户是否喜欢
			IsFavorite: true,
		})

	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		list := make([]types.Video, len(likeVideoList.VideoList))
		for item := range pipe {
			videoInfo := item.(types.Video)
			i, ok := orderMp[int(videoInfo.Id)]
			if !ok {
				cancel(apiErr.InternalError(l.ctx, "获取视频在列表中的原始位置失败"))
				return
			}

			list[i] = videoInfo
		}

		writer.Write(list)
	})

	if err != nil {
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	return &types.FavoriteListReply{
		BasicReply: types.BasicReply(apiErr.Success),
		VideoList:  videoList.([]types.Video),
	}, nil
}
