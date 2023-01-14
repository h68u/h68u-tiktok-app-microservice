package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/logic"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/svc"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/types"
)

func GetHistoryMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetHistoryMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetHistoryMessageLogic(r.Context(), svcCtx)
		resp, err := l.GetHistoryMessage(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
