package contact

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/contact/types/contact"
	"strconv"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageRequest) (resp *types.SendMessageReply, err error) {
	// 参数检查
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}
	if UserId == req.ToUserId {
		return nil, apiErr.InvalidParams.WithDetails("不能发消息给自己")
	}
	//发送
	if req.ActionType == 1 {
		_, err = l.svcCtx.ContactRpc.CreateMessage(l.ctx, &contact.CreateMessageRequest{
			FromId:  UserId,
			ToId:    req.ToUserId,
			Content: req.Content,
		})
		if rpcErr.Is(err, rpcErr.UserNotExist) {
			return nil, apiErr.UserNotFound
		} else if err != nil {
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}
	} else {
		return nil, apiErr.MessageActionUnknown.WithDetails(strconv.Itoa(int(req.ActionType)))
	}
	return &types.SendMessageReply{
		BasicReply: types.BasicReply(apiErr.Success),
	}, nil

}
