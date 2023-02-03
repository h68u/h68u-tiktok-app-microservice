package logic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
	var messages [2]model.Message
	for _, friend := range friends.Follows {
		var message model.Message
		var msgType int32
		var choose int32
		var f1 bool
		var f2 bool
		err := l.svcCtx.DBList.Mysql.Where("from_id = ?", in.Id).Where("to_user_id = ?", friend.ID).Last(&messages[0]).Error
		//判断记录是否存在与判断哪个成为message
		if err == gorm.ErrRecordNotFound {
			choose = 1
			f1 = false
		} else {
			choose = 0
			f2 = true
			if err != nil && err != sql.ErrNoRows {
				return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}
		err = l.svcCtx.DBList.Mysql.Where("from_id = ?", friend.ID).Where("to_user_id = ?", in.Id).Last(&messages[1]).Error
		if err == gorm.ErrRecordNotFound {
			choose = 0
			f1 = false
		} else {
			choose = 1
			f2 = true
			if err != nil && err != sql.ErrNoRows {
				return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}
		if f1 && f2 {
			//两者均有消息
			if messages[0].CreatedAt.Unix() > messages[1].CreatedAt.Unix() {
				message = messages[0]
				msgType = 1
			} else {
				message = messages[1]
				msgType = 0
			}
		} else if f1 || f2 {
			message = messages[choose]
			msgType = 1 - choose
		} else {
			//内容均为空
			msgType = 2 //如果没有发过信息,我把这个设为2
			friendList = append(friendList, &contact.UserInfo{
				Id:          int32(friend.ID),
				Name:        friend.Username,
				FollowCount: int32(friend.FollowCount),
				FansCount:   int32(friend.FanCount),
				NewMessage:  "",
				MsgType:     msgType,
			})
			continue
		}
		friendList = append(friendList, &contact.UserInfo{
			Id:          int32(friend.ID),
			Name:        friend.Username,
			FollowCount: int32(friend.FollowCount),
			FansCount:   int32(friend.FanCount),
			NewMessage:  message.Content,
			MsgType:     msgType,
		})
	}

	return &contact.GetFriendsIdResponse{
		FriendsId: friendList,
	}, nil

}
