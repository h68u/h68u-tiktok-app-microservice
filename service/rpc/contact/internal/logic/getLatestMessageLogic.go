package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/model"

	"h68u-tiktok-app-microservice/service/rpc/contact/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLatestMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestMessageLogic {
	return &GetLatestMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLatestMessageLogic) GetLatestMessage(in *contact.GetLatestMessageRequest) (*contact.GetLatestMessageResponse, error) {
	result := model.Message{}

	l.svcCtx.DBList.Mysql.
		Where("from_id = ? and to_user_id = ?", in.UserAId, in.UserBId).
		Or("from_id = ? and to_user_id = ?", in.UserAId, in.UserBId).
		Order("created_at desc").
		First(&result)

	return &contact.GetLatestMessageResponse{
		Message: &contact.Message{
			Id:         int32(result.ID),
			Content:    result.Content,
			CreateTime: int32(result.CreatedAt.Unix()),
			FromId:     int32(result.FromId),
			ToId:       int32(result.ToUserId),
		},
	}, nil
}
