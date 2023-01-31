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
	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.UserNotLogin
	}
	//获取发布列表数据
	publishedList, err := l.svcCtx.VideoRpc.GetVideoListByAuthor(l.ctx, &videoclient.GetVideoListByAuthorRequest{
		AuthorId: int32(req.UserId),
	})
	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}
	//获取用户信息
	authorInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
		Id: int32(req.UserId),
	})
	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}
	//封装数据
	videoList := make([]types.Video, 0, len(publishedList.VideoList))
	if UserId == int64(req.UserId) {
		for _, video := range publishedList.VideoList {

			videoList = append(videoList, types.Video{
				Id:    int(video.Id),
				Title: video.Title,
				Author: types.User{
					Id:            int(authorInfo.Id),
					Name:          authorInfo.Name,
					FollowCount:   int(authorInfo.FollowCount),
					FollowerCount: int(authorInfo.FanCount),
					// 这里查询的是用户自己的发布列表,无需获取用户是否关注
					IsFollow: false,
				},
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: int(video.FavoriteCount),
				CommentCount:  int(video.CommentCount),
				// 这里查询的是用户自己的发布列表,无需获取用户是否点赞
				IsFavorite: true,
			})
		}
	} else {
		//是否关注作者
		isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
			UserId:       int32(UserId),
			FollowUserId: int32(req.UserId),
		})
		if err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}

		author := types.User{
			Id:            int(authorInfo.Id),
			Name:          authorInfo.Name,
			FollowCount:   int(authorInfo.FollowCount),
			FollowerCount: int(authorInfo.FanCount),
			IsFollow:      isFollowRes.IsFollow,
		}
		for _, video := range publishedList.VideoList {
			//是否点赞
			isFavoriteRes, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoclient.IsFavoriteVideoRequest{
				UserId:  int32(UserId),
				VideoId: video.Id,
			})
			if err != nil {
				return nil, apiErr.RPCFailed.WithDetails(err.Error())
			}
			videoList = append(videoList, types.Video{
				Id:            int(video.Id),
				Title:         video.Title,
				Author:        author,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: int(video.FavoriteCount),
				CommentCount:  int(video.CommentCount),
				IsFavorite:    isFavoriteRes.IsFavorite,
			})
		}

	}

	return &types.PublishedListReply{
		BasicReply: types.BasicReply(apiErr.Success),
		VideoList:  videoList,
	}, nil
}
