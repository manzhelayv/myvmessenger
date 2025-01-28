package managers

import (
	"chat/database/drivers"
	"chat/models"
	redisCli "chat/servers/redis"
	"context"
	"encoding/json"
	"gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"time"
)

type Chat struct {
	dbMongo  drivers.DbInterfase
	redisCli *redisCli.Redis
	userGrpc protobuf.UsersClient
}

func NewChat(dbMongo drivers.DbInterfase, redisCli *redisCli.Redis, userGrpc protobuf.UsersClient) *Chat {
	return &Chat{
		dbMongo:  dbMongo,
		redisCli: redisCli,
		userGrpc: userGrpc,
	}
}

// AddMessage Сохранение сообщения
func (c *Chat) AddMessage(userFrom string, message []byte) (string, string, string, string, error) {
	if len(message) == 0 || userFrom == "" {
		return "", "", "", "", models.EmptyAddMessage
	}

	ctx := context.Background()

	now := time.Now().Format(time.RFC3339)
	date, _ := time.Parse(time.RFC3339, now)

	newMessage := &models.MessageWS{}
	if err := json.Unmarshal(message, newMessage); err != nil {
		return "", "", "", "", err
	}

	returnMessage := newMessage.Message
	if newMessage.File != "" {
		newMessage.Message = ""
	}

	chat := &models.Messages{
		UserFrom:      userFrom,
		UserTo:        newMessage.UserTo,
		Message:       newMessage.Message,
		File:          newMessage.File,
		CreatedAt:     date,
		UpdatedAt:     date,
		DateTimestamp: 0,
	}

	err := c.dbMongo.SetMessage(ctx, chat)
	if err != nil {
		return "", "", "", "", err
	}

	return returnMessage, newMessage.File, newMessage.Status, newMessage.UserTo, nil
}

// GetStatusUser Получение tdid контактов пользователя, с которыми есть переписка
func (c *Chat) GetStatusUser(tdid string) (map[string]bool, error) {
	ctx := context.Background()

	messages, err := c.dbMongo.GetLastMessages(ctx, tdid)
	if err != nil {
		return nil, err
	}

	if messages == nil {
		return nil, nil
	}

	usersTo := make(map[string]bool, len(messages))
	for _, userTo := range messages {
		usersTo[userTo.UserTo] = true
	}

	return usersTo, nil
}
