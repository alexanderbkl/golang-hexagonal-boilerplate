# Golang Hexagonal Boilerplate

A production-ready Go boilerplate implementing hexagonal (ports and adapters) architecture with GraphQL, gRPC, PostgreSQL, Redis, sqlc, and Docker support.

## Features

- ✅ **Hexagonal Architecture** (Ports & Adapters pattern)
- ✅ **GraphQL API** with `gqlgen`
- ✅ **gRPC API** with Protocol Buffers
- ✅ **PostgreSQL** database with connection pooling
- ✅ **sqlc** for type-safe SQL queries
- ✅ **Redis** for caching
- ✅ **Docker & Docker Compose** for containerization
- ✅ **Database migrations** with `golang-migrate`
- ✅ **Clean separation** of concerns (domain, ports, adapters)

## Architecture

```
.
├── api/                    # API definitions
│   ├── graphql/           # GraphQL schemas
│   └── grpc/              # Protocol Buffer definitions
├── cmd/                   # Application entrypoints
│   ├── server/            # Main server
│   └── migrate/           # Migration tool
├── internal/              # Private application code
│   ├── domain/            # Business logic and entities
│   ├── ports/             # Interface definitions
│   └── adapters/          # External service implementations
│       ├── db/            # PostgreSQL adapter
│       ├── graphql/       # GraphQL adapter
│       ├── grpc/          # gRPC adapter
│       └── redis/         # Redis adapter
├── pkg/                   # Public libraries
│   ├── config/            # Configuration management
│   └── logger/            # Logging utilities
├── migrations/            # Database migrations
└── db/queries/            # SQL queries for sqlc
```

### Hexagonal Architecture

This project follows the hexagonal architecture pattern:

- **Domain Layer**: Contains business logic, entities, and domain services
- **Ports**: Define interfaces for communication between layers
  - Primary ports: Interfaces for driving adapters (GraphQL, gRPC)
  - Secondary ports: Interfaces for driven adapters (Database, Cache)
- **Adapters**: Implement the ports
  - Primary adapters: GraphQL and gRPC handlers
  - Secondary adapters: PostgreSQL and Redis implementations

## Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL 16 (if running locally)
- Redis 7 (if running locally)
- Protocol Buffer compiler (protoc)
- Make

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/alexanderbkl/golang-hexagonal-boilerplate.git
cd golang-hexagonal-boilerplate
```

2. Start all services:
```bash
make docker-up
```

This will:
- Start PostgreSQL database
- Start Redis cache
- Run database migrations
- Start the application

3. Access the services:
- GraphQL Playground: http://localhost:8080
- gRPC Server: localhost:9090

4. Stop the services:
```bash
make docker-down
```

### Local Development

1. Install dependencies:
```bash
make deps
make install-tools
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start PostgreSQL and Redis:
```bash
docker-compose up -d postgres redis
```

4. Run migrations:
```bash
make migrate-up
```

5. Run the application:
```bash
make run
```

## API Usage

### GraphQL

The GraphQL playground is available at http://localhost:8080

**Create a user:**
```graphql
mutation {
  createUser(input: {
    email: "user@example.com"
    name: "John Doe"
  }) {
    id
    email
    name
    createdAt
    updatedAt
  }
}
```

**Get a user:**
```graphql
query {
  user(id: "user-id") {
    id
    email
    name
    createdAt
    updatedAt
  }
}
```

**List users:**
```graphql
query {
  users(limit: 10, offset: 0) {
    id
    email
    name
    createdAt
    updatedAt
  }
}
```

**Update a user:**
```graphql
mutation {
  updateUser(id: "user-id", input: {
    name: "Jane Doe"
  }) {
    id
    email
    name
    updatedAt
  }
}
```

**Delete a user:**
```graphql
mutation {
  deleteUser(id: "user-id")
}
```

### gRPC

The gRPC server runs on port 9090. You can use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) or [BloomRPC](https://github.com/bloomrpc/bloomrpc) to interact with it.

**Example using grpcurl:**

```bash
# List services
grpcurl -plaintext localhost:9090 list

# Create a user
grpcurl -plaintext -d '{"email": "user@example.com", "name": "John Doe"}' \
  localhost:9090 user.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{"id": "user-id"}' \
  localhost:9090 user.UserService/GetUser

# List users
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:9090 user.UserService/ListUsers
```

## Development

### Database Migrations

Create a new migration:
```bash
# Manually create files in migrations/ directory
# Format: XXX_description.up.sql and XXX_description.down.sql
```

Apply migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

### Code Generation

**Generate sqlc code** (after modifying SQL queries):
```bash
make sqlc-generate
```

**Generate GraphQL code** (after modifying schema):
```bash
make graphql-generate
```

**Generate gRPC code** (after modifying proto files):
```bash
make grpc-generate
```

### Building

Build the application:
```bash
make build
```

Build Docker image:
```bash
make docker-build
```

### Testing

Run tests:
```bash
make test
```

## Project Structure Details

### Domain Layer

- `internal/domain/user.go`: User entity and value objects
- `internal/domain/user_service.go`: Business logic implementation
- `internal/domain/errors.go`: Domain-specific errors

### Ports Layer

- `internal/ports/service.go`: Business logic interface
- `internal/ports/repository.go`: Data persistence interface
- `internal/ports/cache.go`: Caching interface

### Adapters Layer

- `internal/adapters/db/`: PostgreSQL implementation using sqlc
- `internal/adapters/graphql/`: GraphQL resolvers
- `internal/adapters/grpc/`: gRPC service implementation
- `internal/adapters/redis/`: Redis cache implementation

## Environment Variables

See `.env.example` for all available configuration options.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

This boilerplate demonstrates a clean architecture approach for Go applications, suitable for microservices and complex business logic applications.

