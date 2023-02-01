package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
		videoList = append(videoList, &video.VideoInfo{
			Id:            int32(v.ID),
			AuthorId:      int32(v.AuthorId),
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int32(v.FavoriteCount),
			CommentCount:  int32(v.CommentCount),
			CreateTime:    v.CreatedAt.Unix(),
		})
	}

	return &video.GetVideoListResponse{
		VideoList: videoList,
	}, nil
}
