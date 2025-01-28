package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"server/config"
	"server/database"
	"server/database/drivers"
	application "server/servers"
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
	app.RegisterGrpcServer(application.GrpcRegister)
	app.RegisterRedisServer()
	app.RegisterHTTPServer()

	app.Start(errChan)

	wg.Wait()
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

func errorHandler(ctx context.Context, cancel context.CancelFunc, errChan chan error, logs *Logger) {
	select {
	case err := <-errChan:
		logs.Fatal(err)

		cancel()
		os.Exit(1)
	case <-ctx.Done():
		log.Println("Done")
		return
	}
}
