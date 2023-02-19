package mq

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const TypeTryMakeFriends = "friend:make"

type TryMakeFriendsPayload struct {
	UserAId int64
	UserBId int64
}

func NewTryMakeFriendsTask(userAId int64, userBId int64) (*asynq.Task, error) {
	payload, err := json.Marshal(TryMakeFriendsPayload{UserAId: userAId, UserBId: userBId})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload for task %q", TypeTryMakeFriends)
	}
	return asynq.NewTask(TypeTryMakeFriends, payload), nil
}
