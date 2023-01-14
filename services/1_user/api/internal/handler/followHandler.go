package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/logic"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/svc"
	"h68u-tiktok-app-microservice/services/1_user/api/internal/types"
)

func FollowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FollowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.Follow(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
