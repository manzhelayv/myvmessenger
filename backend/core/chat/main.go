package main

import (
	"chat/config"
	"chat/database"
	"chat/database/drivers"
	application "chat/servers"
	"context"
	logs "log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	opt := config.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := NewLogger()

	app := application.Instance(ctx, log.infoLog, log.errorLog, opt)

	errChan := make(chan error)
	go errorHandler(ctx, cancel, errChan, log)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		cancel()
	}()

	wg := &sync.WaitGroup{}

	app.RegisterWaitGroup(wg)
	app.RegisterDatabasePostgres(database.POSTGRES, dbFactoryPostgres)
	app.RegisterDatabaseMongo(database.MONGO, dbFactoryMongo)
	app.RegisterRedisServer()
	app.RegisterKafkaProducer(&opt.KafkaProducer)
	app.RegisterKafkaConsumer(&opt.KafkaConsumer)
	app.RegisterClickHouseServer(opt)
	app.RegisterHTTPServer()

	app.Start(errChan)

	wg.Wait()

	defer func() {
		if r := recover(); r != nil {
			logs.Println("Recovered in f", r)
		}
	}()
}

func dbFactoryPostgres(ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error) {
	return database.NewDatabase(database.POSTGRES, drivers.Config{
		DbUrl:    opt.POSTGRESURL,
		DbName:   opt.POSTGRESDB,
		DNEngine: opt.POSTGRESEngine,
	}, ctx, opt)
}

func dbFactoryMongo(ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error) {
	return database.NewDatabase(database.MONGO, drivers.Config{
		DbUrl:    opt.MONGOURL,
		DbName:   opt.MONGODB,
		DNEngine: opt.MONGOEngine,
	}, ctx, opt)
}

func errorHandler(ctx context.Context, cancel context.CancelFunc, errChan chan error, log *Logger) {
	select {
	case err := <-errChan:
		log.Fatal(err)

		cancel()
		os.Exit(1)
	case <-ctx.Done():
		logs.Println("Done")
		return
	}
}
