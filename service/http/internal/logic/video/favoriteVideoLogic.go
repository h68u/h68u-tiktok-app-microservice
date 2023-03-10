package video

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	FavoriteVideoAction   = 1
	UnFavoriteVideoAction = 2
)

type FavoriteVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteVideoLogic {
	return &FavoriteVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteVideoLogic) FavoriteVideo(req *types.FavoriteVideoRequest) (resp *types.FavoriteVideoReply, err error) { // 获取登录用户id
	logx.WithContext(l.ctx).Infof("收藏视频: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	switch req.ActionType {
	case FavoriteVideoAction:
		if _, err = l.svcCtx.VideoRpc.FavoriteVideo(l.ctx, &videoclient.FavoriteVideoRequest{
			UserId:  UserId,
			VideoId: req.VideoId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("收藏视频失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

	case UnFavoriteVideoAction:
		if _, err = l.svcCtx.VideoRpc.UnFavoriteVideo(l.ctx, &videoclient.UnFavoriteVideoRequest{
			UserId:  UserId,
			VideoId: req.VideoId,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("取消收藏视频失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

	default:
		return nil, apiErr.FavouriteActionUnknown
	}

	return &types.FavoriteVideoReply{
		BasicReply: types.BasicReply(apiErr.Success),
	}, nil
}
