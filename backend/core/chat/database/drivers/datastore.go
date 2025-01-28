package drivers

import (
	"chat/models"
	"context"
)

type DbInterfase interface {
	Base
	Messages
}

type Base interface {
	Connect() error
	Close()
}

type Messages interface {
	// SetMessage Добавление сообщения в чат
	SetMessage(ctx context.Context, message *models.Messages) error

	// GetMessages Получение сообщений в чате
	GetMessages(ctx context.Context, userFrom, userTo string) ([]*models.Messages, error)

	// GetLastMessages Получение последних сообщений в чатах
	GetLastMessages(ctx context.Context, userFrom string) ([]*models.Messages, error)
}
