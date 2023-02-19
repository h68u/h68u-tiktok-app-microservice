package contact

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"
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

const (
	MsgTypeReceived = 0
	MsgTypeSent     = 1
)

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListRequest) (resp *types.GetFriendListReply, err error) {
	l.Logger.Infof("获取好友列表: %v", req)

	// 拿到好友id
	friendsId, err := l.svcCtx.ContactRpc.GetFriendsList(l.ctx, &contact.GetFriendsListRequest{
		UserId: int32(req.UserId),
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取历史消息列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	l.Logger.Infof("friendsId: %v", friendsId)

	// 调用user服务的GetUserById接口
	var friendsInfo []*user.GetUserReply
	for _, v := range friendsId.FriendsId {
		info, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
			Id: v,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetUserById failed, err:%v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}
		friendsInfo = append(friendsInfo, info)
	}
	l.Logger.Infof("friendsInfo: %v", friendsInfo)

	// 拿到最近消息
	var friendsMsg []*contact.GetLatestMessageResponse
	for _, v := range friendsId.FriendsId {
		msg, err := l.svcCtx.ContactRpc.GetLatestMessage(l.ctx, &contact.GetLatestMessageRequest{
			UserAId: int32(req.UserId),
			UserBId: v,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetLatestMessage failed, err:%v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}
		friendsMsg = append(friendsMsg, msg)
	}
	l.Logger.Infof("friendsMsg: %v", friendsMsg)

	// 封装
	friendList := make([]types.Friend, len(friendsId.FriendsId))
	for i := range friendList {
		friendList[i] = types.Friend{
			Id:            int(friendsInfo[i].Id),
			Name:          friendsInfo[i].Name,
			FollowCount:   int(friendsInfo[i].FollowCount),
			FollowerCount: int(friendsInfo[i].FanCount),
			IsFollow:      true,
			Message:       friendsMsg[i].Message.Content,
			MsgType: func() int {
				if friendsMsg[i].Message.FromId == int32(req.UserId) {
					return MsgTypeSent
				} else {
					return MsgTypeReceived
				}
			}(),
		}
	}

	return &types.GetFriendListReply{
		BasicReply: types.BasicReply(apiErr.Success),
		FriendList: friendList,
	}, nil

}

// GetFriendList 获取好友列表 认定好友是关注和粉丝的并集
//func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListRequest) (resp *types.GetFriendListReply, err error) {
//
//	// 声明一个set记录好友列表
//	friendSet := set.New[*user.UserInfo](func(a, b *user.UserInfo) int {
//		if a.Id == b.Id {
//			return 0
//		} else if a.Id > b.Id {
//			return 1
//		} else {
//			return -1
//		}
//	}, set.WithGoroutineSafe())
//
//	// 记录有没有关注某个用户
//	followedMap := make(map[int32]bool)
//
//	//拿到关注列表的数据
//	GetFollowListReply, err := l.svcCtx.UserRpc.GetFollowList(l.ctx, &user.GetFollowListRequest{
//		UserId: int32(req.UserId),
//	})
//	if err != nil {
//		logx.WithContext(l.ctx).Errorf("GetFollowList failed, err:%v", err)
//		return nil, apiErr.InternalError(l.ctx, err.Error())
//	}
//	for _, v := range GetFollowListReply.FollowList {
//		friendSet.Insert(v)
//		followedMap[v.Id] = true
//	}
//
//	//拿到粉丝列表的信息
//	GetFansListReply, err := l.svcCtx.UserRpc.GetFansList(l.ctx, &user.GetFansListRequest{
//		UserId: int32(req.UserId),
//	})
//	if err != nil {
//		logx.WithContext(l.ctx).Errorf("FansListLogic.FansList GetFansList err: %v", err)
//		return nil, apiErr.InternalError(l.ctx, err.Error())
//	}
//	for _, v := range GetFansListReply.FansList {
//		friendSet.Insert(v)
//	}
//
//	//遍历set，确认哪些用户的关注信息不在map中，需要查询
//	var needQueryUserIds []int32
//	for iter := friendSet.Begin(); iter.IsValid(); iter.Next() {
//		if _, ok := followedMap[iter.Value().Id]; !ok {
//			needQueryUserIds = append(needQueryUserIds, iter.Value().Id)
//		}
//	}
//
//	//查询用户关注信息
//	IsFollowV2Reply, err := l.svcCtx.UserRpc.IsFollowV2(l.ctx, &user.IsFollowV2Request{
//		FollowList: func() []*user.IsFollowRequest {
//			var followList []*user.IsFollowRequest
//			for _, v := range needQueryUserIds {
//				followList = append(followList, &user.IsFollowRequest{
//					UserId:       int32(req.UserId),
//					FollowUserId: v,
//				})
//			}
//			return followList
//		}(),
//	})
//	if err != nil {
//		logx.WithContext(l.ctx).Errorf("FansListLogic.FansList IsFollowV2 err: %v", err)
//		return nil, apiErr.InternalError(l.ctx, err.Error())
//	}
//
//	//遍历查询结果，更新map
//	for _, v := range IsFollowV2Reply.IsFollowedUserId {
//		followedMap[v] = true
//	}
//
//	// 准备返回值
//	var friends []types.Friend
//	for iter := friendSet.Begin(); iter.IsValid(); iter.Next() {
//
//		// 查询最近消息
//		GetGetLatestMessageReply, err := l.svcCtx.ContactRpc.GetLatestMessage(l.ctx, &contact.GetLatestMessageRequest{
//			UserAId: int32(req.UserId),
//			UserBId: iter.Value().Id,
//		})
//		if err != nil {
//			logx.WithContext(l.ctx).Errorf("FansListLogic.FansList GetLatestMessageV2 err: %v", err)
//			return nil, apiErr.InternalError(l.ctx, err.Error())
//		}
//
//		friends = append(friends, types.Friend{
//			Id:            int(iter.Value().Id),
//			Name:          iter.Value().Name,
//			FollowCount:   int(iter.Value().FollowCount),
//			FollowerCount: int(iter.Value().FansCount),
//			IsFollow:      followedMap[iter.Value().Id],
//			Message:       GetGetLatestMessageReply.Message.Content,
//			MsgType: func() int {
//				if GetGetLatestMessageReply.Message.FromId == int32(req.UserId) {
//					return MsgTypeSent
//				} else {
//					return MsgTypeReceived
//				}
//			}(),
//		})
//	}
//
//	return &types.GetFriendListReply{
//		BasicReply: types.BasicReply(apiErr.Success),
//		FriendList: friends,
//	}, nil
//}
