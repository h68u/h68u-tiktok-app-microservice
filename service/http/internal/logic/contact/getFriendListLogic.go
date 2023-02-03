package contact

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListRequest) (resp *types.GetFriendListReply, err error) {
	logx.WithContext(l.ctx).Infof("获取朋友列表: %v", req)
	//拿数据
	GetFriendListResponse, err := l.svcCtx.ContactRpc.GetFriendsId(l.ctx, &contact.GetFriendsIdRequest{
		Id: int32(req.UserId),
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取朋友列表失败:%v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	var friendList []types.Friend
	for _, friend := range GetFriendListResponse.FriendsId {
		//判断关注者是否关注了你
		isFollowReply, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       friend.Id,
			FollowUserId: int32(req.UserId),
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("IsFollow failed, err:%v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}
		friendList = append(friendList, types.Friend{
			Id:            int(friend.Id),
			Name:          friend.Name,
			FollowCount:   int(friend.FollowCount),
			FollowerCount: int(friend.FansCount),
			IsFollow:      isFollowReply.IsFollow,
			NewMessage:    friend.NewMessage,
			MsgType:       int(friend.MsgType),
		})
	}
	return &types.GetFriendListReply{
		Code:       0,
		Msg:        "Success",
		FriendList: friendList,
	}, nil

}
