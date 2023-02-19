package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/model"

	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMakeFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeFriendsLogic {
	return &MakeFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MakeFriendsLogic) MakeFriends(in *contact.MakeFriendsRequest) (*contact.Empty, error) {

	newFriendsA := model.Friend{
		UserId:   in.UserAId,
		FriendId: in.UserBId,
	}

	newFriendsB := model.Friend{
		UserId:   in.UserBId,
		FriendId: in.UserAId,
	}

	tx := l.svcCtx.DBList.Mysql.Begin()

	if err := tx.Create(&newFriendsA).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(&newFriendsB).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &contact.Empty{}, nil
}
