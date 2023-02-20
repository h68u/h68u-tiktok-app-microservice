package handler

import (
	"context"
	"github.com/hibiken/asynq"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
)

func (l *AsynqServer) syncVideoInfoCacheHandler(ctx context.Context, t *asynq.Task) error {
	res, err := l.svcCtx.Redis.LRange(ctx, utils.GenPopVideoListCacheKey(), 0, -1).Result()
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	for _, v := range res {
		videoId := utils.Str2Int64(v)
		// 读取缓存
		videoInfo, err := l.svcCtx.Redis.HGetAll(ctx, utils.GenVideoInfoCacheKey(videoId)).Result()
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
		// 更新视频信息
		_, err = l.svcCtx.VideoRpc.UpdateVideo(ctx, &video.UpdateVideoRequest{
			Video: &video.VideoInfo{
				Id:            videoId,
				AuthorId:      utils.Str2Int64(videoInfo["AuthorId"]),
				Title:         videoInfo["Title"],
				PlayUrl:       videoInfo["PlayUrl"],
				CoverUrl:      videoInfo["CoverUrl"],
				FavoriteCount: utils.Str2Int64(videoInfo["FavoriteCount"]),
				CommentCount:  utils.Str2Int64(videoInfo["CommentCount"]),
				CreateTime:    utils.Str2Int64(videoInfo["CreateTime"]),
			},
		})
		if err != nil {
			l.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
