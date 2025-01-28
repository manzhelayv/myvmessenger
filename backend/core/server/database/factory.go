package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/config"
	"server/database/drivers"
	"server/database/drivers/mongo"
	"server/database/drivers/postgres"
)

const MONGO = "mongo"
const POSTGRES = "postgres"

type DbEngine = func(config drivers.Config, ctx *context.Context, opt *config.Configuration) drivers.DbInterfase

type DbFactory string

var dbFactories = make(map[DbFactory]DbEngine, 2)

func init() {
	Register(MONGO, mongo.New)
	Register(POSTGRES, postgres.New)
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
