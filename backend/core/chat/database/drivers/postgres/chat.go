package postgres

import (
	"chat/models"
	"context"
)

func (p Postgres) GetLastMessages(ctx context.Context, userFromOrUserTo string) ([]*models.Messages, error) {
	return nil, nil
}

func (p Postgres) GetMessages(ctx context.Context, userFrom, userTo string) ([]*models.Messages, error) {
	return nil, nil
}

func (p Postgres) SetMessage(ctx context.Context, message *models.Messages) error {
	return nil
}
