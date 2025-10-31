package ports

import (
	"context"

	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/domain"
)

// UserRepository defines the interface for user storage operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Update(ctx context.Context, id string, input *domain.UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}
