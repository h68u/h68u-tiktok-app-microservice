package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
)

type GetFriendsIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsIdLogic {
	return &GetFriendsIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsIdLogic) GetFriendsId(in *contact.GetFriendsIdRequest) (*contact.GetFriendsIdResponse, error) {
	var friends model.User
	err := l.svcCtx.DBList.Mysql.Where("id = ?", in.Id).Preload("Follows").Find(&friends).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	var friendList []*contact.UserInfo
	for _, friend := range friends.Follows {
		friendList = append(friendList, &contact.UserInfo{
			Id:          int32(friend.ID),
			Name:        friend.Username,
			FollowCount: int32(friend.FollowCount),
			FansCount:   int32(friend.FanCount),
		})
	}

	return &contact.GetFriendsIdResponse{
		FriendsId: friendList,
	}, nil

}
