package apiErr

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest"
)

//const (
//	SuccessCode int = iota
//	PermissionDeniedCode
//	InvalidParamsCode
//	RPCFailedCode
//)
//
//const (
//	UserNotFoundCode int = iota + 1000
//	UserAlreadyExistCode
//	PasswordIncorrectCode
//	GenerateTokenFailedCode
//	TokenParseFailedCode
//	InvalidTokenCode
//	AlreadyFollowedCode
//	NotFollowedCode
//	UserNotLoginCode
//	InsufficientPermissionsCode
//)
//
//const (
//	FavouriteActionUnknownCode int = iota + 2000
//	CommentActionUnknownCode
//)
//
//const (
//	FileUploadFailedCode int = iota + 3000
//	FileIsNotVideoCode
//)
//
//var errCodeMap = map[int]string{
//	SuccessCode:          "成功",
//	PermissionDeniedCode: "权限不足",
//	InvalidParamsCode:    "参数错误",
//	RPCFailedCode:        "RPC错误",
//
//	UserNotFoundCode:        "该用户不存在",
//	UserAlreadyExistCode:    "该用户已存在",
//	PasswordIncorrectCode:   "密码错误",
//	GenerateTokenFailedCode: "生成token失败",
//	UserNotLoginCode:        "用户还未登录",
//
//	TokenParseFailedCode: "解析token失败",
//	InvalidTokenCode:     "无效的token",
//	AlreadyFollowedCode:  "已关注",
//	NotFollowedCode:      "未关注",
//
//	FavouriteActionUnknownCode: "未知的喜欢操作",
//	CommentActionUnknownCode:   "未知的评论操作",
//
//	InsufficientPermissionsCode: "权限不足",
//
//	FileUploadFailedCode: "文件上传失败",
//	FileIsNotVideoCode:   "文件不是视频",
//}
//
//var (
//	Success                 = NewApiError(SuccessCode)
//	PermissionDenied        = NewApiError(PermissionDeniedCode)
//	InvalidParams           = NewApiError(InvalidParamsCode)
//	RPCFailed               = NewApiError(RPCFailedCode)
//	UserNotFound            = NewApiError(UserNotFoundCode)
//	UserNotLogin            = NewApiError(UserNotLoginCode)
//	InsufficientPermissions = NewApiError(InsufficientPermissionsCode)
//
//	UserAlreadyExist    = NewApiError(UserAlreadyExistCode)
//	PasswordIncorrect   = NewApiError(PasswordIncorrectCode)
//	GenerateTokenFailed = NewApiError(GenerateTokenFailedCode)
//	TokenParseFailed    = NewApiError(TokenParseFailedCode)
//	InvalidToken        = NewApiError(InvalidTokenCode)
//	AlreadyFollowed     = NewApiError(AlreadyFollowedCode)
//	NotFollowed         = NewApiError(NotFollowedCode)
//
//	FavouriteActionUnknown = NewApiError(FavouriteActionUnknownCode)
//	CommentActionUnknown   = NewApiError(CommentActionUnknownCode)
//
//	FileUploadFailed = NewApiError(FileUploadFailedCode)
//	FileIsNotVideo   = NewApiError(FileIsNotVideoCode)
//)

type ApiErr struct {
	Code int
	Msg  string
}

type ErrInternal struct {
	ApiErr
	Details string
	TraceId string
}

type ErrResponse struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

func newError(code int, msg string) ApiErr {
	return ApiErr{
		Code: code,
		Msg:  msg,
	}
}

func (e ApiErr) Error() string {
	return e.Msg
}

// WithDetails 在基础错误上追加详细信息，例如：密码错误，密码长度不足6位
func (e ApiErr) WithDetails(detail string) ApiErr {
	return ApiErr{
		Code: e.Code,
		Msg:  e.Msg + ": " + detail,
	}
}

// InternalError 用于返回内部错误，例如RPC连接故障，数据库报错等，
// 该方法会自动从上下文中获取 traceId，
// debugDetail 用于调试，生成环境不会返回给客户端，
// 当 Mode 为 service.ProMode 仅返回 traceId，
// 当 Mode 为 service.DevMode 同时返回 debugDetail 和 traceId，帮助调试
func InternalError(ctx context.Context, debugDetail string) ErrInternal {
	return ErrInternal{
		ApiErr:  ServerInternal,
		Details: debugDetail,
		TraceId: trace.TraceIDFromContext(ctx),
	}
}

func (e ApiErr) Response() *ErrResponse {
	return &ErrResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

func (e ErrInternal) Response(cfg rest.RestConf) *ErrResponse {
	switch cfg.Mode {
	case service.ProMode:
		return &ErrResponse{
			Code: e.Code,
			Msg:  fmt.Sprintf("%s，请求 ID: %s", e.Msg, e.TraceId),
		}
	default:
		return &ErrResponse{
			Code: e.Code,
			Msg:  fmt.Sprintf("%s, Details: %s, TraceId: %s", e.Msg, e.Details, e.TraceId),
		}
	}
}
