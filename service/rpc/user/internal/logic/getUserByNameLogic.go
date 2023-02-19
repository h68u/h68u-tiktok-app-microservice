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

type GetUserByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByNameLogic {
	return &GetUserByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByNameLogic) GetUserByName(in *user.GetUserByNameRequest) (*user.GetUserReply, error) {
	// 准备数据
	result := &model.User{}

	// 查询数据
	err := l.svcCtx.DBList.Mysql.Where("username = ?", in.Name).First(result).Error

	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(rpcErr.UserNotExist.Code, rpcErr.UserNotExist.Message)
	} else if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &user.GetUserReply{
		Id:          int64(result.ID),
		Name:        result.Username,
		Password:    result.Password,
		FollowCount: result.FollowCount,
		FanCount:    result.FanCount,
	}, nil
}
