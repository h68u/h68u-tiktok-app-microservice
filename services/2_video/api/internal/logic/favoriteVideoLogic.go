package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/services/2_video/rpc/videoclient"

	"h68u-tiktok-app-microservice/services/2_video/api/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *FavoriteVideoLogic) FavoriteVideo(req *types.FavoriteVideoRequest) (resp *types.FavoriteVideoReply, err error) {
	// 通过jwt token 获取用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.UserNotLogin
	}

	l.Logger.Debugf("用户喜欢视频, 用户id:%d\n", UserId)

	// 根据action决定点赞操作
	if req.ActionType == 1 {
		if _, err = l.svcCtx.VideoRpc.FavoriteVideo(l.ctx, &videoclient.FavoriteVideoRequest{
			UserId:  int32(UserId),
			VideoId: int32(req.VideoId),
		}); err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}

	} else if req.ActionType == 0 {
		if _, err = l.svcCtx.VideoRpc.UnFavoriteVideo(l.ctx, &videoclient.UnFavoriteVideoRequest{
			UserId:  int32(UserId),
			VideoId: int32(req.VideoId),
		}); err != nil {
			return nil, apiErr.RPCFailed.WithDetails(err.Error())
		}

	} else { // 未知的点赞操作
		return nil, apiErr.FavouriteActionUnknown
	}

	return &types.FavoriteVideoReply{
		Code: apiErr.SuccessCode,
		Msg:  apiErr.Success.Msg,
	}, nil
}
