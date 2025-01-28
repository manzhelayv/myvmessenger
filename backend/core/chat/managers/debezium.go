package managers

import (
	"chat/models"
	clickhouseCli "chat/servers/clickhouse"
	"context"
)

type Debezium struct {
	clickhouseCli *clickhouseCli.ClickHouseClient
}

func NewDebezium(clickhouseCli *clickhouseCli.ClickHouseClient) *Debezium {
	return &Debezium{
		clickhouseCli: clickhouseCli,
	}
}

// AddMessageDebeziumChats Добавление сообщений в clickhouse из postgresql из коннектора кафки
func (c *Debezium) AddMessageDebeziumChats(chat *models.ChatsChanges) error {
	if chat == nil {
		return models.EmptyAddMessage
	}

	ctx := context.Background()

	err := c.clickhouseCli.SetMessageDebeziumChats(ctx, chat)
	if err != nil {
		return err
	}

	return nil
}
