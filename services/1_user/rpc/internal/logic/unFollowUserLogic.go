package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/services/1_user/model"

	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowUserLogic {
	return &UnFollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFollowUserLogic) UnFollowUser(in *user.UnFollowUserRequest) (*user.Empty, error) {
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		//取消关注
		var users *model.User
		tx.Where("id = ?", in.UserId).First(&users)
		users.FollowCount--
		err := tx.Save(&users).Error
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		//被关注用户
		var unFollowUser *model.User
		tx.Where("id = ?", in.UnFollowUserId).First(&unFollowUser)
		unFollowUser.FanCount--
		err = tx.Save(&unFollowUser).Error
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		//解除关注关系
		err = tx.Model(users).Association("Follows").Delete(unFollowUser)
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &user.Empty{}, nil
}
