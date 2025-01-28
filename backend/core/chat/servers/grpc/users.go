package grpc

import (
	"chat/database/drivers"
	grpcResource "gitlab.com/myvmessenger/client/server-client/user/protobuf"
)

type UserServer struct {
	db drivers.DbInterfase
	grpcResource.UnimplementedUsersServer
}

func NewGrpcUserServer(db drivers.DbInterfase) *UserServer {
	return &UserServer{
		db: db,
	}
}
