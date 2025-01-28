package database

import (
	"chat/config"
	"chat/database/drivers"
	"chat/database/drivers/mongo"
	"chat/database/drivers/postgres"
	"context"
	"errors"
	"fmt"
	"log"
)

const MONGO = "mongo"
const POSTGRES = "postgres"

//const CLICKHOUSE = "clickhouse"

type DbEngine = func(config drivers.Config, ctx *context.Context, opt *config.Configuration) drivers.DbInterfase

type DbFactory string

var dbFactories = make(map[DbFactory]DbEngine, 2)

func init() {
	Register(MONGO, mongo.New)
	Register(POSTGRES, postgres.New)
	//Register(CLICKHOUSE, clickhouse.New)
}

func Register(name DbFactory, fn DbEngine) {
	if _, ok := dbFactories[name]; ok {
		log.Println("Ok DB: ", name)
	} else {
		dbFactories[name] = fn
	}
}

func NewDatabase(fabric DbFactory, config drivers.Config, ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error) {
	if engine, ok := dbFactories[fabric]; ok {
		return engine(config, ctx, opt), nil
	} else {
		return nil, errors.New(fmt.Sprintf("Error connect %s", fabric))
	}
}
