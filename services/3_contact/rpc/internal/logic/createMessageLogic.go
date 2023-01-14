package logic

import (
	"context"

	"h68u-tiktok-app-microservice/services/3_contact/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/3_contact/rpc/types/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMessageLogic) CreateMessage(in *contact.CreateMessageRequest) (*contact.Empty, error) {
	// todo: add your logic here and delete this line

	return &contact.Empty{}, nil
}
