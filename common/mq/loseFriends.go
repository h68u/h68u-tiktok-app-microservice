package mq

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const TypeLoseFriends = "friend:lose"

type LoseFriendsPayload struct {
	UserAId int64
	UserBId int64
}

func NewLoseFriendsTask(userAId int64, userBId int64) (*asynq.Task, error) {
	payload, err := json.Marshal(LoseFriendsPayload{UserAId: userAId, UserBId: userBId})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload for task %q", TypeLoseFriends)
	}
	return asynq.NewTask(TypeLoseFriends, payload), nil
}
