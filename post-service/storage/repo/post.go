package repo

import (
	pb "github.com/template-service/genproto"
)

//PostStorageI ...
type PostStorageI interface {
	CreatePost(*pb.Post) (*pb.Post, error)
	GetPostById(id string) (*pb.Post, error)
	GetAllUserPosts(userID string) ([]*pb.Post, error)
	GetAllPostList(limit, page int64) ([]*pb.Post, int64, error)
}
