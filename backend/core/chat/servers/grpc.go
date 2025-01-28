package servers

import (
	"chat/config"
	"chat/database/drivers"
	userGrpc "chat/servers/grpc"
	"context"
	grpcResource "gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type GrpcServer struct {
	handler    func(*grpc.Server, drivers.DbInterfase)
	server     *grpc.Server
	ctx        context.Context
	wg         *sync.WaitGroup
	usersGRPC  grpcResource.UsersServer
	dbPostgres drivers.DbInterfase
	dbMongo    drivers.DbInterfase
	opt        *config.Configuration
	grpcResource.UnimplementedUsersServer
}

func NewGrpcServer(handler func(server *grpc.Server, db drivers.DbInterfase), a *application) *GrpcServer {
	return &GrpcServer{
		handler:    handler,
		ctx:        a.ctx,
		wg:         a.wg,
		dbPostgres: a.dbPostgres,
		dbMongo:    a.dbMongo,
		opt:        a.opt,
	}
}

func GrpcRegister(server *grpc.Server, db drivers.DbInterfase) {
	srv := userGrpc.NewGrpcUserServer(db)
	grpcResource.RegisterUsersServer(server, srv)
}

func (g *GrpcServer) Start() error {
	l, err := net.Listen("tcp", g.opt.GrpcListenAddr)
	if err != nil {
		return err
	}

	g.server = grpc.NewServer()
	g.handler(g.server, g.dbPostgres)

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
