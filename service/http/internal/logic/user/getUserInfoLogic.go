package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

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
	logx.WithContext(l.ctx).Infof("获取用户信息: %v", req)

	//从token获取自己的id
	id, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	//获取用户信息(名字与id)
	getUserByIdReply, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
		Id: req.UserId,
	})
	if rpcErr.Is(err, rpcErr.UserNotExist) {
		return nil, apiErr.UserNotFound
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	//判断是否关注了该用户
	isFollowReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
		UserId:       id,
		FollowUserId: getUserByIdReply.Id,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取是否关注失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	return &types.GetUserInfoReply{
		BasicReply: types.BasicReply(apiErr.Success),
		User: types.User{
			Id:            getUserByIdReply.Id,
			Name:          getUserByIdReply.Name,
			FollowCount:   getUserByIdReply.FollowCount,
			FollowerCount: getUserByIdReply.FanCount,
			IsFollow:      isFollowReply.IsFollow,
		},
	}, nil
}
