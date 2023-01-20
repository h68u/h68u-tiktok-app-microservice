package logic

import (
	"context"
	"h68u-tiktok-app-microservice/services/1_user/model"
	"h68u-tiktok-app-microservice/services/1_user/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"

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
	//判断你是否关注了这个人
	var aUser model.User //
	l.svcCtx.DBList.Mysql.Where("id = ?", in.UserId).Preload("Follows", "id = ?", in.FollowUserId).Find(&aUser)

	if len(aUser.Follows) > 0 {
		return &user.IsFollowReply{
			IsFollow: true,
		}, nil
	} else {
		return &user.IsFollowReply{
			IsFollow: false,
		}, nil
	}

}
