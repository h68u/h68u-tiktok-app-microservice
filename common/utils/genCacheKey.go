package utils

func GenFavoriteVideoCacheKey(userId, videoId int64) string {
	return "favorite_video_" + Int642Str(userId) + "_" + Int642Str(videoId)
}

func GenFollowUserCacheKey(userId, followUserId int64) string {
	return "follow_user_" + Int642Str(userId) + "_" + Int642Str(followUserId)
}

func GenPopUserListCacheKey() string {
	return "pop_user_list"
}

func GenUserInfoCacheKey(userId int64) string {
	return "user_info_" + Int642Str(userId)
}

func GenPopVideoListCacheKey() string {
	return "pop_video_list"
}

func GenVideoInfoCacheKey(videoId int64) string {
	return "video_info_" + Int642Str(videoId)
}
