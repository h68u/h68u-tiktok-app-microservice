package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/rpcErr"
	"h68u-tiktok-app-microservice/services/2_video/model"
	"h68u-tiktok-app-microservice/services/2_video/rpc/internal/svc"
	"h68u-tiktok-app-microservice/services/2_video/rpc/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteVideoListLogic {
	return &GetFavoriteVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteVideoListLogic) GetFavoriteVideoList(in *video.GetFavoriteVideoListRequest) (*video.GetFavoriteVideoListResponse, error) {
	// 获取用户点赞的视频id （要根据时间排序，从新到旧）
	var favorites []model.Favorite

	if err := l.svcCtx.DBList.Mysql.
		Where("user_id = ?", in.UserId).
		Preload("Video").
		Order("created_at desc").
		Find(&favorites).Error; err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	// 封装查询结果
	videoList := make([]*video.Video, 0, len(favorites))
	for _, v := range favorites {
		// 可能存在脏数据,需要判断视频是否存在
		if v.Video.ID == 0 {
			// TODO:异步删除脏数据
			continue
		}

		videoInfo := &video.Video{
			Id:            int32(v.Video.ID),
			AuthorId:      int32(v.Video.AuthorId),
			Title:         v.Video.Title,
			PlayUrl:       v.Video.PlayUrl,
			CoverUrl:      v.Video.CoverUrl,
			FavoriteCount: int32(v.Video.FavoriteCount),
			CommentCount:  int32(v.Video.CommentCount),
		}

		videoList = append(videoList, videoInfo)
	}

	return &video.GetFavoriteVideoListResponse{
		VideoList: videoList,
	}, nil
}
