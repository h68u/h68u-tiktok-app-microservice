package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFansListLogic {
	return &GetFansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFansListLogic) GetFansList(in *user.GetFansListRequest) (*user.GetFansListReply, error) {
	var fans model.User
	//获取粉丝列表
	err := l.svcCtx.DBList.Mysql.Where("id = ?", in.UserId).Preload("Fans").Find(&fans).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	var fanList []*user.UserInfo
	for _, fan := range fans.Fans {
		fanList = append(fanList, &user.UserInfo{
			Id:          int64(fan.ID),
			Name:        fan.Username,
			FollowCount: fan.FollowCount,
			FansCount:   fan.FanCount,
		})
	}
	return &user.GetFansListReply{
		FansList: fanList,
	}, nil
}
