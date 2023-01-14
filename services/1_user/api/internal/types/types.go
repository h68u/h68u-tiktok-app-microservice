// Code generated by goctl. DO NOT EDIT.
package types

type RegisterRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type RegisterReply struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginReply struct {
	Code   int    `json:"status_code"`
	Msg    string `json:"status_msg"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type User struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type GetUserInfoRequest struct {
	UserId int    `form:"user_id"`
	Token  string `form:"token"`
}

type GetUserInfoReply struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
	User User   `json:"user"`
}

type FollowRequest struct {
	UserId     int    `form:"user_id"`
	Token      string `form:"token"`
	ToUserId   int    `form:"to_user_id"`
	ActionType int    `form:"action_type"`
}

type FollowReply struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

type FollowListRequest struct {
	UserId int    `form:"user_id"`
	Token  string `form:"token"`
}

type FollowListReply struct {
	Code  int    `json:"status_code"`
	Msg   string `json:"status_msg"`
	Users []User `json:"user_list"`
}

type FansListRequest struct {
	UserId int    `form:"user_id"`
	Token  string `form:"token"`
}

type FansListReply struct {
	Code  int    `json:"status_code"`
	Msg   string `json:"status_msg"`
	Users []User `json:"user_list"`
}
