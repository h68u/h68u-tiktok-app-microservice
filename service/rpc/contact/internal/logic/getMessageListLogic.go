package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *contact.GetMessageListRequest) (*contact.GetMessageListResponse, error) {
	var messages []model.Message
	err := l.svcCtx.DBList.Mysql.Where("from_id = ?", in.FromId).Where("to_user_id = ?", in.ToId).Find(&messages).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	var messageList []*contact.Message
	for _, message := range messages {
		messageList = append(messageList, &contact.Message{
			Id:         int32(message.FromId),
			Content:    string(message.Content),
			CreateTime: int32(message.CreatedAt.Unix()),
		})
	}
	return &contact.GetMessageListResponse{
		Messages: messageList,
	}, nil
}
