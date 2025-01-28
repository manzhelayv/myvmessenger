package config

import (
	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	MONGOEngine string `short:"n" env:"DATASTORE_MONGO" default:"mongo"`
	MONGODB     string `short:"d" env:"DATASTORE_MONGO_DB" default:"test"`
	MONGOURL    string `short:"u" env:"DATASTORE_MONGO_URL" default:"mongodb://192.168.1.239:30102/?directConnection=true"`

	POSTGRESEngine   string `short:"p" env:"DATASTORE_POSTGRES" default:"postgres"`
	POSTGRESDB       string `short:"b" env:"DATASTORE_POSTGRES_DB" default:"postgres"`
	POSTGRESUSER     string `short:"c" env:"DATASTORE_POSTGRES_USER" default:"postgres"`
	POSTGRESPASSWORD string `short:"a" env:"DATASTORE_POSTGRES_PASSWORD" default:"postgres"`
	POSTGRESURL      string `short:"r" env:"DATASTORE_POSTGRES_URL" default:"192.168.1.239:30101"`

	CLICKHOUSEUSER     string `short:"P" env:"DATASTORE_CLICKHOUSE_USER" default:"default"`
	CLICKHOUSEPASSWORD string `short:"R" env:"DATASTORE_CLICKHOUSE_PASSWORD" default:""`
	CLICKHOUSEHOST     string `short:"S" env:"DATASTORE_CLICKHOUSE_HOST" default:"192.168.1.239"`
	CLICKHOUSEPORT     string `short:"T" env:"DATASTORE_CLICKHOUSE_PORT" default:"30104"`
	CLICKHOUSEDB       string `short:"U" env:"DATASTORE_CLICKHOUSE_DB" default:"default"`

	REDISDB       int    `short:"e" env:"DATASTORE_REDIS" default:"1"`
	REDISPASSWORD string `short:"i" env:"DATASTORE_REDIS_REDISPASSWORD" default:""`
	REDISURL      string `short:"k" env:"DATASTORE_REDIS_URL" default:"192.168.1.239:30105"`

	ListenAddr string `short:"l" env:"LISTEN" default:":8095"`

	GrpcListenAddr       string `short:"g" env:"GRPC_LISTEN" default:":4051"`
	F3GrpcListenAddr     string `short:"G" env:"GRPC_CLIENT_F3" default:":4052"`
	ServerGrpcListenAddr string `short:"K" env:"GRPC_CLIENT_SERVER" default:":4050"`

	KafkaProducer

	KafkaConsumer

	KafkaProducerDebeziumChats

	KafkaConsumerDebeziumChats
}

// KafkaProducer Производитель kafka
type KafkaProducer struct {
	Brokers  []string `short:"o" env:"BROKERS_PRODUCER_KAFKA" default:"192.168.1.239:30094"`
	Username string   `short:"s" env:"USERNAME_PRODUCER_KAFKA" default:""`
	Password string   `short:"t" env:"PASSWORD_PRODUCER_KAFKA" default:""`
	Topic    string   `short:"x" env:"TOPIC_PRODUCER_KAFKA" default:"mongo_chats"`
}

// KafkaConsumer Потребитель kafka
type KafkaConsumer struct {
	Brokers  []string `short:"H" env:"BROKERS_CONSUMER_KAFKA" default:"192.168.1.239:30094"`
	Username string   `short:"y" env:"USERNAME_CONSUMER_KAFKA" default:""`
	Password string   `short:"z" env:"PASSWORD_CONSUMER_KAFKA" default:""`
	Topic    string   `short:"O" env:"TOPIC_CONSUMER_KAFKA" default:"mongo_chats"`
	GroupID  string   `short:"v" env:"TOPIC_GROUP_ID_KAFKA" default:"chat-group-id"`
}

// KafkaConsumerDebeziumChats Производитель kafka debezium
type KafkaProducerDebeziumChats struct {
	Brokers  []string `short:"B" env:"BROKERS_PRODUCER_KAFKA_DEBEZIUM" default:"192.168.1.239:30094"`
	Username string   `short:"L" env:"USERNAME_PRODUCER_KAFKA_DEBEZIUM" default:""`
	Password string   `short:"E" env:"PASSWORD_PRODUCER_KAFKA_DEBEZIUM" default:""`
	Topic    string   `short:"X" env:"TOPIC_PRODUCER_KAFKA_DEBEZIUM" default:"chat"`
}

// KafkaConsumerDebeziumChats Потребитель kafka debezium
type KafkaConsumerDebeziumChats struct {
	Brokers  []string `short:"Z" env:"BROKERS_CONSUMER_KAFKA_DEBEZIUM" default:"192.168.1.239:30094"`
	Username string   `short:"Y" env:"USERNAME_CONSUMER_KAFKA_DEBEZIUM" default:""`
	Password string   `short:"D" env:"PASSWORD_CONSUMER_KAFKA_DEBEZIUM" default:""`
	Topic    string   `short:"I" env:"TOPIC_CONSUMER_KAFKA_DEBEZIUM" default:"chat"`
	GroupID  string   `short:"V" env:"TOPIC_GROUP_ID_KAFKA_DEBEZIUM" default:"debezium-chat-group-id"`
}

// Parse Парсит параметры и опции
func Parse() *Configuration {
	var opt Configuration

	p := flags.NewParser(&opt, flags.Default|flags.IgnoreUnknown)
	if _, err := p.Parse(); err != nil {
		panic(err)
	}

	return &opt
}
