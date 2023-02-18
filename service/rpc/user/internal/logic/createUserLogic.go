package logic

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *user.CreateUserRequest) (*user.CreatUserReply, error) {
	tx := l.svcCtx.DBList.Mysql.Begin()

	// 检查是否已经存在
	var count int64
	if err := tx.Model(&model.User{}).Where("username = ?", in.Name).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	if count > 0 {
		tx.Rollback()
		return nil, status.Error(rpcErr.UserAlreadyExist.Code, rpcErr.UserAlreadyExist.Message)
	}

	// 密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.PassWordEncryptFailed.Code, err.Error())
	}

	// 准备数据
	newUser := &model.User{
		Username: in.Name,
		Password: string(password),
	}

	// 插入数据
	if err := tx.Create(newUser).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &user.CreatUserReply{
		Id: int32(newUser.ID),
	}, nil
}
