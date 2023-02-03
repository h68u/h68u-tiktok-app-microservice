package middleware

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/http/internal/config"
	"net/http"
)

type AuthMiddleware struct {
	Config config.Config
}

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		Config: c,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var token string
		if token = r.URL.Query().Get("token"); token == "" {
			token = r.PostFormValue("token")
		}
		if token == "" {
			httpx.OkJson(w, apiErr.NotLogin)
			return
		}
		isTimeOut, err := utils.ValidToken(token, m.Config.Auth.AccessSecret)
		if err != nil || isTimeOut {
			httpx.OkJson(w, apiErr.InvalidToken)
			return
		}

		next(w, r)
	}
}
