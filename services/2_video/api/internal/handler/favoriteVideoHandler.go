package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/logic"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/types"
)

func FavoriteVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FavoriteVideoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFavoriteVideoLogic(r.Context(), svcCtx)
		resp, err := l.FavoriteVideo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
