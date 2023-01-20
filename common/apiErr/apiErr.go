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
	TokenParseFailedCode
	InvalidTokenCode
	AlreadyFollowedCode
	NotFollowedCode
	UserNotLoginCode
)

const (
	FavouriteActionUnknownCode int = iota + 2000
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
	UserNotLoginCode:        "用户还未登录",

	FavouriteActionUnknownCode: "未知的喜欢操作",
	TokenParseFailedCode:    "解析token失败",
	InvalidTokenCode:        "无效的token",
	AlreadyFollowedCode:     "已关注",
	NotFollowedCode:         "未关注",
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
