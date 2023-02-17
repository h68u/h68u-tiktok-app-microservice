package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/model"

	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoseFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoseFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoseFriendsLogic {
	return &LoseFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoseFriendsLogic) LoseFriends(in *contact.LoseFriendsRequest) (*contact.Empty, error) {
	friendsA := model.Friend{
		UserId:   int64(in.UserAId),
		FriendId: int64(in.UserBId),
	}

	friendsB := model.Friend{
		UserId:   int64(in.UserBId),
		FriendId: int64(in.UserAId),
	}

	tx := l.svcCtx.DBList.Mysql.Begin()

	if err := tx.Delete(&friendsA).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Delete(&friendsB).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &contact.Empty{}, nil
}
