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

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListRequest) (resp *types.FavoriteListReply, err error) {
	// 判断用户是否登录
	UserId, err := utils.GetUserIDFormToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.UserNotLogin
	}

	l.Logger.Debugf("获取用户喜欢视频列表, 用户id:%d\n", UserId)

	// 获取用户喜欢视频列表
	res, err := l.svcCtx.VideoRpc.GetFavoriteVideoList(l.ctx, &videoclient.GetFavoriteVideoListRequest{
		UserId: int32(req.UserId),
	})

	if err != nil {
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	l.Logger.Infof("获取到的点赞视频数量为: %d\n", len(res.VideoList))

	// 封装列表数据
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	videoList := make([]types.Video, len(res.VideoList))
	for i, v := range res.VideoList {
		wg.Add(1)

		// 使用闭包将参数传入协程
		index := i
		videoRow := v

		threading.GoSafe(func() {
			defer wg.Done()

			// 获取作者信息
			// todo:作者信息可以传切片查
			authorInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &userclient.GetUserByIdRequest{
				Id: videoRow.AuthorId,
			})

			if err != nil {
				errChan <- apiErr.RPCFailed.WithDetails(err.Error())
			}

			// 获取用户是否关注该作者
			isFollowRes, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &userclient.IsFollowRequest{
				UserId:       int32(UserId),
				FollowUserId: authorInfo.Id,
			})

			if err != nil {
				errChan <- apiErr.RPCFailed.WithDetails(err.Error())
			}

			author := types.User{
				Id:            int(authorInfo.Id),
				Name:          authorInfo.Name,
				FollowCount:   int(authorInfo.FollowCount),
				FollowerCount: int(authorInfo.FanCount),
				IsFollow:      isFollowRes.IsFollow,
			}

			videoInfo := types.Video{
				Id:            int(videoRow.Id),
				Title:         videoRow.Title,
				Author:        author,
				PlayUrl:       videoRow.PlayUrl,
				CoverUrl:      videoRow.CoverUrl,
				FavoriteCount: int(videoRow.FavoriteCount),
				CommentCount:  int(videoRow.CommentCount),
				// 这里查询的是用户喜欢的视频列表,无需获取用户是否喜欢
				IsFavorite: true,
			}

			videoList[index] = videoInfo
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

	return &types.FavoriteListReply{
		Code:      apiErr.SuccessCode,
		Msg:       apiErr.Success.Msg,
		VideoList: videoList,
	}, nil
}
