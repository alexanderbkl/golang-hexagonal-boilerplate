package domain

import (
	"time"
)

// User represents the core user entity in the domain
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}
