package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/config"

	//"database/sql"
	"github.com/go-pg/pg/v10"
	"server/database/drivers"
	"time"
)

const (
	collectionUsers    = "users"
	collectionContacts = "contacts"
	collectionProfile  = "profiles"
)

type Postgres struct {
	config *pg.Options
	client *pg.DB

	ctx context.Context

	retries           int
	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	opt *config.Configuration
}

func New(config drivers.Config, ctx *context.Context, opt *config.Configuration) drivers.DbInterfase {
	configOptions := &pg.Options{
		Addr:     opt.POSTGRESURL,
		User:     opt.POSTGRESUSER,
		Password: opt.POSTGRESPASSWORD,
		Database: opt.POSTGRESDB,
	}

	return &Postgres{
		config: configOptions,
		ctx:    *ctx,
		opt:    opt,
	}
}

func (p *Postgres) Connect() error {
	var err error

	p.client = pg.Connect(p.config)

	if err = p.Ping(); err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Connect postgres: %s", err.Error()))
	}

	if err = p.createTablesIfNotExists(); err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Connect create table postgres: %s ", err.Error()))
	}

	go p.gracefulShutdown()

	log.Println("[INFO] Starting PostgresSQL ", p.opt.POSTGRESURL)

	return nil
}

func (p *Postgres) Ping() error {
	return p.client.Ping(context.Background())
}

func (p *Postgres) Name() string {
	return "Postgres"
}

func (p *Postgres) Close() {
	err := p.client.Close()
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] Connect postgres: %s", err.Error()))
	}
}

func (p *Postgres) gracefulShutdown() {
	<-p.ctx.Done()
	p.Close()
	log.Println("Shutting down PostgreSQL database")
}

func (p *Postgres) createTablesIfNotExists() error {
	for collection, indexes := range p.tablesToCreating() {
		sqlStatement := "CREATE TABLE IF NOT EXISTS " + collection + "("

		for name, model := range indexes {
			sqlStatement += name + " " + model + ", "
		}

		sqlStatement = sqlStatement[:len(sqlStatement)-2]
		sqlStatement += ")"

		_, err := p.client.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}

	if err := p.alterTable(); err != nil {
		return err
	}

	return p.ensureIndexes()
}

func (m *Postgres) tablesToCreating() map[string]map[string]string {
	return map[string]map[string]string{
		collectionUsers: {
			"id":         "serial PRIMARY KEY",
			"tdid":       "VARCHAR(10) NOT NULL UNIQUE",
			"email":      "VARCHAR(20) NOT NULL UNIQUE",
			"login":      "VARCHAR(20) NOT NULL UNIQUE",
			"name":       "VARCHAR(20) NOT NULL",
			"password":   "VARCHAR(512) NOT NULL",
			"created_at": "TIMESTAMP NOT NULL",
			"updated_at": "TIMESTAMP NOT NULL",
		},
		collectionContacts: {
			"id":         "serial PRIMARY KEY",
			"user_from":  "VARCHAR(10) NOT NULL",
			"user_to":    "VARCHAR(10) NOT NULL",
			"created_at": "TIMESTAMP NOT NULL",
			"updated_at": "TIMESTAMP NOT NULL",
		},
		collectionProfile: {
			"profile_id":           "serial PRIMARY KEY",
			"image":                "VARCHAR(50) NULL",
			"user_id":              "integer NOT NULL",
			"FOREIGN KEY(user_id)": " REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE",
		},
	}
}

func (p *Postgres) alterTable() error {
	for collection, indexes := range p.alterToCreating() {
		for name, model := range indexes {
			sqlStatement := "ALTER TABLE " + collection + " ADD COLUMN IF NOT EXISTS " + name + " " + model

			_, err := p.client.Exec(sqlStatement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Postgres) alterToCreating() map[string]map[string]string {
	return map[string]map[string]string{
		collectionContacts: {
			"user_to_name": "VARCHAR(20) NOT NULL",
		},
		collectionUsers: {
			"phone": "VARCHAR(20) NULL",
		},
	}
}

func (p *Postgres) ensureIndexes() error {
	for collection, indexes := range p.indexesToCreating() {
		for name, model := range indexes {
			sqlStatement := "CREATE INDEX IF NOT EXISTS " + name + " ON " + collection + " (" + model + ")"

			_, err := p.client.Exec(sqlStatement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Postgres) indexesToCreating() map[string]map[string]string {
	return map[string]map[string]string{
		collectionUsers: {
			"tdid":               "tdid",
			"email":              "email",
			"email_and_tdid":     "email, tdid",
			"email_and_password": "email, password",
		},
	}
}
