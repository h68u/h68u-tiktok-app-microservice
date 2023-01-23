package middleware

import (
	"fmt"
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

		token := r.URL.Query().Get("token")

		if token == "" {
			token = r.PostFormValue("token")
		}

		timeOut, err := utils.ValidToken(token, m.Config.Auth.AccessSecret)
		if err != nil || timeOut {
			//  返回json
			//w.Header().Set("Content-Type", "application/json; charset=utf-8")
			//w.WriteHeader(http.StatusOK)
			//
			//reply, _ := json.Marshal(struct {
			//	Code int    `json:"status_code"`
			//	Msg  string `json:"status_msg"`
			//}{
			//	Code: apiErr.PermissionDenied.Code,
			//	Msg:  apiErr.PermissionDenied.Msg,
			//})
			//

			//w.Write(reply)
			fmt.Println("unpass!")
			httpx.OkJson(w, apiErr.PermissionDenied)
			return
		}
		fmt.Println("pass!")
		// Passthrough to next handler if need
		next(w, r)
	}
}
