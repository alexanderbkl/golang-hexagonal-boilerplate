package db

import (
	"context"
	"time"

	sqlcdb "github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/adapters/db/sqlc"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/domain"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements the UserRepository interface using PostgreSQL
type PostgresRepository struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) ports.UserRepository {
	return &PostgresRepository{
		db:      db,
		queries: sqlcdb.New(db),
	}
}

// Helper function to convert time.Time to pgtype.Timestamp
func toPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

// Helper function to convert pgtype.Timestamp to time.Time
func fromPgTimestamp(t pgtype.Timestamp) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

// Create creates a new user
func (r *PostgresRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.queries.CreateUser(ctx, sqlcdb.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: toPgTimestamp(user.CreatedAt),
		UpdatedAt: toPgTimestamp(user.UpdatedAt),
	})
	return err
}

// GetByID retrieves a user by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: fromPgTimestamp(user.CreatedAt),
		UpdatedAt: fromPgTimestamp(user.UpdatedAt),
	}, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: fromPgTimestamp(user.CreatedAt),
		UpdatedAt: fromPgTimestamp(user.UpdatedAt),
	}, nil
}

// List retrieves a list of users
func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := r.queries.ListUsers(ctx, sqlcdb.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = &domain.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: fromPgTimestamp(u.CreatedAt),
			UpdatedAt: fromPgTimestamp(u.UpdatedAt),
		}
	}

	return result, nil
}

// Update updates a user
func (r *PostgresRepository) Update(ctx context.Context, id string, input *domain.UpdateUserInput) (*domain.User, error) {
	// Get current user
	currentUser, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Use current values if not provided in input
	email := currentUser.Email
	name := currentUser.Name

	if input.Email != nil {
		email = *input.Email
	}
	if input.Name != nil {
		name = *input.Name
	}

	user, err := r.queries.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:        id,
		Email:     email,
		Name:      name,
		UpdatedAt: toPgTimestamp(time.Now()),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: fromPgTimestamp(user.CreatedAt),
		UpdatedAt: fromPgTimestamp(user.UpdatedAt),
	}, nil
}

// Delete deletes a user
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteUser(ctx, id)
}
