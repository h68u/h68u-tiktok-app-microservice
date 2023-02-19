package video

import (
	"bytes"
	"github.com/h2non/filetype"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"h68u-tiktok-app-microservice/common/error/apiErr"
	"h68u-tiktok-app-microservice/common/oss"
	"h68u-tiktok-app-microservice/common/utils"
	logic "h68u-tiktok-app-microservice/service/http/internal/logic/video"
	"h68u-tiktok-app-microservice/service/http/internal/svc"
	"h68u-tiktok-app-microservice/service/http/internal/types"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func PublishVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishVideoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 获取用户id
		userId, err := utils.GetUserIDFormToken(req.Token, svcCtx.Config.Auth.AccessSecret)
		if err != nil {
			httpx.Error(w, apiErr.InvalidToken)
			return
		}

		// 获取文件
		file, fileHeader, err := r.FormFile("data")
		if err != nil {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails(err.Error()))
			return
		}
		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		// 判断是否为视频
		tmpFile, err := fileHeader.Open()
		if err != nil {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails(err.Error()))
			return
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, tmpFile); err != nil {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails(err.Error()))
			return
		}
		if !filetype.IsVideo(buf.Bytes()) {
			httpx.Error(w, apiErr.FileIsNotVideo)
			return
		}
		if err = tmpFile.Close(); err != nil {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails(err.Error()))
			return
		}

		// 使用uuid重新生成文件名
		fileName := utils.GetUUID() + filepath.Ext(fileHeader.Filename)

		// 存储到oss
		ok, err := oss.UploadVideoToOss(svcCtx.AliyunClient, svcCtx.Config.OSS.BucketName, fileName, file)
		if err != nil {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails(err.Error()))
			return
		} else if !ok {
			httpx.Error(w, apiErr.FileUploadFailed.WithDetails("upload video to oss error"))
			return
		}

		videoUrl, imgUrl := oss.GetOssVideoUrlAndImgUrl(svcCtx.Config.OSS, fileName)

		// 调用rpc服务
		_, err = svcCtx.VideoRpc.PublishVideo(r.Context(), &video.PublishVideoRequest{
			Video: &video.VideoInfo{
				AuthorId: userId,
				Title:    req.Title,
				PlayUrl:  videoUrl,
				CoverUrl: imgUrl,
			},
		})

		if err != nil {
			logx.WithContext(r.Context()).Errorf("PublishVideo rpc error: %v", err)
			httpx.Error(w, apiErr.InternalError(r.Context(), err.Error()))
			return
		}

		//httpx.OkJson(w, apiErr.Success)

		l := logic.NewPublishVideoLogic(r.Context(), svcCtx)
		resp, err := l.PublishVideo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
