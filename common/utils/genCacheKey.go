package utils

import "strconv"

func GenFavoriteVideoCacheKey(userId, videoId int64) string {
	return "favorite_video_" + strconv.Itoa(int(userId)) + "_" + strconv.Itoa(int(videoId))
}

func GenFollowUserCacheKey(userId, followUserId int64) string {
	return "follow_user_" + strconv.Itoa(int(userId)) + "_" + strconv.Itoa(int(followUserId))
}
