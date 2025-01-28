package server

import (
	"chat/config"
	"chat/database/drivers"
	"chat/managers"
	"chat/servers/clickhouse"
	"context"
	"errors"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	collectionChat        = "chats"
	collectionChatMv      = "chats_mv"
	collectionChatChanges = "chats_changes"
)

type ClickHouse struct {
	config drivers.Config

	clickHouseCli *clickhouse.ClickHouseClient
	debeziumChat  *managers.Debezium
	ctx           context.Context
	opt           *config.Configuration
}

func NewClickHouseServer(ctx *context.Context, opt *config.Configuration) *ClickHouse {
	config := drivers.Config{
		DbUser:     opt.CLICKHOUSEUSER,
		DbPassword: opt.CLICKHOUSEPASSWORD,
		DbHost:     opt.CLICKHOUSEHOST,
		DbPort:     opt.CLICKHOUSEPORT,
		DbName:     opt.CLICKHOUSEDB,
	}

	c := &ClickHouse{
		config: config,
		ctx:    *ctx,
		opt:    opt,
	}

	c.clickHouseCli = clickhouse.NewClickHouseClient(c.opt)

	c.debeziumChat = managers.NewDebezium(c.clickHouseCli)

	return c
}

func (c *ClickHouse) Start() error {
	var err error

	c.clickHouseCli.Client, err = sqlx.Open("clickhouse",
		fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
			c.config.DbUser,
			c.config.DbPassword,
			c.config.DbHost,
			c.config.DbPort,
			c.config.DbName))

	if err != nil {
		return errors.New(fmt.Sprintf("[ERRORS] Connect clickhouse: %s ", err))
	}

	if pingErr := c.clickHouseCli.Client.Ping(); pingErr != nil {
		return errors.New(fmt.Sprintf("[ERROR] Ping clickhouse: %s ", pingErr.Error()))
	}

	if err = c.createTablesIfNotExists(); err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Connect создание таблицы clickhouse: %s ", err.Error()))
	}

	if err = c.createMaterializedTablesIfNotExists(); err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Connect создание представления materialized clickhouse : %s ", err.Error()))
	}

	go c.StartConsumerDebeziumChats()

	return nil
}

func (c *ClickHouse) Stop() error {
	err := c.clickHouseCli.Client.Close()
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] Connect clickhouse: %s", err.Error()))
	}

	return nil
}

func (c *ClickHouse) createTablesIfNotExists() error {
	for collection, indexes := range c.tablesToCreating() {
		sqlStatement := "CREATE TABLE IF NOT EXISTS " + collection + "("

		for name, model := range indexes {
			sqlStatement += name + " " + model + ", "
		}

		sqlStatement = sqlStatement[:len(sqlStatement)-2]
		sqlStatement += ")\n"

		driver := c.getTableDriver(collection)

		sqlStatement += driver

		_, err := c.clickHouseCli.Client.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}

	if err := c.alterTable(); err != nil {
		return err
	}

	return c.ensureIndexes()
}

func (c *ClickHouse) tablesToCreating() map[string]map[string]string {
	return map[string]map[string]string{
		collectionChat: {
			"id":         "varchar(30)",
			"user_from":  "varchar(10)",
			"user_to":    "varchar(10)",
			"message":    "String",
			"file":       "String",
			"created_at": "DateTime",
			"updated_at": "DateTime",
		},
		collectionChatChanges: {
			"`before.id`":         "Nullable(varchar(30))",
			"`before.user_from`":  "Nullable(varchar(10))",
			"`before.user_to`":    "Nullable(varchar(10))",
			"`before.message`":    "Nullable(String)",
			"`before.file`":       "Nullable(String)",
			"`before.created_at`": "Nullable(DateTime)",
			"`before.updated_at`": "Nullable(DateTime)",
			"`after.id`":          "Nullable(varchar(30))",
			"`after.user_from`":   "Nullable(varchar(10))",
			"`after.user_to`":     "Nullable(varchar(10))",
			"`after.message`":     "Nullable(String)",
			"`after.file`":        "Nullable(String)",
			"`after.created_at`":  "Nullable(DateTime)",
			"`after.updated_at`":  "Nullable(DateTime)",
			"`op`":                "LowCardinality(String)",
		},
	}
}

func (c *ClickHouse) createMaterializedTablesIfNotExists() error {
	for collection, indexes := range c.tablesMaterializedToCreating() {
		sqlStatement := "CREATE MATERIALIZED VIEW IF NOT EXISTS " + collection + " TO " + collectionChat + "("

		for name, model := range indexes {
			sqlStatement += name + " " + model + ", "
		}

		sqlStatement = sqlStatement[:len(sqlStatement)-2]
		sqlStatement += ")\n"

		driver := c.getTableDriver(collection)

		sqlStatement += driver

		_, err := c.clickHouseCli.Client.Exec(sqlStatement)
		if err != nil {
			return err
		}
	}

	if err := c.alterTable(); err != nil {
		return err
	}

	return c.ensureIndexes()
}

func (c *ClickHouse) tablesMaterializedToCreating() map[string]map[string]string {
	return map[string]map[string]string{
		collectionChatMv: {
			"id":         "Nullable(varchar(30))",
			"user_from":  "Nullable(varchar(10))",
			"user_to":    "Nullable(varchar(10))",
			"message":    "Nullable(String)",
			"file":       "Nullable(String)",
			"created_at": "Nullable(DateTime)",
			"updated_at": "Nullable(DateTime)",
		},
	}
}

func (c *ClickHouse) getTableDriver(driver string) string {
	str := ""
	if driver == collectionChat {
		str = " ENGINE = ReplacingMergeTree(updated_at) ORDER BY (id);"
	}

	if driver == collectionChatChanges {
		str = " ENGINE = MergeTree ORDER BY tuple();"
	}

	if driver == collectionChatMv {
		str = "AS SELECT " +
			"if(op = 'd', before.id, after.id) AS id,\n   " +
			"if(op = 'd', before.user_from, after.user_from) AS user_from,\n" +
			"if(op = 'd', before.user_to, after.user_to) AS user_to,\n" +
			"if(op = 'd', before.message, after.message) AS message,\n   " +
			"if(op = 'd', before.file, after.file) AS file,\n   " +
			"if(op = 'd', before.created_at, after.created_at) AS created_at,\n  " +
			"if(op = 'd', before.updated_at, after.updated_at) AS updated_at\n  " +
			"FROM " + collectionChatChanges + "\n" +
			"WHERE (op = 'c') OR (op = 'r') OR (op = 'u') OR (op = 'd')"
	}

	return str
}

func (c *ClickHouse) alterTable() error {
	for collection, indexes := range c.alterToCreating() {
		for name, model := range indexes {
			sqlStatement := "ALTER TABLE " + collection + " ADD COLUMN IF NOT EXISTS " + name + " " + model

			_, err := c.clickHouseCli.Client.Exec(sqlStatement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ClickHouse) alterToCreating() map[string]map[string]string {
	return map[string]map[string]string{}
}

func (c *ClickHouse) ensureIndexes() error {
	for collection, indexes := range c.indexesToCreating() {
		for name, model := range indexes {
			sqlStatement := "CREATE INDEX IF NOT EXISTS " + name + " ON " + collection + " (" + model + ")"

			_, err := c.clickHouseCli.Client.Exec(sqlStatement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ClickHouse) indexesToCreating() map[string]map[string]string {
	return map[string]map[string]string{}
}
