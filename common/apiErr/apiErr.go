package apiErr

const (
	SuccessCode int = iota
	ErrUnknownCode
	ErrPermissionDeniedCode
	ErrInvalidParamsCode
	ErrRPCFailedCode
)

const (
	ErrUserNotFoundCode int = iota + 1000
	ErrPasswordIncorrectCode
	ErrGenerateTokenFailedCode
	ErrRegisterFailedCode
)

const (
	ErrClassAlreadyExistCode int = iota + 2000
	ErrClassNotFoundCode
)

const (
	ErrFileAlreadyExistCode int = iota + 3000
	ErrDirNotExistCode
)

var errCodeMap = map[int]string{
	SuccessCode:             "Success",
	ErrUnknownCode:          "Unknown Error",
	ErrPermissionDeniedCode: "Permission Denied",
	ErrInvalidParamsCode:    "Invalid Params",
	ErrRPCFailedCode:        "RPC Failed",

	ErrUserNotFoundCode:        "该用户不存在",
	ErrPasswordIncorrectCode:   "密码错误",
	ErrGenerateTokenFailedCode: "生成token失败",
	ErrRegisterFailedCode:      "有用户注册失败：\n",

	ErrClassAlreadyExistCode: "班级已存在",
	ErrClassNotFoundCode:     "班级不存在",

	ErrFileAlreadyExistCode: "同名文件已存在",
	ErrDirNotExistCode:      "目录不存在",
}

var (
	Success             = NewApiError(SuccessCode)
	ErrUnknown          = NewApiError(ErrUnknownCode)
	ErrPermissionDenied = NewApiError(ErrPermissionDeniedCode)
	ErrInvalidParams    = NewApiError(ErrInvalidParamsCode)
	ErrRPCFailed        = NewApiError(ErrRPCFailedCode)

	ErrUserNotFound        = NewApiError(ErrUserNotFoundCode)
	ErrPasswordIncorrect   = NewApiError(ErrPasswordIncorrectCode)
	ErrGenerateTokenFailed = NewApiError(ErrGenerateTokenFailedCode)
	ErrRegisterFailed      = NewApiError(ErrRegisterFailedCode)

	ErrClassAlreadyExist = NewApiError(ErrClassAlreadyExistCode)
	ErrClassNotFound     = NewApiError(ErrClassNotFoundCode)

	ErrFileAlreadyExist = NewApiError(ErrFileAlreadyExistCode)
	ErrDirNotExist      = NewApiError(ErrDirNotExistCode)
)

type ApiError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewApiError(code int) *ApiError {
	return &ApiError{
		Code: code,
		Msg:  errCodeMap[code],
	}
}

func (e *ApiError) Error() string {
	return e.Msg
}

func (e *ApiError) WithDetails(detail string) *ApiError {
	return &ApiError{
		Code: e.Code,
		Msg:  e.Msg + ", Details: " + detail,
	}
}
