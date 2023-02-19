package contact

import (
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/service/http/internal/logic/contact"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
)

func SendMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendMessageRequest
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.Error(w, err)
		//	return
		//}

		// post 居然用query传参
		req.Token = r.URL.Query().Get("token")
		req.Content = r.URL.Query().Get("content")

		ToUserId, err := strconv.ParseInt(r.URL.Query().Get("to_user_id"), 10, 64)
		if err != nil {
			httpx.OkJson(w, apiErr.InvalidParams)
			return
		}
		req.ToUserId = ToUserId

		ActionType, err := strconv.ParseInt(r.URL.Query().Get("action_type"), 10, 64)
		if err != nil {
			httpx.OkJson(w, apiErr.InvalidParams)
			return
		}
		req.ActionType = ActionType

		l := contact.NewSendMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendMessage(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
