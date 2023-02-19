package mq

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

const TypeDelCache = "cache:del"

type DelCachePayload struct {
	Key string
}

func NewDelCacheTask(key string) (*asynq.Task, error) {
	payload, err := json.Marshal(DelCachePayload{Key: key})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload for task %q", TypeDelCache)
	}
	return asynq.NewTask(TypeDelCache, payload), nil
}
