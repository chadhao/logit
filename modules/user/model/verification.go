package model

import (
	"time"
)

type Redis struct {
	Key            string
	ExpireDuration time.Duration
}

func (r *Redis) Get() (string, error) {
	return redisClient.Get(r.Key).Result()
}

func (r *Redis) Set(value string) {
	redisClient.Set(r.Key, value, r.ExpireDuration)
}

func (r *Redis) Expire() {
	redisClient.ExpireAt(r.Key, time.Now())
}
