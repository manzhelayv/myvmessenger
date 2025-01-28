package clickhouse

import (
	"chat/config"
	kafkaCli "chat/servers/kafka"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type ClickHouseClient struct {
	Client                     *sqlx.DB
	KafkaProducerDebeziumChats *kafka.Writer
	KafkaConsumerDebeziumChats *kafka.Reader
}

func NewClickHouseClient(opt *config.Configuration) *ClickHouseClient {
	c := &ClickHouseClient{}

	c.RegisterKafkaProducerDebeziumChats(&opt.KafkaProducerDebeziumChats)
	c.RegisterKafkaConsumerDebeziumChats(&opt.KafkaConsumerDebeziumChats)

	return c
}

// RegisterKafkaProducerDebeziumChats Регистрация производителя kafka debezium, таблицы mongo chats
func (c *ClickHouseClient) RegisterKafkaProducerDebeziumChats(opt *config.KafkaProducerDebeziumChats) {
	c.KafkaProducerDebeziumChats = kafkaCli.NewProducerDebeziumChats(opt)
}

// RegisterKafkaConsumerDebeziumChats Регистрация потребителя kafka debezium, таблицы mongo chats
func (c *ClickHouseClient) RegisterKafkaConsumerDebeziumChats(opt *config.KafkaConsumerDebeziumChats) {
	c.KafkaConsumerDebeziumChats = kafkaCli.NewConsumerDebeziumChats(opt)
}
