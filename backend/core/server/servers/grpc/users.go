package grpc

import (
	"context"
	grpcResource "gitlab.com/myvmessenger/client/server-client/user/protobuf"
	"server/database/drivers"
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

// GetUserTdid Получение tdid пользователя по email
func (u *UserServer) GetUserTdid(ctx context.Context, grpc *grpcResource.User) (*grpcResource.User, error) {
	email := grpc.Email

	user, err := u.db.GetUserTdidForEmailOrLogin(ctx, email)
	if err != nil {
		return nil, err
	}

	return &grpcResource.User{
		Tdid: user.Tdid,
	}, nil
}

// GetUsersFromTdid Получение списка пользователей по массиву tdid
func (u *UserServer) GetUsersFromTdid(ctx context.Context, tdids *grpcResource.UserTdids) (*grpcResource.UsersItem, error) {
	users, err := u.db.GetUsersFromTdid(ctx, tdids.Tdid)
	if err != nil {
		return nil, err
	}

	usersProto := &grpcResource.UsersItem{}
	for _, user := range users {
		usersProto.Users = append(usersProto.Users, &grpcResource.UserItem{
			Email:  user.Email,
			Tdid:   user.Tdid,
			Name:   user.Name,
			Login:  user.Login,
			Avatar: user.Image,
		})
	}

	return usersProto, nil
}
