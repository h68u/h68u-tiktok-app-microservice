package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *user.IsFollowRequest) (*user.IsFollowReply, error) {
	// 查询缓存是否存在
	if l.svcCtx.DBList.Redis.
		Exists(l.ctx, utils.GenFollowUserCacheKey(in.UserId, in.FollowUserId)).
		Val() == 1 {
		return &user.IsFollowReply{
			IsFollow: true,
		}, nil
	}

	//判断你是否关注了这个人
	var aUser model.User
	l.svcCtx.DBList.Mysql.Where("id = ?", in.UserId).Preload("Follows", "id = ?", in.FollowUserId).Find(&aUser)

	if len(aUser.Follows) > 0 {
		// 记录存在，设置缓存
		err := l.svcCtx.DBList.Redis.
			Set(l.ctx, utils.GenFollowUserCacheKey(in.UserId, in.FollowUserId), 1, 0).Err()
		if err != nil {
			return nil, status.Error(rpcErr.CacheError.Code, err.Error())
		}
		return &user.IsFollowReply{
			IsFollow: true,
		}, nil
	} else {
		return &user.IsFollowReply{
			IsFollow: false,
		}, nil
	}

}
