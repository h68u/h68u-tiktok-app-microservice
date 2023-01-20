// Code generated by goctl. DO NOT EDIT!
// Source: video.proto

package videoclient

import (
	"context"

	"h68u-tiktok-app-microservice/service/rpc/video/types/video"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Comment                      = video.Comment
	CommentVideoRequest          = video.CommentVideoRequest
	Empty                        = video.Empty
	FavoriteVideoRequest         = video.FavoriteVideoRequest
	GetCommentListRequest        = video.GetCommentListRequest
	GetCommentListResponse       = video.GetCommentListResponse
	GetFavoriteVideoListRequest  = video.GetFavoriteVideoListRequest
	GetFavoriteVideoListResponse = video.GetFavoriteVideoListResponse
	GetVideoListByAuthorRequest  = video.GetVideoListByAuthorRequest
	GetVideoListByAuthorResponse = video.GetVideoListByAuthorResponse
	GetVideoListRequest          = video.GetVideoListRequest
	GetVideoListResponse         = video.GetVideoListResponse
	IsFavoriteVideoRequest       = video.IsFavoriteVideoRequest
	IsFavoriteVideoResponse      = video.IsFavoriteVideoResponse
	PublishVideoRequest          = video.PublishVideoRequest
	UnFavoriteVideoRequest       = video.UnFavoriteVideoRequest
	VideoInfo                    = video.VideoInfo

	Video interface {
		GetVideoList(ctx context.Context, in *GetVideoListRequest, opts ...grpc.CallOption) (*GetVideoListResponse, error)
		PublishVideo(ctx context.Context, in *PublishVideoRequest, opts ...grpc.CallOption) (*Empty, error)
		GetVideoListByAuthor(ctx context.Context, in *GetVideoListByAuthorRequest, opts ...grpc.CallOption) (*GetVideoListByAuthorResponse, error)
		FavoriteVideo(ctx context.Context, in *FavoriteVideoRequest, opts ...grpc.CallOption) (*Empty, error)
		UnFavoriteVideo(ctx context.Context, in *UnFavoriteVideoRequest, opts ...grpc.CallOption) (*Empty, error)
		GetFavoriteVideoList(ctx context.Context, in *GetFavoriteVideoListRequest, opts ...grpc.CallOption) (*GetFavoriteVideoListResponse, error)
		IsFavoriteVideo(ctx context.Context, in *IsFavoriteVideoRequest, opts ...grpc.CallOption) (*IsFavoriteVideoResponse, error)
		CommentVideo(ctx context.Context, in *CommentVideoRequest, opts ...grpc.CallOption) (*Empty, error)
		GetCommentList(ctx context.Context, in *GetCommentListRequest, opts ...grpc.CallOption) (*GetCommentListResponse, error)
	}

	defaultVideo struct {
		cli zrpc.Client
	}
)

func NewVideo(cli zrpc.Client) Video {
	return &defaultVideo{
		cli: cli,
	}
}

func (m *defaultVideo) GetVideoList(ctx context.Context, in *GetVideoListRequest, opts ...grpc.CallOption) (*GetVideoListResponse, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.GetVideoList(ctx, in, opts...)
}

func (m *defaultVideo) PublishVideo(ctx context.Context, in *PublishVideoRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.PublishVideo(ctx, in, opts...)
}

func (m *defaultVideo) GetVideoListByAuthor(ctx context.Context, in *GetVideoListByAuthorRequest, opts ...grpc.CallOption) (*GetVideoListByAuthorResponse, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.GetVideoListByAuthor(ctx, in, opts...)
}

func (m *defaultVideo) FavoriteVideo(ctx context.Context, in *FavoriteVideoRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.FavoriteVideo(ctx, in, opts...)
}

func (m *defaultVideo) UnFavoriteVideo(ctx context.Context, in *UnFavoriteVideoRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.UnFavoriteVideo(ctx, in, opts...)
}

func (m *defaultVideo) GetFavoriteVideoList(ctx context.Context, in *GetFavoriteVideoListRequest, opts ...grpc.CallOption) (*GetFavoriteVideoListResponse, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.GetFavoriteVideoList(ctx, in, opts...)
}

func (m *defaultVideo) IsFavoriteVideo(ctx context.Context, in *IsFavoriteVideoRequest, opts ...grpc.CallOption) (*IsFavoriteVideoResponse, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.IsFavoriteVideo(ctx, in, opts...)
}

func (m *defaultVideo) CommentVideo(ctx context.Context, in *CommentVideoRequest, opts ...grpc.CallOption) (*Empty, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.CommentVideo(ctx, in, opts...)
}

func (m *defaultVideo) GetCommentList(ctx context.Context, in *GetCommentListRequest, opts ...grpc.CallOption) (*GetCommentListResponse, error) {
	client := video.NewVideoClient(m.cli.Conn())
	return client.GetCommentList(ctx, in, opts...)
}