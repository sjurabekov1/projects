package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
	l "github.com/template-service/pkg/logger"
	"github.com/template-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//PostService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
}

//NewPostService ...
func NewPostService(db *sqlx.DB, log l.Logger) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid for new post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while generating uuid")
	}
	req.Id = id.String()
	user, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("failed while inserting post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while inserting post")
	}

	return user, nil
}
func (s *PostService) GetPostById(ctx context.Context, req *pb.GetByPostIdRequest) (*pb.Post, error) {
	post, err := s.storage.Post().GetPostById(req.UserId)
	if err != nil {
		s.logger.Error("failed get post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get user")
	}

	return post, err
}

func (s *PostService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	posts, err := s.storage.Post().GetAllUserPosts(req.UserId)
	if err != nil {
		s.logger.Error("failed get all user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get all user posts")
	}

	return &pb.GetUserPosts{
		Posts: posts,
	}, err
}

func (s *PostService) GetAllPostList(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	posts, count, err := s.storage.Post().GetAllPostList(req.Limit, req.Page)
	if err != nil {
		s.logger.Error("failed getting all user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed getting all user posts")
	}

	return &pb.GetPostResponse{
		Posts: posts,
		Count: count,
	}, err
}

