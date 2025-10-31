package ports

import (
	"context"

	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/domain"
)

// UserService defines the business logic interface
type UserService interface {
	CreateUser(ctx context.Context, input *domain.CreateUserInput) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id string, input *domain.UpdateUserInput) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}
