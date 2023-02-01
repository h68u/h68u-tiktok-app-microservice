package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
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
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		//创建并增加消息记录
		message := model.Message{
			FromId:   int64(in.FromId),
			ToUserId: int64(in.ToId),
			Content:  in.Content,
		}

		if err := l.svcCtx.DBList.Mysql.Create(&message).Error; err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	return &contact.Empty{}, nil
}
