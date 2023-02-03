package apiErr

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest"
)

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
