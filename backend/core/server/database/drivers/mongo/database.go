package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"server/config"
	"server/database/drivers"
)

type Mongo struct {
	DbUrl    string
	DbName   string
	DNEngine string

	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	opt    *config.Configuration
}

const CollectionUsers = "users"

func New(config drivers.Config, ctx *context.Context, opt *config.Configuration) drivers.DbInterfase {
	return &Mongo{
		DbUrl:    config.DbUrl,
		DbName:   config.DbName,
		DNEngine: config.DNEngine,
		ctx:      *ctx,
		opt:      opt,
	}
}

func (m *Mongo) Connect() error {
	var err error

	m.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(m.DbUrl))
	if err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Connect mongo: %s ", err.Error()))
	}

	err = m.client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Ping mongo: %s ", err.Error()))
	}

	log.Println("[INFO] Starting MongoDB ", m.opt.MONGOURL)

	m.db = m.client.Database(m.DbName)

	go m.gracefulShutdown()

	return nil
}

func (m *Mongo) Close() {
	err := m.client.Disconnect(m.ctx)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] Connect mongodb: %s", err.Error()))
	}
}

func (p *Mongo) gracefulShutdown() {
	<-p.ctx.Done()
	p.Close()
	log.Println("Shutting down MongoDB database")
}
