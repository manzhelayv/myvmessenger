package servers

import (
	"context"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"server/config"
	"server/database/drivers"
	redisCli "server/servers/redis"
	"sync"
)

const MAX_SIZE = 10 * 1024 * 1024

type HttpServer struct {
	handler      http.Handler
	server       *http.Server
	ctx          context.Context
	wg           *sync.WaitGroup
	dbPostgres   drivers.DbInterfase
	dbMongo      drivers.DbInterfase
	redisCli     *redisCli.Redis
	f3GrpcClient protobuf.FileClientClient
	opt          *config.Configuration
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
}

func NewHTTPServer(a *application) *HttpServer {
	srv := &HttpServer{
		ctx:        a.ctx,
		wg:         a.wg,
		redisCli:   &a.redisCli,
		dbPostgres: a.dbPostgres,
		dbMongo:    a.dbMongo,
		opt:        a.opt,
		InfoLog:    a.InfoLog,
		ErrorLog:   a.errorLog,
	}

	err := srv.userServiceGrpc()
	if err != nil {
		log.Println(err)
	}

	return srv
}

func (h *HttpServer) userServiceGrpc() error {
	opts := make([]grpc.DialOption, 0, 2)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_SIZE), grpc.MaxCallSendMsgSize(MAX_SIZE)))

	conn, err := grpc.NewClient(h.opt.F3GrpcListenAddr, opts...)
	if err != nil {
		return err
	}

	client := protobuf.NewFileClientClient(conn)

	h.f3GrpcClient = client

	return nil
}

func (h *HttpServer) Start() error {
	go h.gracefulShutdown()

	h.handler = h.setupRouter()

	h.server = &http.Server{
		Addr:    h.opt.ListenAddr,
		Handler: h.handler,
	}

	log.Println("[INFO] Server HTTP started port", h.opt.ListenAddr)

	return h.server.ListenAndServe()
}

func (h *HttpServer) Stop() error {
	return h.server.Shutdown(h.ctx)
}

func (h *HttpServer) gracefulShutdown() {
	defer h.wg.Done()
	<-h.ctx.Done()
	log.Println("Shutting down HTTP server")

	grCtx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	if err := h.server.Shutdown(grCtx); err != nil {
		log.Println(err, "[ERROR] HTTP server Shutdown")
	}
}
