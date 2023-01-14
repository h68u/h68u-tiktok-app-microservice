package rpcErr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rpc 错误，参考了 https://junedayday.github.io/2021/05/07/go-tip/go-tip-3/

const (
	ErrPersonNotFoundCode codes.Code = iota + 1000
	ErrPersonAlreadyExistCode
	ErrPersonCreateFailedCode
	ErrPersonUpdateFailedCode
	ErrPersonDeleteFailedCode
)

const (
	ErrClassNotFoundCode codes.Code = iota + 2000
	ErrClassAlreadyExistCode
	ErrClassCreateFailedCode
	ErrClassUpdateFailedCode
	ErrClassDeleteFailedCode
	ErrClassGetListFailedCode
)

const (
	ErrJoinClassFailedCode codes.Code = iota + 3000
	ErrLeaveClassFailedCode
)

const (
	ErrFileAlreadyExistCode codes.Code = iota + 4000
	ErrFileCreateFailedCode
	ErrFileNotExistCode
	ErrDirNotExistCode
)

var errCodeMap = map[codes.Code]string{
	ErrPersonNotFoundCode:     "用户不存在",
	ErrPersonAlreadyExistCode: "用户已存在",
	ErrPersonCreateFailedCode: "用户创建失败",
	ErrPersonUpdateFailedCode: "用户更新失败",
	ErrPersonDeleteFailedCode: "用户删除失败",

	ErrClassNotFoundCode:      "班级不存在",
	ErrClassAlreadyExistCode:  "班级已存在",
	ErrClassCreateFailedCode:  "班级创建失败",
	ErrClassUpdateFailedCode:  "班级更新失败",
	ErrClassDeleteFailedCode:  "班级删除失败",
	ErrClassGetListFailedCode: "班级列表获取失败",

	ErrJoinClassFailedCode:  "加入班级失败",
	ErrLeaveClassFailedCode: "退出班级失败",

	ErrFileAlreadyExistCode: "同目录下已存在同名文件",
	ErrFileCreateFailedCode: "文件创建失败",
	ErrFileNotExistCode:     "文件不存在",
	ErrDirNotExistCode:      "目录不存在",
}

var (
	ErrPersonNotFound     = NewRpcError(ErrPersonNotFoundCode)
	ErrPersonAlreadyExist = NewRpcError(ErrPersonAlreadyExistCode)
	ErrPersonCreateFailed = NewRpcError(ErrPersonCreateFailedCode)
	ErrPersonUpdateFailed = NewRpcError(ErrPersonUpdateFailedCode)
	//ErrPersonDeleteFailed = NewRpcError(ErrPersonDeleteFailedCode)

	ErrClassNotFound     = NewRpcError(ErrClassNotFoundCode)
	ErrClassAlreadyExist = NewRpcError(ErrClassAlreadyExistCode)
	ErrClassCreateFailed = NewRpcError(ErrClassCreateFailedCode)
	ErrClassUpdateFailed = NewRpcError(ErrClassUpdateFailedCode)
	//ErrClassDeleteFailed  = NewRpcError(ErrClassDeleteFailedCode)
	//ErrClassGetListFailed = NewRpcError(ErrClassGetListFailedCode)

	ErrJoinClassFailed  = NewRpcError(ErrJoinClassFailedCode)
	ErrLeaveClassFailed = NewRpcError(ErrLeaveClassFailedCode)

	ErrFileAlreadyExist = NewRpcError(ErrFileAlreadyExistCode)
	ErrFileCreateFailed = NewRpcError(ErrFileCreateFailedCode)
	ErrFileNotExist     = NewRpcError(ErrFileNotExistCode)
	ErrDirNotExist      = NewRpcError(ErrDirNotExistCode)
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
