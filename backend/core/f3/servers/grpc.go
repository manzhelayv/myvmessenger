package servers

import (
	"context"
	"f3/config"
	"f3/database/drivers/file"
	"f3/manager"
	userGrpc "f3/servers/grpc"
	grpcResource "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

const MAX_SIZE = 10 * 1024 * 1024

type Option func(server *GrpcServer)

type GrpcServer struct {
	handler      func(*grpc.Server, manager.FileService)
	server       *grpc.Server
	ctx          context.Context
	wg           *sync.WaitGroup
	usersGRPC    grpcResource.FileClientServer
	filesStorage manager.FileService
	minio        *MinioServer
	opt          *config.Configuration
	grpcResource.UnimplementedFileClientServer
}

func NewGrpcServer(handler func(server *grpc.Server, db manager.FileService), a *application) *GrpcServer {
	return &GrpcServer{
		handler: handler,
		ctx:     a.ctx,
		wg:      a.wg,
		opt:     a.opt,
		minio:   &a.minio,
	}
}

func GrpcRegister(server *grpc.Server, db manager.FileService) {
	opt := config.Parse()

	ds := file.New(
		opt.MinioApiAddr,
		opt.MinioAccessKeyID,
		opt.MinioSecretAccessKey,
		opt.MinioLocation,
		file.InitBuckets(opt.MinioDefaultBucketName, map[string]string{MINIO: opt.MinioDefaultBucketName}),
	)
	if err := ds.Connect(); err != nil {
		log.Println(err)
	}

	filesStorage := WithFileStorageManager(manager.NewFileStorageManager(ds.FilesRepository()))

	srv := userGrpc.NewFileStorageResource(filesStorage)

	grpcResource.RegisterFileClientServer(server, srv)
}

func WithFileStorageManager(man manager.FileService) manager.FileService {
	//return func(s *GrpcServer) {
	//	s.filesStorage = man
	//}

	return man
}

func (g *GrpcServer) Start() error {
	l, err := net.Listen("tcp", g.opt.GrpcListenAddr)
	if err != nil {
		return err
	}

	opts := make([]grpc.ServerOption, 0, 2)

	opts = append(opts, grpc.MaxRecvMsgSize(MAX_SIZE))
	opts = append(opts, grpc.MaxSendMsgSize(MAX_SIZE))

	g.server = grpc.NewServer(opts...)
	g.handler(g.server, g.filesStorage)

	go g.gracefulShutdown()

	log.Println("[INFO] Server GRPC started port", g.opt.GrpcListenAddr)

	return g.server.Serve(l)
}

func (g *GrpcServer) Stop() error {
	g.server.GracefulStop()

	return nil
}

func (g *GrpcServer) gracefulShutdown() {
	defer g.wg.Done()
	<-g.ctx.Done()
	log.Println("Shutting down GRPC server")
	g.server.GracefulStop()
}
