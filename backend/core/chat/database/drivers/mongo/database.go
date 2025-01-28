package mongo

import (
	"chat/config"
	"chat/database/drivers"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 15 * time.Second

	CollectionChat = "chat"
)

type Mongo struct {
	DbUrl    string
	DbName   string
	DNEngine string

	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	opt    *config.Configuration

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration
}

func New(config drivers.Config, ctx *context.Context, opt *config.Configuration) drivers.DbInterfase {
	return &Mongo{
		DbUrl:             config.DbUrl,
		DbName:            config.DbName,
		DNEngine:          config.DNEngine,
		ctx:               *ctx,
		opt:               opt,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
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

	return m.ensureIndexes()
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
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

// ensureIndexes Убеждается что все индексы построены
func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.ensureIdxTimeout)
	defer cancel()

	for collection, indexes := range m.indexesToCreating() {
		col := m.db.Collection(collection)

		needCreate := make([]mongo.IndexModel, 0)

		for name, model := range indexes {
			exists, err := m.indexExistsByName(ctx, col, name)
			if err != nil {
				return err
			}
			if !exists {
				needCreate = append(needCreate, model)
			}
		}

		if len(needCreate) == 0 {
			continue
		}

		createOpts := options.CreateIndexes().SetMaxTime(m.ensureIdxTimeout)
		_, err := col.Indexes().CreateMany(ctx, needCreate, createOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

// indexExistsByName проверяет существование индекса с именем name.
func (m *Mongo) indexExistsByName(ctx context.Context, collection *mongo.Collection, name string) (bool, error) {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	for cur.Next(ctx) {
		if name == cur.Current.Lookup("name").StringValue() {
			return true, nil
		}
	}

	return false, nil
}

// ensureUsersIndexes Строит индексы для коллекции
func (m *Mongo) indexesToCreating() map[string]map[string]mongo.IndexModel {
	return map[string]map[string]mongo.IndexModel{
		CollectionChat: {
			"user_from":             {Keys: bson.D{{"user_from", 1}, {"_id", 1}}, Options: options.Index().SetUnique(true).SetName("user_from")},
			"user_to":               {Keys: bson.D{{"user_to", 1}, {"_id", 1}}, Options: options.Index().SetUnique(true).SetName("user_to")},
			"user_from_and_user_to": {Keys: bson.D{{"user_from", 1}, {"user_to", 1}, {"_id", 1}}, Options: options.Index().SetName("user_from_and_user_to")},
		},
	}
}
