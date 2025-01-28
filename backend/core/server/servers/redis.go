package servers

import (
	"context"
	"github.com/go-redis/redis/v7"
	"server/config"
	redisCli "server/servers/redis"
	"sync"

	"log"
)

type RedisServer struct {
	client   redis.Client
	ctx      context.Context
	wg       *sync.WaitGroup
	redisCli *redisCli.Redis
	opt      *config.Configuration
	config   *redis.Options
}

func NewRedisServer(ctx context.Context, wg *sync.WaitGroup, a *application) *RedisServer {
	config := &redis.Options{
		Addr:     a.opt.REDISURL,
		Password: a.opt.REDISPASSWORD,
		DB:       a.opt.REDISDB,
	}

	cli := &RedisServer{
		ctx:      ctx,
		wg:       wg,
		redisCli: &a.redisCli,
		opt:      a.opt,
		config:   config,
	}

	return cli
}

func (r *RedisServer) Start() error {
	client := redis.NewClient(r.config)

	if _, err := client.Ping().Result(); err != nil {
		log.Fatalf("Connect to redis client failed, err: %v\n", err)
	}

	r.client = *client

	if _, err := r.client.Ping().Result(); err != nil {
		log.Fatalf("Connect to redis client failed, err: %v\n", err)
	}

	go r.gracefulShutdown()

	res := redisCli.NewRedis(client)
	*r.redisCli = *res

	log.Println("[INFO] Server Redis started port", r.opt.REDISURL)

	return nil
}

func (r *RedisServer) Stop() error {
	r.client.Shutdown()

	return nil
}

func (r *RedisServer) gracefulShutdown() {
	defer r.wg.Done()
	<-r.ctx.Done()
	log.Println("Shutting down Redis server")
	r.client.Shutdown()
}
