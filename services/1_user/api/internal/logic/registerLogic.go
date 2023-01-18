package logic

import (
	"context"
	"h68u-tiktok-app-microservice/common/apiErr"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/services/1_user/rpc/types/user"
	"time"

	"h68u-tiktok-app-microservice/services/1_user/api/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/types"

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
		return nil, apiErr.RPCFailed.WithDetails(err.Error())
	}

	// 生成 token
	var jwtToken string
	var payload = make(map[string]interface{})
	payload["Id"] = CreateUserReply.Id
	payload["Username"] = req.Username
	payload["exp"] = time.Now().Unix() + l.svcCtx.Config.Auth.AccessExpire
	payload["iat"] = time.Now().Unix()

	if jwtToken, err = utils.NewJwtToken(payload, l.svcCtx.Config.Auth.AccessSecret); err != nil {
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.RegisterReply{
		Code:   apiErr.SuccessCode,
		Msg:    apiErr.Success.Msg,
		UserId: int(CreateUserReply.Id),
		Token:  jwtToken,
	}, nil
}
