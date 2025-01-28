package servers

import (
	"chat/config"
	"chat/database/drivers"
	"chat/managers"
	redisCli "chat/servers/redis"
	"context"
	"github.com/segmentio/kafka-go"
	protobufF3 "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"sync"
)

const MAX_SIZE = 10 * 1024 * 1024

type HttpServer struct {
	handler       http.Handler
	server        *http.Server
	ctx           context.Context
	wg            *sync.WaitGroup
	dbPostgres    drivers.DbInterfase
	dbMongo       drivers.DbInterfase
	redisCli      *redisCli.Redis
	userGrpc      protobuf.UsersClient
	f3GrpcClient  protobufF3.FileClientClient
	opt           *config.Configuration
	manChat       *managers.Chat
	kafkaProducer *kafka.Writer
	kafkaConsumer *kafka.Reader
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}

func NewHTTPServer(a *application) *HttpServer {
	srv := &HttpServer{
		ctx:           a.ctx,
		wg:            a.wg,
		redisCli:      &a.redisCli,
		dbPostgres:    a.dbPostgres,
		dbMongo:       a.dbMongo,
		opt:           a.opt,
		kafkaProducer: a.kafkaProducer,
		kafkaConsumer: a.kafkaConsumer,
		InfoLog:       a.InfoLog,
		ErrorLog:      a.ErrorLog,
	}

	err := srv.userServiceGrpc()
	if err != nil {
		log.Println(err)
	}

	err = srv.f3ServiceGrpc()
	if err != nil {
		log.Println(err)
	}

	return srv
}

func (h *HttpServer) userServiceGrpc() error {
	conn, err := grpc.NewClient(h.opt.ServerGrpcListenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := protobuf.NewUsersClient(conn)

	h.userGrpc = client

	return nil
}

func (h *HttpServer) f3ServiceGrpc() error {
	opts := make([]grpc.DialOption, 0, 2)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_SIZE), grpc.MaxCallSendMsgSize(MAX_SIZE)))

	conn, err := grpc.NewClient(h.opt.F3GrpcListenAddr, opts...)
	if err != nil {
		return err
	}

	client := protobufF3.NewFileClientClient(conn)

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

	manChat := managers.NewChat(h.dbMongo, h.redisCli, h.userGrpc)
	h.manChat = manChat

	go h.StartWS()

	go h.startConsumer()

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
