package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/service/rpc/video/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/video/types/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListByAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListByAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListByAuthorLogic {
	return &GetVideoListByAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListByAuthorLogic) GetVideoListByAuthor(in *video.GetVideoListByAuthorRequest) (*video.GetVideoListByAuthorResponse, error) {
	//获得作者的视频列表（由新到旧）
	var videos []model.Video
	err := l.svcCtx.DBList.Mysql.Where("author_id = ?", in.AuthorId).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	// 封装查询结果
	videoList := make([]*video.VideoInfo, 0, len(videos))
	for _, v := range videos {
		videoInfo := &video.VideoInfo{
			Id:            int32(v.ID),
			AuthorId:      int32(v.AuthorId),
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int32(v.FavoriteCount),
			CommentCount:  int32(v.CommentCount),
		}
		videoList = append(videoList, videoInfo)
	}

	return &video.GetVideoListByAuthorResponse{
		VideoList: videoList,
	}, nil
}
