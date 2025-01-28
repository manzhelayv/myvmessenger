package main

import (
	"context"
	"f3/config"
	application "f3/servers"
	"log"
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
	app.RegisterGrpcServer(application.GrpcRegister)
	app.RegisterRedisServer()
	app.RegisterHTTPServer()
	app.RegisterMinioServer()

	app.Start(errChan)

	wg.Wait()
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
