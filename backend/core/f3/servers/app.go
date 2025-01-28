package servers

import (
	"context"
	"f3/config"
	"f3/manager"
	redisCli "f3/servers/redis"
	//minio "f3/servers/minio"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	once      sync.Once
	singleton *application
)

// Основная структура приложения
type application struct {
	servers  []Servers
	wg       *sync.WaitGroup
	ctx      context.Context
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	redisCli redisCli.Redis
	db       manager.FileService
	minio    MinioServer
	opt      *config.Configuration
}

func Instance(ctx context.Context, infoLog *log.Logger, errorLog *log.Logger, opt *config.Configuration) *application {
	once.Do(func() {
		singleton = &application{
			ctx:      ctx,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
			opt:      opt,
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
func (a *application) RegisterGrpcServer(fn func(server *grpc.Server, db manager.FileService)) {
	srv := NewGrpcServer(fn, a)

	a.servers = append(a.servers, srv)
}

// RegisterMinioServer Регистрация Minio сервера
func (a *application) RegisterMinioServer() {
	srv := NewMinioServer(a)

	a.minio = *srv

	a.servers = append(a.servers, srv)
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
