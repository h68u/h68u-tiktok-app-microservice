package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/3_contact/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/3_contact/rpc/types/contact"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetFriendsIdLogic) GetFriendsId(in *contact.GetFriendsIdRequest) (*contact.GetFriendsIdResponse, error) {
	// todo: add your logic here and delete this line

	return &contact.GetFriendsIdResponse{}, nil
}
