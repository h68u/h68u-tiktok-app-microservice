package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"

	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginReply, err error) {

	// 调用rpc
	GetUserByNameReply, err := l.svcCtx.UserRpc.GetUserByName(l.ctx, &user.GetUserByNameRequest{
		Name: req.Username,
	})
	if rpcErr.Is(err, rpcErr.UserNotExist) {
		return nil, apiErr.UserNotFound
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("LoginLogic.Login GetUserByName err: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(GetUserByNameReply.Password), []byte(req.Password))
	if err != nil {
		return nil, apiErr.PasswordIncorrect
	}

	// 生成 token
	jwtToken, err := utils.CreateToken(
		GetUserByNameReply.Id,
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)

	if err != nil {
		logx.WithContext(l.ctx).Errorf("LoginLogic.Login CreateToken err: %v", err)
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.LoginReply{
		BasicReply: types.BasicReply(apiErr.Success),
		UserId:     GetUserByNameReply.Id,
		Token:      jwtToken,
	}, nil
}
