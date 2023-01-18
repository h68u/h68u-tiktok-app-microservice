package apiErr

const (
	SuccessCode int = iota
	PermissionDeniedCode
	InvalidParamsCode
	RPCFailedCode
)

const (
	UserNotFoundCode int = iota + 1000
	UserAlreadyExistCode
	PasswordIncorrectCode
	GenerateTokenFailedCode
)

var errCodeMap = map[int]string{
	SuccessCode:          "成功",
	PermissionDeniedCode: "权限不足",
	InvalidParamsCode:    "参数错误",
	RPCFailedCode:        "RPC错误",

	UserNotFoundCode:        "该用户不存在",
	UserAlreadyExistCode:    "该用户已存在",
	PasswordIncorrectCode:   "密码错误",
	GenerateTokenFailedCode: "生成token失败",
}

var (
	Success          = NewApiError(SuccessCode)
	PermissionDenied = NewApiError(PermissionDeniedCode)
	InvalidParams    = NewApiError(InvalidParamsCode)
	RPCFailed        = NewApiError(RPCFailedCode)

	UserNotFound        = NewApiError(UserNotFoundCode)
	UserAlreadyExist    = NewApiError(UserAlreadyExistCode)
	PasswordIncorrect   = NewApiError(PasswordIncorrectCode)
	GenerateTokenFailed = NewApiError(GenerateTokenFailedCode)
)

type ApiError struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
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
