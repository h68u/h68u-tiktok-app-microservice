package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/services/1_user/model"

	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

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
	var fanlist []*user.UserInfo
	for _, fan := range fans.Fans {
		fanlist = append(fanlist, &user.UserInfo{
			Id:          int32(fan.ID),
			Name:        fan.Username,
			FollowCount: int32(fan.FollowCount),
			FansCount:   int32(fan.FanCount),
		})
	}
	return &user.GetFansListReply{
		FansList: fanlist,
	}, nil
}
