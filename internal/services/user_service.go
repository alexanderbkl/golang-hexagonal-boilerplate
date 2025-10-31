package services

import (
	"context"
	"time"

	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/domain"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/ports"
	"github.com/google/uuid"
)

// UserService implements the UserService interface
type UserService struct {
	repo  ports.UserRepository
	cache ports.CacheRepository
}

// NewUserService creates a new user service
func NewUserService(repo ports.UserRepository, cache ports.CacheRepository) ports.UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, input *domain.CreateUserInput) (*domain.User, error) {
	// Validate input
	if input.Email == "" || input.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create user
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

// ListUsers retrieves a list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.repo.List(ctx, limit, offset)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id string, input *domain.UpdateUserInput) (*domain.User, error) {
	// Check if user exists
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Update user
	updatedUser, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}
