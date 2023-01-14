package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/logic"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/svc"
	"h68u-tiktok-app-microservice/services/3_contact/api/internal/types"
)

func SendMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSendMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendMessage(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
