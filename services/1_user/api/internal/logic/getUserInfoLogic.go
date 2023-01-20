package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/apiErr"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/types"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (resp *types.GetUserInfoReply, err error) {
	//验证用户token
	valid, err := utils.ValidToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.TokenParseFailed
	}
	if valid {
		return nil, apiErr.InvalidToken
	}

	//从token获取自己的id
	id, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.TokenParseFailed
	}

	//获取用户信息(名字与id)
	getUserByIdReply, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
		Id: int32(req.UserId),
	})
	if rpcErr.Is(err, rpcErr.UserNotExist) {
		return nil, apiErr.UserNotFound
	}
	//判断是否关注了该用户
	isFollowReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
		UserId:       int32(id),
		FollowUserId: getUserByIdReply.Id,
	})
	//获取用户的关注与粉丝
	getFollowListReply, err := l.svcCtx.UserRpc.GetFollowList(l.ctx, &user.GetFollowListRequest{
		UserId: int32(req.UserId),
	})
	getFansListReply, err := l.svcCtx.UserRpc.GetFansList(l.ctx, &user.GetFansListRequest{
		UserId: int32(req.UserId),
	})
	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	return &types.GetUserInfoReply{
		Code: apiErr.SuccessCode,
		Msg:  apiErr.Success.Msg,
		User: types.User{
			Id:            int(getUserByIdReply.Id),
			Name:          getUserByIdReply.Name,
			FollowCount:   len(getFollowListReply.FollowList),
			FollowerCount: len(getFansListReply.FansList),
			IsFollow:      isFollowReply.IsFollow,
		},
	}, nil
}
