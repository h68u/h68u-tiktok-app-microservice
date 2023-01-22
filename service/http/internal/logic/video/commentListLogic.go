package video

import (
	"context"
	"github.com/zeromicro/go-zero/core/threading"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/userclient"
	"h68u-tiktok-app-microservice/service/rpc/video/videoclient"
	"sync"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListRequest) (resp *types.CommentListReply, err error) {
	// 获取登录用户id
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.UserNotLogin
	}

	// 获取评论数据
	res, err := l.svcCtx.VideoRpc.GetCommentList(l.ctx, &videoclient.GetCommentListRequest{
		VideoId: int32(req.VideoId),
	})

	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	// 封装评论列表数据
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	commentList := make([]types.Comment, len(res.CommentList))
	for i, v := range res.CommentList {
		wg.Add(1)

		// 协程函数入参(闭包)
		index := i
		commentInfo := v

		threading.GoSafe(func() {
			defer wg.Done()
			// 获取评论用户信息

			l.Logger.Infof("评论用户id: %d", commentInfo.AuthorId)
			userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
				Id: commentInfo.AuthorId,
			})

			if err != nil {
				errChan <- apiErr.RPCFailed.WithDetails(err.Error())
			}

			// 获取用户是否关注该作者
			isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
				UserId:       int32(UserId),
				FollowUserId: commentInfo.AuthorId,
			})

			if err != nil {
				errChan <- apiErr.RPCFailed.WithDetails(err.Error())
			}

			commentList[index] = types.Comment{
				Id:         int(commentInfo.Id),
				Content:    commentInfo.Content,
				CreateTime: int(commentInfo.CreateTime),
				User: types.User{
					Id:            int(userInfo.Id),
					Name:          userInfo.Name,
					FollowCount:   int(userInfo.FollowCount),
					FollowerCount: int(userInfo.FanCount),
					IsFollow:      isFollowRes.IsFollow,
				},
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

	return &types.CommentListReply{
		Code:        apiErr.SuccessCode,
		Msg:         apiErr.Success.Msg,
		CommentList: commentList,
	}, nil

}
