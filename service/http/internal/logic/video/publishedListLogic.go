package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
)

type PublishedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishedListLogic {
	return &PublishedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishedListLogic) PublishedList(req *types.PublishedListRequest) (resp *types.PublishedListReply, err error) {
	logx.WithContext(l.ctx).Infof("获取发布列表: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}
	//获取发布列表数据
	publishedList, err := l.svcCtx.VideoRpc.GetVideoListByAuthor(l.ctx, &videoclient.GetVideoListByAuthorRequest{
		AuthorId: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取发布列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}
	//获取用户信息
	authorInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
		Id: req.UserId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}
	//封装数据
	videoList := make([]types.Video, 0, len(publishedList.VideoList))
	if UserId == req.UserId {
		for _, video := range publishedList.VideoList {

			videoList = append(videoList, types.Video{
				Id:    video.Id,
				Title: video.Title,
				Author: types.User{
					Id:            authorInfo.Id,
					Name:          authorInfo.Name,
					FollowCount:   authorInfo.FollowCount,
					FollowerCount: authorInfo.FanCount,
					// 这里查询的是用户自己的发布列表,无需获取用户是否关注
					IsFollow: false,
				},
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				// 这里查询的是用户自己的发布列表,无需获取用户是否点赞
				IsFavorite: true,
			})
		}
	} else {
		//是否关注作者
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
			UserId:       UserId,
			FollowUserId: req.UserId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取用户是否关注失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		author := types.User{
			Id:            authorInfo.Id,
			Name:          authorInfo.Name,
			FollowCount:   authorInfo.FollowCount,
			FollowerCount: authorInfo.FanCount,
			IsFollow:      isFollowRes.IsFollow,
		}
		for _, video := range publishedList.VideoList {
			//是否点赞
			isFavoriteRes, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoclient.IsFavoriteVideoRequest{
				UserId:  UserId,
				VideoId: video.Id,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取用户是否点赞失败: %v", err)
				return nil, apiErr.InternalError(l.ctx, err.Error())
			}
			videoList = append(videoList, types.Video{
				Id:            video.Id,
				Title:         video.Title,
				Author:        author,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    isFavoriteRes.IsFavorite,
			})
		}

	}

	return &types.PublishedListReply{
		BasicReply: types.BasicReply(apiErr.Success),
		VideoList:  videoList,
	}, nil
}
