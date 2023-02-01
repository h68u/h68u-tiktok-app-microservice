package contact

import (
	"context"
	"github.com/zeromicro/go-zero/core/threading"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/contact/contactclient"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.UserNotLogin
	}
	res, err := l.svcCtx.ContactRpc.GetMessageList(l.ctx, &contactclient.GetMessageListRequest{
		FromId: int32(UserId),
		ToId:   int32(req.ToUserId),
	})

	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	// 封装评论列表数据
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	messageList := make([]types.Message, len(res.Messages))
	for i, c := range res.Messages {
		wg.Add(1)

		// 协程函数入参(闭包)
		index := i
		messageInfo := c

		threading.GoSafe(func() {
			defer wg.Done()
			// 获取对面用户信息

			l.Logger.Infof("消息用户id: %d", messageInfo.Id)

			if err != nil {
				errChan <- apiErr.RPCFailed.WithDetails(err.Error())
			}

			messageList[index] = types.Message{
				Id:         int(messageInfo.Id),
				Content:    messageInfo.Content,
				CreateTime: string(messageInfo.CreateTime),
			}

		})
	}

	threading.GoSafe(func() {
		wg.Wait()
		close(finished)
	})

	select {
	case <-finished: // 正常退出
	case err := <-errChan: // 获取时发生错误
		return nil, err
	}

	return &types.GetHistoryMessageReply{
		BasicReply:  types.BasicReply(apiErr.Success),
		MessageList: messageList,
	}, nil
}
