package user

import (
	"context"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterReply, err error) {
	logx.WithContext(l.ctx).Infof("注册: %v", req)

	// 参数检查
	if len(req.Username) > 32 {
		return nil, apiErr.InvalidParams.WithDetails("用户名最长32个字符")
	} else if len(req.Password) > 32 {
		return nil, apiErr.InvalidParams.WithDetails("密码最长32个字符")
	}

	// 调用rpc
	CreateUserReply, err := l.svcCtx.UserRpc.CreateUser(l.ctx, &user.CreateUserRequest{
		Name:     req.Username,
		Password: req.Password,
	})

	if rpcErr.Is(err, rpcErr.UserAlreadyExist) {
		return nil, apiErr.UserAlreadyExist
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("调用rpc CreateUser 失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	// 生成 token
	jwtToken, err := utils.CreateToken(
		int64(CreateUserReply.Id),
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)

	if err != nil {
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.RegisterReply{
		BasicReply: types.BasicReply(apiErr.Success),
		UserId:     int(CreateUserReply.Id),
		Token:      jwtToken,
	}, nil
}
