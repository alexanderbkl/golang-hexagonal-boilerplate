# Contributing to Golang Hexagonal Boilerplate

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL 16 (for local development)
- Redis 7 (for local development)
- Protocol Buffer compiler (protoc)
- Make

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/alexanderbkl/golang-hexagonal-boilerplate.git
cd golang-hexagonal-boilerplate
```

2. Install development tools:
```bash
make install-tools
```

3. Install dependencies:
```bash
make deps
```

4. Copy environment variables:
```bash
cp .env.example .env
```

5. Start dependencies:
```bash
docker-compose up -d postgres redis
```

6. Run migrations:
```bash
make migrate-up
```

7. Run the application:
```bash
make run
```

## Project Structure

See [ARCHITECTURE.md](ARCHITECTURE.md) for a detailed explanation of the project structure and architecture.

## Making Changes

### Code Style

- Follow Go best practices and conventions
- Run `go fmt` before committing
- Run `go vet` to check for common issues
- Keep functions small and focused
- Write self-documenting code with clear variable names

### Adding New Features

When adding new features, follow the hexagonal architecture pattern:

1. **Define Domain Entities** in `internal/domain/`
   - Add business entities
   - Define value objects
   - Add domain-specific errors

2. **Define Ports** in `internal/ports/`
   - Create service interfaces
   - Create repository interfaces
   - Define other adapter interfaces

3. **Implement Services** in `internal/services/`
   - Implement business logic
   - Keep services focused on use cases
   - Depend only on interfaces (ports)

4. **Implement Adapters** in `internal/adapters/`
   - Add database repositories
   - Add API handlers (GraphQL/gRPC)
   - Add external service clients

### Database Changes

1. Create migration files in `migrations/`:
```bash
# Format: XXX_description.up.sql and XXX_description.down.sql
migrations/002_add_feature.up.sql
migrations/002_add_feature.down.sql
```

2. Add SQL queries in `db/queries/`:
```sql
-- name: GetFeatureByID :one
SELECT * FROM features WHERE id = $1;
```

3. Regenerate sqlc code:
```bash
make sqlc-generate
```

4. Implement repository methods in `internal/adapters/db/`

### GraphQL Changes

1. Update schema in `api/graphql/schema.graphql`

2. Regenerate GraphQL code:
```bash
make graphql-generate
```

3. Implement resolvers in `internal/adapters/graphql/schema.resolvers.go`

### gRPC Changes

1. Update proto files in `api/grpc/`

2. Regenerate gRPC code:
```bash
make grpc-generate
```

3. Implement service methods in `internal/adapters/grpc/`

## Testing

### Running Tests

```bash
make test
```

### Writing Tests

- Write unit tests for business logic in services
- Write integration tests for repositories
- Mock dependencies using interfaces
- Aim for high test coverage on critical paths

Example test structure:
```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    mockCache := &MockCacheRepository{}
    service := services.NewUserService(mockRepo, mockCache)
    
    // Test
    user, err := service.CreateUser(context.Background(), &domain.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

## Code Generation

This project uses several code generators:

- **sqlc**: Generates type-safe Go code from SQL
- **gqlgen**: Generates GraphQL server code
- **protoc**: Generates gRPC code from Protocol Buffers

After modifying schemas or queries, regenerate code:
```bash
make sqlc-generate      # After SQL query changes
make graphql-generate   # After GraphQL schema changes
make grpc-generate      # After protobuf changes
```

## Building

### Local Build

```bash
make build
```

### Docker Build

```bash
make docker-build
```

## Commit Guidelines

- Write clear, concise commit messages
- Use present tense ("Add feature" not "Added feature")
- Reference issue numbers when applicable
- Keep commits focused on a single change

Example commit message:
```
Add user profile endpoint to GraphQL API

- Extend User type with profile fields
- Add GetUserProfile query
- Implement resolver logic
- Add database migration for profile table

Closes #123
```

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes following the guidelines above
3. Ensure all tests pass
4. Update documentation if needed
5. Submit a pull request with a clear description

## Questions?

Feel free to open an issue for questions, bugs, or feature requests.
