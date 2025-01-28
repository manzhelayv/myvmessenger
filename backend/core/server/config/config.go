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

	REDISDB       int    `short:"e" env:"DATASTORE_REDIS" default:"1"`
	REDISPASSWORD string `short:"i" env:"DATASTORE_REDIS_REDISPASSWORD" default:""`
	REDISURL      string `short:"k" env:"DATASTORE_REDIS_URL" default:"192.168.1.239:30105"`

	ListenAddr string `short:"l" env:"LISTEN" default:":8096"`

	GrpcListenAddr   string `short:"g" env:"GRPC_LISTEN" default:":4050"`
	F3GrpcListenAddr string `short:"G" env:"GRPC_CLIENT_F3" default:":4052"`
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
