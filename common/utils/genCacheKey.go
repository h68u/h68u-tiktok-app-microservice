package utils

import "strconv"

func GenFavoriteVideoCacheKey(userId, videoId int32) string {
	return "favorite_video_" + strconv.Itoa(int(userId)) + "_" + strconv.Itoa(int(videoId))
}

func GenFollowUserCacheKey(userId, followUserId int32) string {
	return "follow_user_" + strconv.Itoa(int(userId)) + "_" + strconv.Itoa(int(followUserId))
}
