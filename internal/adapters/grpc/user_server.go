package grpc

import (
	"context"

	pb "github.com/alexanderbkl/golang-hexagonal-boilerplate/api/grpc"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/domain"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer implements the gRPC UserService server
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService ports.UserService
}

// NewUserServiceServer creates a new gRPC user service server
func NewUserServiceServer(userService ports.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// CreateUser creates a new user
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := s.userService.CreateUser(ctx, &domain.CreateUserInput{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.userService.GetUser(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// ListUsers retrieves a list of users
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)
	if limit == 0 {
		limit = 10
	}

	users, err := s.userService.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	grpcUsers := make([]*pb.User, len(users))
	for i, user := range users {
		grpcUsers[i] = &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &pb.ListUsersResponse{
		Users: grpcUsers,
	}, nil
}

// UpdateUser updates a user
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	input := &domain.UpdateUserInput{}
	if req.Email != nil {
		input.Email = req.Email
	}
	if req.Name != nil {
		input.Name = req.Name
	}

	user, err := s.userService.UpdateUser(ctx, req.Id, input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// DeleteUser deletes a user
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}
