package servers

import (
	"context"
	"errors"
	"f3/config"
	"f3/database/drivers/file"
	redisCli "f3/servers/redis"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/minio/minio-go/v7"
	"sync"

	"log"
)

const MINIO = "minio"

type MinioServer struct {
	client   *minio.Client
	ctx      context.Context
	wg       *sync.WaitGroup
	redisCli *redisCli.Redis
	opt      *config.Configuration
	config   *redis.Options
}

func NewMinioServer(a *application) *MinioServer {
	srv := &MinioServer{
		ctx:      a.ctx,
		wg:       a.wg,
		redisCli: &a.redisCli,
		opt:      a.opt,
	}

	return srv
}

func (m *MinioServer) Start() error {
	ds := file.New(
		m.opt.MinioApiAddr,
		m.opt.MinioAccessKeyID,
		m.opt.MinioSecretAccessKey,
		m.opt.MinioLocation,
		file.InitBuckets(m.opt.MinioDefaultBucketName, map[string]string{MINIO: m.opt.MinioDefaultBucketName}),
	)
	if err := ds.Connect(); err != nil {
		errMsg := fmt.Sprintf("[ERROR] Cannot connect to database Minio S3 Storage: %v", err)
		return errors.New(errMsg)
	}

	log.Println("[INFO] Connected to File Minio Storage")

	return nil
}

func (m *MinioServer) Stop() error {
	//m.client.Shutdown()

	return nil
}

func (m *MinioServer) gracefulShutdown() {
	defer m.wg.Done()
	<-m.ctx.Done()
	log.Println("Shutting down Redis server")
}
