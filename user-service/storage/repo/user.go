package repo

import (
	pb "github.com/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	UpdateUser(*pb.User) (string, error)
	GetUserById(id string) (*pb.User, error)
	GetUserList(limit, page int64) ([]*pb.User, int64, error)
}
