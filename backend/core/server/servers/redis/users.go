package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"server/models"
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

// GetUser Получение пользователя по email или login
func (c *Redis) GetUser(emailOrLogin string) *models.Users { //, password string
	userRedis := c.Client.HGetAll("users:user:" + emailOrLogin).Val()

	if userRedis["Email"] == emailOrLogin { // && userRedis["Password"] == password
		u := &models.Users{
			Email:    userRedis["Email"],
			Password: userRedis["Password"],
			Name:     userRedis["Name"],
			Login:    userRedis["Login"],
			Tdid:     userRedis["Tdid"],
		}

		return u
	}

	return nil
}

// CreateUser Добавление пользователя
func (c *Redis) CreateUser(data map[string]interface{}) (string, error) {
	pipeline := c.Client.TxPipeline()
	pipeline.HMSet(fmt.Sprintf("users:user:%s", data["Email"]), data)
	pipeline.HMSet(fmt.Sprintf("users:user:%s", data["Login"]), data)

	if _, err := pipeline.Exec(); err != nil {
		log.Println("pipeline err in CreateUser: ", err)
		return "", err
	}

	return "", nil
}
