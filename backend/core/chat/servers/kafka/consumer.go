package kafka

import (
	"chat/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"time"
)

// NewConsumer Регистрация потребителя kafka
func NewConsumer(cfg *config.KafkaConsumer) *kafka.Reader {
	readerConf := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	}

	if cfg.Username != "" && cfg.Password != "" {
		mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
		if err != nil {
			panic(err)
		}

		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		readerConf.Dialer = dialer
	}

	kafkaConsumer := kafka.NewReader(readerConf)

	return kafkaConsumer
}

// NewConsumerDebeziumChats Регистрация потребителя kafka debezium
func NewConsumerDebeziumChats(cfg *config.KafkaConsumerDebeziumChats) *kafka.Reader {
	readerConf := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	}

	if cfg.Username != "" && cfg.Password != "" {
		mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
		if err != nil {
			panic(err)
		}

		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		readerConf.Dialer = dialer
	}

	kafkaConsumer := kafka.NewReader(readerConf)

	return kafkaConsumer
}
