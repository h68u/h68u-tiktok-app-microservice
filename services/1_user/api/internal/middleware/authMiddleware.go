package middleware

import (
	"encoding/json"
	"h68u-tiktok-app-microservice/common/apiErr"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/config"
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

		token := r.URL.Query().Get("token")

		timeOut, err := utils.ValidToken(token, m.Config.Auth.AccessSecret)
		if err != nil || timeOut {
			//  返回json
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			reply, _ := json.Marshal(struct {
				Code int    `json:"status_code"`
				Msg  string `json:"status_msg"`
			}{
				Code: apiErr.PermissionDenied.Code,
				Msg:  apiErr.PermissionDenied.Msg,
			})

			// TODO: 处理错误
			w.Write(reply)
			return
		}

		// Passthrough to next handler if need
		next(w, r)
	}
}
