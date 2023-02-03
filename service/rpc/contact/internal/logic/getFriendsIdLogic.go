package logic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
	"sync"
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

const (
	MsgTypeSend    = 1
	MSgTypeReceive = 0
	MsgNone        = 2
)

func (l *GetFriendsIdLogic) GetFriendsId(in *contact.GetFriendsIdRequest) (*contact.GetFriendsIdResponse, error) {
	var friends model.User
	err := l.svcCtx.DBList.Mysql.Where("id = ?", in.Id).Preload("Follows").Find(&friends).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	var friendList []*contact.UserInfo
	var messages [2]model.Message
	// 封装评论列表数据
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)
	for _, friend := range friends.Follows {
		wg.Add(1)

		var message model.Message
		var msgType int32
		threading.GoSafe(func() {
			defer wg.Done()
			//判断哪个消息合适
			var choose int32 //判断选择哪个消息作为message
			var f1 bool      //判断是否有发消息给对面
			var f2 bool      //判断是否有接受过消息
			err := l.svcCtx.DBList.Mysql.Where("from_id = ?", in.Id).Where("to_user_id = ?", friend.ID).Last(&messages[0]).Error
			//判断记录是否存在与判断哪个成为message
			if err == gorm.ErrRecordNotFound {
				choose = 1
				f1 = false
			} else {
				choose = 0
				f2 = true
				if err != nil && err != sql.ErrNoRows {
					logx.WithContext(l.ctx).Errorf("获取接收的消息失败: %v", err)
					errChan <- apiErr.InternalError(l.ctx, err.Error())
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
					logx.WithContext(l.ctx).Errorf("获取发送的消息失败: %v", err)
					errChan <- apiErr.InternalError(l.ctx, err.Error())
				}
			}
			if f1 && f2 {
				//两者均有消息
				if messages[0].CreatedAt.Unix() > messages[1].CreatedAt.Unix() {
					message = messages[0]
					msgType = MsgTypeSend
				} else {
					message = messages[1]
					msgType = MSgTypeReceive
				}
			} else if f1 || f2 {
				message = messages[choose]
				if choose == 0 {
					msgType = MsgTypeSend
				} else {
					msgType = MSgTypeReceive
				}
			} else {
				//内容均为空
				msgType = MsgNone //如果没有发过信息,我把这个设为2
				message.Content = ""
			}
			friendList = append(friendList, &contact.UserInfo{
				Id:          int32(friend.ID),
				Name:        friend.Username,
				FollowCount: int32(friend.FollowCount),
				FansCount:   int32(friend.FanCount),
				NewMessage:  message.Content,
				MsgType:     msgType,
			})
		})

	}
	threading.GoSafe(func() {
		wg.Wait()
		close(finished)
	})
	select {
	case <-finished: // 正常退出
	case err := <-errChan: // 获取时发生错误
		return nil, err
	}

	return &contact.GetFriendsIdResponse{
		FriendsId: friendList,
	}, nil

}
