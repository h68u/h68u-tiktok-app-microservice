package contact

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/contact/contactclient"
)

type GetHistoryMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryMessageLogic {
	return &GetHistoryMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryMessageLogic) GetHistoryMessage(req *types.GetHistoryMessageRequest) (resp *types.GetHistoryMessageReply, err error) {
	logx.WithContext(l.ctx).Infof("获取历史消息列表请求参数: %v", req)

	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	// 获取我的发送的消息
	MessageSent, err := l.svcCtx.ContactRpc.GetMessageList(l.ctx, &contactclient.GetMessageListRequest{
		FromId: UserId,
		ToId:   req.ToUserId,
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取历史消息列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	// 获取我的接收的消息
	MessageReceived, err := l.svcCtx.ContactRpc.GetMessageList(l.ctx, &contactclient.GetMessageListRequest{
		FromId: req.ToUserId,
		ToId:   UserId,
	})

	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取历史消息列表失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	var messageList []types.Message

	for _, message := range MessageSent.Messages {
		messageList = append(messageList, types.Message{
			Id:         message.Id,
			Content:    message.Content,
			CreateTime: message.CreateTime,
			FromUserId: message.FromId,
			ToUserId:   message.ToId,
		})
	}

	for _, message := range MessageReceived.Messages {
		messageList = append(messageList, types.Message{
			Id:         message.Id,
			Content:    message.Content,
			CreateTime: message.CreateTime,
			FromUserId: message.FromId,
			ToUserId:   message.ToId,
		})
	}

	return &types.GetHistoryMessageReply{
		BasicReply:  types.BasicReply(apiErr.Success),
		MessageList: messageList,
	}, nil
}
