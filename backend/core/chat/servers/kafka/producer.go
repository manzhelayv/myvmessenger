package kafka

import (
	"chat/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// NewProducer Регистрация производителя kafka
func NewProducer(cfg *config.KafkaProducer) *kafka.Writer {
	writer := kafka.Writer{
		Addr:       kafka.TCP(cfg.Brokers...),
		Topic:      cfg.Topic,
		BatchBytes: 1048576000,
	}

	if cfg.Username != "" && cfg.Password != "" {
		mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
		if err != nil {
			panic(err)
		}

		sharedTransport := &kafka.Transport{
			SASL: mechanism,
		}

		writer.Transport = sharedTransport
	}

	return &writer
}

// NewProducerDebeziumChats Регистрация производителя kafka debezium
func NewProducerDebeziumChats(cfg *config.KafkaProducerDebeziumChats) *kafka.Writer {
	writer := kafka.Writer{
		Addr:  kafka.TCP(cfg.Brokers...),
		Topic: cfg.Topic,
	}

	if cfg.Username != "" && cfg.Password != "" {
		mechanism, err := scram.Mechanism(scram.SHA512, cfg.Username, cfg.Password)
		if err != nil {
			panic(err)
		}

		sharedTransport := &kafka.Transport{
			SASL: mechanism,
		}

		writer.Transport = sharedTransport
	}

	return &writer
}
