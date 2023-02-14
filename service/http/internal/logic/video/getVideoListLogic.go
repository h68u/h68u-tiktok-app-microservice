package video

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
	"time"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// lastTime 值0时不限制；限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
// nextTime 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time

func (l *GetVideoListLogic) GetVideoList(req *types.GetVideoListRequest) (resp *types.GetVideoListReply, err error) {
	logx.WithContext(l.ctx).Infof("GetVideoList req: %+v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		UserId = 0
	}

	// 获取视频列表
	GetVideoListReply, err := l.svcCtx.VideoRpc.GetVideoList(l.ctx, &videoclient.GetVideoListRequest{
		Num: 20,
		LatestTime: func() int64 {
			if req.LatestTime == 0 {
				return time.Now().Unix()
			} else {
				return req.LatestTime / 1000 // 前端传入的时间戳精确到毫秒，转换为秒
			}
		}(),
	})

	// 封装返回体
	resp = &types.GetVideoListReply{}
	resp.BasicReply = types.BasicReply(apiErr.Success)
	if len(GetVideoListReply.VideoList) != 0 {
		resp.NextTime = GetVideoListReply.VideoList[len(GetVideoListReply.VideoList)-1].CreateTime
	}

	for _, v := range GetVideoListReply.VideoList {

		// 获取视频作者信息
		GetUserInfoReply, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
			Id: v.AuthorId,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetUserById err: %+v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// 获取用户关注状态
		isFollowed := false
		if UserId != 0 {
			IsFollowedReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
				UserId:       int32(UserId),
				FollowUserId: v.AuthorId,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("IsFollow err: %+v", err)
				return nil, apiErr.InternalError(l.ctx, err.Error())
			}
			isFollowed = IsFollowedReply.IsFollow
		}

		// 获取视频收藏状态
		isFavorite := false
		if UserId != 0 {
			IsFavoriteVideoReply, err := l.svcCtx.VideoRpc.IsFavoriteVideo(l.ctx, &videoclient.IsFavoriteVideoRequest{
				UserId:  int32(UserId),
				VideoId: v.Id,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("IsFavoriteVideo err: %+v", err)
				return nil, apiErr.InternalError(l.ctx, err.Error())
			}
			isFavorite = IsFavoriteVideoReply.IsFavorite
		}

		// 封装返回体
		resp.VideoList = append(resp.VideoList, types.Video{
			Id:    int(v.Id),
			Title: v.Title,
			Author: types.User{
				Id:            int(v.AuthorId),
				Name:          GetUserInfoReply.Name,
				FollowCount:   int(GetUserInfoReply.FollowCount),
				FollowerCount: int(GetUserInfoReply.FanCount),
				IsFollow:      isFollowed,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int(v.FavoriteCount),
			CommentCount:  int(v.CommentCount),
			IsFavorite:    isFavorite,
		})

	}

	return resp, nil

}
