package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/logic"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/api/internal/types"
)

func CommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCommentListLogic(r.Context(), svcCtx)
		resp, err := l.CommentList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
