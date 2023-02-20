package logic

import (
	"context"
	"gorm.io/gorm/clause"
	"h68u-tiktok-app-microservice/common/model"

	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserRequest) (*user.Empty, error) {

	// 开启事务
	tx := l.svcCtx.DBList.Mysql.Begin()

	var newUser *model.User
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.Id).First(&newUser).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//newUser.Username = in.Name
	//newUser.Password = in.Password
	newUser.FollowCount = in.FollowCount
	newUser.FanCount = in.FanCount

	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&newUser).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &user.Empty{}, nil
}
