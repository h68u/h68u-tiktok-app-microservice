package logic

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"h68u-tiktok-app-microservice/common/apiErr"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"
	"time"

	"h68u-tiktok-app-microservice/services/1_user/api/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/types"

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
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(GetUserByNameReply.Password), []byte(req.Password))
	if err != nil {
		return nil, apiErr.PasswordIncorrect
	}

	// 生成 token
	var jwtToken string
	var payload = make(map[string]interface{})
	payload["Id"] = GetUserByNameReply.Id
	payload["Username"] = req.Username
	payload["exp"] = time.Now().Unix() + l.svcCtx.Config.Auth.AccessExpire
	payload["iat"] = time.Now().Unix()

	if jwtToken, err = utils.NewJwtToken(payload, l.svcCtx.Config.Auth.AccessSecret); err != nil {
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.LoginReply{
		Code:   apiErr.SuccessCode,
		Msg:    apiErr.Success.Msg,
		UserId: int(GetUserByNameReply.Id),
		Token:  jwtToken,
	}, nil
}
