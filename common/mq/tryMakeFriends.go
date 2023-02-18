package mq

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const TypeTryMakeFriends = "friend:make"

type TryMakeFriendsPayload struct {
	UserAId int32
	UserBId int32
}

func NewTryMakeFriendsTask(userAId int32, userBId int32) (*asynq.Task, error) {
	payload, err := json.Marshal(TryMakeFriendsPayload{UserAId: userAId, UserBId: userBId})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload for task %q", TypeTryMakeFriends)
	}
	return asynq.NewTask(TypeTryMakeFriends, payload), nil
}
