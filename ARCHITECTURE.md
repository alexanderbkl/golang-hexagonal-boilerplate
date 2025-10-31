# Project Structure

This document provides an overview of the complete project structure and architecture.

## Directory Layout

```
.
├── api/                          # API definitions and contracts
│   ├── graphql/                  # GraphQL schema definitions
│   │   └── schema.graphql        # User GraphQL schema
│   └── grpc/                     # gRPC protocol buffer definitions
│       ├── user.proto            # User service protobuf definition
│       ├── user.pb.go            # Generated Go protobuf code
│       └── user_grpc.pb.go       # Generated gRPC server/client code
│
├── cmd/                          # Application entry points
│   ├── server/                   # Main application server
│   │   └── main.go              # Server startup and initialization
│   └── migrate/                  # Database migration tool
│       └── main.go              # Migration runner
│
├── internal/                     # Private application code
│   ├── domain/                   # Business domain layer
│   │   ├── user.go              # User entity and value objects
│   │   ├── errors.go            # Domain-specific errors
│   │   └── ...                  # Other domain entities
│   │
│   ├── ports/                    # Interface definitions (hexagonal ports)
│   │   ├── service.go           # Business logic interfaces
│   │   ├── repository.go        # Data persistence interfaces
│   │   └── cache.go             # Caching interfaces
│   │
│   ├── services/                 # Business logic implementation
│   │   └── user_service.go      # User service implementation
│   │
│   └── adapters/                 # External service adapters
│       ├── db/                   # Database adapter (PostgreSQL)
│       │   ├── postgres.go      # Repository implementation
│       │   └── sqlc/            # Generated sqlc code
│       │       ├── db.go
│       │       ├── models.go
│       │       ├── querier.go
│       │       └── users.sql.go
│       │
│       ├── graphql/              # GraphQL adapter
│       │   ├── resolver.go      # GraphQL resolver setup
│       │   ├── schema.resolvers.go  # GraphQL resolver implementations
│       │   ├── generated.go     # Generated GraphQL code
│       │   └── models_gen.go    # Generated GraphQL models
│       │
│       ├── grpc/                 # gRPC adapter
│       │   └── user_server.go   # gRPC service implementation
│       │
│       └── redis/                # Redis cache adapter
│           └── redis.go         # Cache implementation
│
├── pkg/                          # Public/shared packages
│   ├── config/                   # Configuration management
│   │   └── config.go            # Environment-based configuration
│   └── logger/                   # Logging utilities
│       └── logger.go            # Simple logger implementation
│
├── migrations/                   # Database migrations
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
│
├── db/queries/                   # SQL queries for sqlc
│   └── users.sql                # User CRUD queries
│
├── Dockerfile                    # Docker image definition
├── docker-compose.yml            # Multi-container setup
├── Makefile                      # Development tasks automation
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── sqlc.yaml                     # sqlc configuration
├── gqlgen.yml                    # GraphQL generator configuration
├── .env.example                  # Example environment variables
├── .gitignore                    # Git ignore rules
└── README.md                     # Project documentation
```

## Hexagonal Architecture Layers

### 1. Domain Layer (`internal/domain/`)
- **Purpose**: Core business logic and entities
- **Dependencies**: None (pure business logic)
- **Contents**:
  - Business entities (User)
  - Value objects
  - Domain errors
  - Business rules

### 2. Ports Layer (`internal/ports/`)
- **Purpose**: Define interfaces between layers
- **Dependencies**: Domain layer only
- **Contents**:
  - Service interfaces (business logic contracts)
  - Repository interfaces (data persistence contracts)
  - Cache interfaces (caching contracts)

### 3. Services Layer (`internal/services/`)
- **Purpose**: Implement business logic
- **Dependencies**: Domain and Ports layers
- **Contents**:
  - Service implementations
  - Business logic orchestration
  - Use case implementations

### 4. Adapters Layer (`internal/adapters/`)
- **Purpose**: Implement external integrations
- **Dependencies**: Domain, Ports, and Services layers
- **Contents**:
  - Database adapters (PostgreSQL with sqlc)
  - API adapters (GraphQL, gRPC)
  - Cache adapters (Redis)
  - External service clients

### 5. Infrastructure Layer (`pkg/`, `cmd/`)
- **Purpose**: Application startup and shared utilities
- **Dependencies**: All layers
- **Contents**:
  - Configuration management
  - Logging utilities
  - Application entry points

## Data Flow

### GraphQL Request Flow
```
GraphQL Client → GraphQL Playground/Query
    ↓
GraphQL Handler (adapter)
    ↓
User Service (business logic)
    ↓
Repository Interface (port)
    ↓
PostgreSQL Repository (adapter)
    ↓
Database
```

### gRPC Request Flow
```
gRPC Client
    ↓
gRPC Server (adapter)
    ↓
User Service (business logic)
    ↓
Repository Interface (port)
    ↓
PostgreSQL Repository (adapter)
    ↓
Database
```

## Key Design Decisions

1. **Hexagonal Architecture**: Ensures clean separation of concerns and testability
2. **Interface-based Design**: All dependencies are through interfaces for flexibility
3. **sqlc**: Type-safe SQL with compile-time query validation
4. **gqlgen**: Code-first GraphQL with strong typing
5. **Protocol Buffers**: Efficient serialization for gRPC
6. **Docker Compose**: Easy local development and deployment
7. **Environment Configuration**: 12-factor app configuration management

## Running the Application

See [README.md](README.md) for detailed instructions on:
- Quick start with Docker
- Local development setup
- API usage examples
- Development tasks
