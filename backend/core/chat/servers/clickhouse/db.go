package clickhouse

import (
	"chat/models"
	"context"
	"fmt"
)

// SetMessageDebeziumChats Добавляет данные из кафки в clickhouse таблицу chats_changes
func (c ClickHouseClient) SetMessageDebeziumChats(ctx context.Context, chat *models.ChatsChanges) error {
	if chat == nil {
		return fmt.Errorf("Message is nil")
	}

	dateAfterCreatedAt := chat.After.CreatedAt.Date / 1000
	dateAfterUpdatedAt := chat.After.UpdatedAt.Date / 1000

	dateBeforeCreatedAt := chat.Before.CreatedAt.Date / 1000
	dateBeforeUpdatedAt := chat.Before.UpdatedAt.Date / 1000

	_, err := c.Client.ExecContext(
		ctx,
		"INSERT INTO chats_changes "+
			"(after.id, after.user_from, after.user_to, after.message, after.file, after.created_at, after.updated_at, "+
			"before.id, before.user_from, before.user_to, before.message, before.file, before.created_at, before.updated_at, "+
			"op) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
		chat.After.Id.Oid,
		chat.After.UserFrom,
		chat.After.UserTo,
		chat.After.Message,
		chat.After.File,
		dateAfterCreatedAt,
		dateAfterUpdatedAt,
		chat.Before.Id.Oid,
		chat.Before.UserFrom,
		chat.Before.UserTo,
		chat.Before.Message,
		chat.Before.File,
		dateBeforeCreatedAt,
		dateBeforeUpdatedAt,
		chat.Op,
	)
	if err != nil {
		return fmt.Errorf(" Обогащение таблицы chats_changes : %#v, :%w", chat, err)
	}

	return nil
}
