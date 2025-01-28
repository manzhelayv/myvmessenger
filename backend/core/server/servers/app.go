package servers

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"server/config"
	"server/database"
	"server/database/drivers"
	redisCli "server/servers/redis"
	"sync"
)

var (
	once      sync.Once
	singleton *application
)

// Основная структура приложения
type application struct {
	servers    []Servers
	wg         *sync.WaitGroup
	ctx        context.Context
	InfoLog    *log.Logger
	errorLog   *log.Logger
	dbPostgres drivers.DbInterfase
	dbMongo    drivers.DbInterfase
	redisCli   redisCli.Redis
	opt        *config.Configuration
}

func Instance(ctx context.Context, infoLog *log.Logger, errorLog *log.Logger, opt *config.Configuration) *application {
	once.Do(func() {
		singleton = &application{
			ctx:      ctx,
			InfoLog:  infoLog,
			opt:      opt,
			errorLog: errorLog,
		}
	})

	return singleton
}

// RegisterWaitGroup Добавление общего WaitGroup
func (a *application) RegisterWaitGroup(wg *sync.WaitGroup) {
	a.wg = wg
}

// RegisterRedisServer Регистрация redis сервера
func (a *application) RegisterRedisServer() {
	srv := NewRedisServer(a.ctx, a.wg, a)

	a.servers = append(a.servers, srv)
}

// RegisterHTTPServer Регистрация HTTP сервера
func (a *application) RegisterHTTPServer() {
	srv := NewHTTPServer(a)

	a.servers = append(a.servers, srv)
}

// RegisterGrpcServer Регистрация GRPC сервера
func (a *application) RegisterGrpcServer(fn func(server *grpc.Server, db drivers.DbInterfase)) {
	srv := NewGrpcServer(fn, a)

	a.servers = append(a.servers, srv)
}

// RegisterDatabasePostgres Регистрация PostgreSQL базу данных
func (a *application) RegisterDatabasePostgres(fabric database.DbFactory, fn func(ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error)) {
	db := a.DbConnect(fabric, fn)

	a.dbPostgres = db
}

// RegisterDatabaseMongo Регистрация MongoDB базу данных
func (a *application) RegisterDatabaseMongo(fabric database.DbFactory, fn func(ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error)) {
	db := a.DbConnect(fabric, fn)

	a.dbMongo = db
}

func (a *application) DbConnect(fabric database.DbFactory, fn func(ctx *context.Context, opt *config.Configuration) (drivers.DbInterfase, error)) drivers.DbInterfase {
	db, err := fn(&a.ctx, a.opt)
	if err != nil {
		log.Fatalf("[ERROR] Connect %s: %s", fabric, err.Error())
		log.Println()
	}

	err = db.Connect()
	if err != nil {
		log.Fatalf("[ERROR] Connect %s: %s", fabric, err.Error())
		log.Println()
	}

	return db
}

// Start Старт приложения, запуск серверов
func (a *application) Start(errChan chan error) {
	for _, server := range a.servers {
		s := server

		a.wg.Add(1)

		go func() {
			if err := s.Start(); err != nil {
				errChan <- err
			}
		}()
	}
}
