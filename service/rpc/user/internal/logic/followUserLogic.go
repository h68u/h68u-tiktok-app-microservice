package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowUserLogic) FollowUser(in *user.FollowUserRequest) (*user.Empty, error) {
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		//关注用户
		var users *model.User
		tx.Where("id = ?", in.UserId).First(&users)
		users.FollowCount++
		err := tx.Save(&users).Error
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		//被关注用户
		var followUser *model.User
		tx.Where("id = ?", in.FollowUserId).First(&followUser)
		followUser.FanCount++
		err = tx.Save(&followUser).Error
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		//建立关注关系
		err = tx.Model(users).Association("Follows").Append(followUser)
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
