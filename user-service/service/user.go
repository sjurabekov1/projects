package service

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
	l "github.com/template-service/pkg/logger"
	cl "github.com/template-service/service/grpc_client"
	"github.com/template-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid for user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while generating uuid for user")
	}
	req.Id = id.String()
	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error("failed while creating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while creating user")
	}
	if req.Posts != nil {
		for _, post := range req.Posts {
			post.UserId = req.Id
			createdPosts, err := s.client.PostService().CreatePost(context.Background(), post)
			if err != nil {
				s.logger.Error("failed while inserting user post", l.Error(err))
				return nil, status.Error(codes.Internal, "failed while inserting user post")
			}
			fmt.Println(createdPosts)
		}
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.UpdateUserResponse, error) {
	id, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error("failed while updating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while updating user")
	}
	return &pb.UpdateUserResponse{
		Id: id,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.UserId)
	if err != nil {
		s.logger.Error("failed while getting by Id user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting by Id user")
	}

	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}

	//	user.Posts = posts.Posts
	return user, err
}

// func (s *UserService) GetAllUser(ctx context.Context, req *pb.Empty) (*pb.GetAllResponse, error) {
// 	users, err := s.storage.User().GetAllUser()
// 	if err != nil {
// 		s.logger.Error("failed while getting All users", l.Error(err))
// 		return nil, status.Error(codes.Internal, "failed while getting All users")
// 	}

// 	for _, user := range users {
// 		posts, err := s.client.PostService().GetAllUserPosts(
// 			ctx,
// 			&pb.GetUserPostsrequest{
// 				UserId: user.Id,
// 			},
// 		)
// 		if err != nil {
// 			s.logger.Error("failed while getting user posts", l.Error(err))
// 			return nil, status.Error(codes.Internal, "failed while getting user posts")
// 		}

// 		user.Posts = posts.Posts
// 	}

// 	return &pb.GetAllResponse{
// 		Users: users,
// 	}, err
// }

func (s *UserService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	res, err := s.client.PostService().GetAllUserPosts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, count, err := s.storage.User().GetUserList(req.Limit, req.Page)
	if err != nil {
		s.logger.Error("failed while getting all users", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting all users")
	}

	for _, user := range users {
		post, err := s.client.PostService().GetAllUserPosts(ctx, &pb.GetUserPostsrequest{UserId: user.Id})
		if err != nil {
			s.logger.Error("failed while getting all user posts", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while getting all posts")
		}

		user.Posts = post.Posts
	}

	return &pb.GetUsersResponse{
		Users: users,
		Count: count,
	}, nil
}
