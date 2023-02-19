package rpcErr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DataBaseErrorCode codes.Code = iota + 1000
	CacheErrorCode
	PasswordEncryptFailedCode
	MQErrorCode
)

const (
	UserAlreadyExistCode codes.Code = iota + 2000
	UserNotExistCode
)

const (
	CommentNotExistCode codes.Code = iota + 3000
)

var errCodeMap = map[codes.Code]string{
	DataBaseErrorCode:         "数据库错误",
	CacheErrorCode:            "缓存错误",
	PasswordEncryptFailedCode: "密码加密失败",
	MQErrorCode:               "消息队列错误",

	UserAlreadyExistCode: "用户已存在",
	UserNotExistCode:     "用户不存在",

	CommentNotExistCode: "评论不存在",
}

var (
	UserAlreadyExist      = NewRpcError(UserAlreadyExistCode)
	DataBaseError         = NewRpcError(DataBaseErrorCode)
	CacheError            = NewRpcError(CacheErrorCode)
	MQError               = NewRpcError(MQErrorCode)
	PassWordEncryptFailed = NewRpcError(PasswordEncryptFailedCode)
	UserNotExist          = NewRpcError(UserNotExistCode)
	CommentNotExist       = NewRpcError(CommentNotExistCode)
)

type RpcError struct {
	Code    codes.Code `json:"code"`
	Message string     `json:"message"`
}

func NewRpcError(code codes.Code) *RpcError {
	return &RpcError{
		Code:    code,
		Message: errCodeMap[code],
	}
}

func Is(err error, target *RpcError) bool {

	if err == nil {
		return false
	}
	s, _ := status.FromError(err)

	return s.Code() == target.Code
}

func (e *RpcError) Error() string {
	return e.Message
}
