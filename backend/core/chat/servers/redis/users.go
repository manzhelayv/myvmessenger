package redis

import (
	"github.com/go-redis/redis/v7"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(redis *redis.Client) *Redis {
	redisCli := &Redis{
		Client: redis,
	}

	return redisCli
}

// GetUser Получение tdid пользователя
func (c *Redis) GetUser(userTo string) string {
	userRedis := c.Client.HGetAll("users:user:" + userTo).Val()

	if userRedis != nil {
		return userRedis["Tdid"]
	}

	return ""
}
