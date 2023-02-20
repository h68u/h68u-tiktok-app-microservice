package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
	"time"
)

type GetVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListLogic) GetVideoList(in *video.GetVideoListRequest) (*video.GetVideoListResponse, error) {
	var videos []model.Video

	err := l.svcCtx.DBList.Mysql.
		Model(&model.Video{}).
		Where("created_at < ?", time.Unix(in.LatestTime, 0)).
		Order("created_at desc"). // 按照创建时间倒序，最新的在前面
		Limit(int(in.Num)).
		Find(&videos).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	var videoList []*video.VideoInfo
	for _, v := range videos {
		//把热门视频放缓存里
		if model.IsPopularVideo(v.FavoriteCount, v.CommentCount) {
			if l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenVideoInfoCacheKey(int64(v.ID))).Val() == 0 {
				// 缓存不存在，写入缓存
				if err := l.svcCtx.DBList.Redis.HSet(l.ctx, utils.GenVideoInfoCacheKey(int64(v.ID)), map[string]interface{}{
					"AuthorId":      v.AuthorId,
					"Title":         v.Title,
					"PlayUrl":       v.PlayUrl,
					"CoverUrl":      v.CoverUrl,
					"FavoriteCount": v.FavoriteCount,
					"CommentCount":  v.CommentCount,
					"CreatedAt":     v.CreatedAt.Unix(),
				}).Err(); err != nil {
					return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
				}
				// 放进热门视频列表
				err = l.svcCtx.DBList.Redis.LPush(l.ctx, utils.GenPopVideoListCacheKey(), v.ID).Err()
				if err != nil {
					return nil, status.Error(rpcErr.CacheError.Code, err.Error())
				}
			}

		}

		videoList = append(videoList, &video.VideoInfo{
			Id:            int64(v.ID),
			AuthorId:      v.AuthorId,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			CreateTime:    v.CreatedAt.Unix(),
		})
	}

	return &video.GetVideoListResponse{
		VideoList: videoList,
	}, nil
}
