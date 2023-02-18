package mq

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const TypeLoseFriends = "friend:lose"

type LoseFriendsPayload struct {
	UserAId int32
	UserBId int32
}

func NewLoseFriendsTask(userAId int32, userBId int32) (*asynq.Task, error) {
	payload, err := json.Marshal(LoseFriendsPayload{UserAId: userAId, UserBId: userBId})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload for task %q", TypeLoseFriends)
	}
	return asynq.NewTask(TypeLoseFriends, payload), nil
}
