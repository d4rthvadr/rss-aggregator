# RSS Aggregator

A Go-based RSS feed aggregator built with Chi router and PostgreSQL. This application allows users to manage RSS feeds, follow feeds, and aggregate content from multiple sources.

## Features

- ✅ User management with auto-generated API keys
- ✅ API key-based authentication
- ✅ Create and manage RSS feeds
- ✅ Follow/unfollow RSS feeds
- ✅ Thread-safe request handling with context-based authentication
- ✅ Type-safe database queries with sqlc
- ✅ Database migrations with goose
- ✅ Clean architecture with domain/database model separation
- ✅ CORS support for web clients
- ✅ PostgreSQL with UUID primary keys and timezone-aware timestamps

## Technology Stack

### Core

- **Go 1.23.2**: Primary programming language
- **PostgreSQL 15**: Database with UUID and timezone support
- **Chi Router v5**: Lightweight HTTP router and middleware framework
- **sqlc**: Type-safe SQL code generation (compile-time SQL validation)
- **goose**: Database migration management

### Libraries

- `github.com/go-chi/chi/v5` - HTTP routing and middleware
- `github.com/go-chi/cors` - CORS middleware
- `github.com/google/uuid` - UUID generation
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/joho/godotenv` - Environment variable management

### Development Tools

- **goenv** - Go version management
- **Docker Compose** - Local database containerization

## Prerequisites

- Go 1.23+ (managed via goenv)
- Docker and Docker Compose
- Git
- goose (database migrations)
- sqlc (SQL code generation)

## System Requirements & Versions

| Tool           | Version | Purpose                       |
| -------------- | ------- | ----------------------------- |
| Go             | 1.23.2+ | Application runtime           |
| PostgreSQL     | 15      | Database                      |
| goenv          | 2.2.30+ | Go version management         |
| goose          | latest  | Database migrations           |
| sqlc           | latest  | Type-safe SQL code generation |
| Docker         | latest  | Container runtime             |
| Docker Compose | latest  | Multi-container orchestration |

## Project Structure

```
rss-aggregator/
├── docker-compose.yml           # PostgreSQL database setup
├── .env                         # Environment variables
├── Makefile                     # Goose migration shortcuts
├── sqlc.yaml                    # sqlc configuration
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── main.go                      # Main application entry point
├── json.go                      # JSON response utilities
├── models.go                    # Domain models (User, Feed, FeedFollow)
├── middleware_auth.go           # Authentication middleware
├── handler_readiness.go         # Health check handler
├── handler_error.go             # Error handler
├── handler_user.go              # User CRUD handlers
├── handler_feed.go              # Feed CRUD handlers
├── handler_feed_follows.go      # Feed follows handlers
├── internal/
│   ├── auth/
│   │   └── auth.go             # API key extraction utilities
│   └── database/
│       ├── db.go               # Generated database connection
│       ├── models.go           # Generated database models
│       ├── users.sql.go        # Generated user queries
│       ├── feeds.sql.go        # Generated feed queries
│       └── feed_follows.sql.go # Generated feed_follows queries
├── sql/
│   ├── schema/                 # Database migrations
│   │   ├── 001_users.sql       # Create users table
│   │   ├── 002_users.sql       # Add api_key to users
│   │   ├── 003_feeds.sql       # Create feeds table
│   │   └── 004_feed_follows.sql # Create feed_follows table
│   └── queries/                # SQL queries for sqlc
│       ├── users.sql           # User queries
│       ├── feeds.sql           # Feed queries
│       └── feed_follows.sql    # Feed follow queries
└── README.md
```

## Installation

### 1. Install Go Version Manager (goenv)

```bash
# Install goenv via Homebrew
arch -arm64 brew install goenv

# Add to your ~/.bash_profile or ~/.zshrc
echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile
echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
echo 'eval "$(goenv init -)"' >> ~/.bash_profile
echo 'export PATH="$HOME/.goenv/shims:$PATH"' >> ~/.bash_profile
source ~/.bash_profile

# Install Go 1.23.2
goenv install 1.23.2
goenv global 1.23.2
```

### 2. Install Development Tools

```bash
# Install goose for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc for type-safe SQL code generation
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installations
goose --version
sqlc version
```

## Getting Started

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd rss-aggregator
```

### 2. Start the Database

```bash
# Start PostgreSQL container
docker-compose up -d

# Verify database is running
docker-compose logs postgres
```

### 3. Configure Environment Variables

Create a `.env` file in the project root:

```env
# Application
APP_PORT=8080

# Database Connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rss_aggregator
DB_SSLMODE=disable

# Database URL (used by goose)
DB_URL=postgresql://postgres:password@localhost:5432/rss_aggregator?sslmode=disable
```

### 4. Run Database Migrations

```bash
# Apply all migrations
make migrate-up

# Check migration status
make migrate-status
```

### 5. Generate Database Code

```bash
# Generate type-safe Go code from SQL queries
sqlc generate
```

### 6. Install Go Dependencies

```bash
go mod tidy
```

### 7. Run the Application

```bash
# Build and run
go build && ./rss-aggregator

# Or run directly
go run .
```

## Architecture & Design

### Clean Architecture Pattern

The application follows a clean architecture with clear separation of concerns:

```
HTTP Layer (handlers)
    ↓
Domain Layer (models.go)
    ↓
Database Layer (internal/database)
```

- **HTTP Layer**: Handles HTTP requests/responses, validation, and JSON serialization
- **Domain Layer**: Business logic and domain models (User, Feed, FeedFollow)
- **Database Layer**: Generated by sqlc, handles database operations

### Authentication & Authorization

The application uses **API key-based authentication**:

1. User creates an account via `POST /v1/users` and receives an API key
2. The API key is automatically generated using SHA256 hash
3. Protected endpoints require `Authorization: ApiKey <key>` header
4. Middleware validates the API key and injects user context into the request
5. Handlers retrieve the authenticated user from the request context

**Key Design Decision**: User authentication state is stored in **request context** (not shared state) to prevent race conditions in concurrent requests.

### Domain Models vs Database Models

The codebase maintains separation between database models and domain models:

- **Database Models** (`internal/database/models.go`): Generated by sqlc, uses `sql.NullTime`, `uuid.NullUUID`
- **Domain Models** (`models.go`): Clean business objects with `time.Time` and `uuid.UUID`
- **Conversion Functions**: `databaseToUser()`, `databaseToFeed()`, `databaseToFeedFollows()` handle the mapping

This separation provides:

- Clean API responses without SQL null types
- Flexibility to change database structure without affecting API contracts
- Better testability of business logic

### Database Schema

The application has three main tables:

#### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### Feeds Table

```sql
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, url)
);
```

#### Feed Follows Table

```sql
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, feed_id)
);
```

### Database Timestamps

All tables use `TIMESTAMP WITH TIME ZONE` for timestamps:

- Stores timestamps in UTC
- Automatically handles timezone conversions
- `DEFAULT CURRENT_TIMESTAMP` auto-populates on insert
- More portable than `TIMESTAMP WITHOUT TIME ZONE`

## API Endpoints

### Public Endpoints

| Method | Endpoint     | Description            | Request Body         | Response              |
| ------ | ------------ | ---------------------- | -------------------- | --------------------- |
| GET    | `/v1/`       | Welcome message        | -                    | Plain text            |
| GET    | `/v1/health` | Health check           | -                    | `{"status": "ok"}`    |
| GET    | `/v1/error`  | Error testing endpoint | -                    | Error response        |
| POST   | `/v1/users`  | Create a new user      | `{"name": "string"}` | User object with UUID |

### Protected Endpoints (Require Authentication)

All protected endpoints require the `Authorization` header with an API key:

```
Authorization: ApiKey <your-api-key>
```

| Method | Endpoint           | Description                    | Request Body                           | Response                    |
| ------ | ------------------ | ------------------------------ | -------------------------------------- | --------------------------- |
| GET    | `/v1/users`        | Get current authenticated user | -                                      | User object                 |
| POST   | `/v1/feeds`        | Create a new RSS feed          | `{"title": "string", "url": "string"}` | Feed object                 |
| GET    | `/v1/feeds`        | Get all feeds                  | -                                      | Array of feed objects       |
| POST   | `/v1/feed_follows` | Follow an RSS feed             | `{"feed_id": "uuid"}`                  | FeedFollow object           |
| GET    | `/v1/feed_follows` | Get user's feed follows        | -                                      | Array of FeedFollow objects |

### Request/Response Examples

#### Create User

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "api_key": "a1b2c3d4e5f6...",
  "created_at": "2025-10-14T10:00:00Z",
  "updated_at": "2025-10-14T10:00:00Z"
}
```

#### Get Current User

```bash
curl -X GET http://localhost:8080/v1/users \
  -H "Authorization: ApiKey a1b2c3d4e5f6..."
```

#### Create Feed

```bash
curl -X POST http://localhost:8080/v1/feeds \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey a1b2c3d4e5f6..." \
  -d '{"title": "Tech Blog", "url": "https://example.com/feed.xml"}'
```

#### Follow Feed

```bash
curl -X POST http://localhost:8080/v1/feed_follows \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey a1b2c3d4e5f6..." \
  -d '{"feed_id": "550e8400-e29b-41d4-a716-446655440000"}'
```

## Development

### Database Management

```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d rss_aggregator

# Stop database
docker-compose down

# Stop and remove volumes (careful: this deletes data)
docker-compose down -v
```

## Database Migrations with Goose

Goose is used to manage database schema changes with version control.

### Migration File Structure

```
sql/
└── schema/
    ├── 001_users.sql
    ├── 002_feeds.sql
    └── 003_posts.sql
```

### Creating a Migration

```bash
# Using Makefile (easier)
make migrate-create NAME=add_users_table

# Or using goose directly
goose -dir sql/schema create add_users_table sql
```

### Migration File Format

```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
```

### Running Migrations

#### Using Makefile (Recommended)

The project includes a Makefile for convenient migration commands:

```bash
# Show available commands
make help

# Apply all pending migrations
make migrate-up

# Check migration status
make migrate-status

# Rollback last migration
make migrate-down

# Reset all migrations (rollback everything)
make migrate-reset

# Create a new migration
make migrate-create NAME=add_feeds_table
```

#### Using Goose Directly

If you prefer to use goose commands directly:

```bash
# Set database connection string (already in .env as DB_URL)
export DB_URL="postgresql://postgres:password@localhost:5432/rss_aggregator?sslmode=disable"

# Apply all pending migrations
goose -dir sql/schema postgres "$DB_URL" up

# Check migration status
goose -dir sql/schema postgres "$DB_URL" status

# Rollback last migration
goose -dir sql/schema postgres "$DB_URL" down

# Reset all migrations
goose -dir sql/schema postgres "$DB_URL" reset

# Apply specific version
goose -dir sql/schema postgres "$DB_URL" up-to 1

# Redo last migration (down + up)
goose -dir sql/schema postgres "$DB_URL" redo
```

### Common Goose Commands

| Command             | Description                          |
| ------------------- | ------------------------------------ |
| `up`                | Apply all pending migrations         |
| `down`              | Rollback the last migration          |
| `status`            | Show migration status                |
| `create <name> sql` | Create a new migration file          |
| `reset`             | Rollback all migrations              |
| `redo`              | Rollback and re-apply last migration |
| `version`           | Show current migration version       |

## SQL Code Generation with sqlc

sqlc generates type-safe Go code from SQL queries, eliminating the need for ORMs.

### sqlc Configuration

Create `sqlc.yaml` in the project root:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/queries/"
    schema: "sql/schema/"
    gen:
      go:
        package: "database"
        out: "internal/database"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: false
        emit_exact_table_names: false
```

### Writing SQL Queries

Create query files in `sql/queries/`:

```sql
-- name: CreateUser :one
INSERT INTO users (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

### Generating Go Code

```bash
# Generate type-safe Go code from SQL
sqlc generate

# This creates Go files in internal/database/ with:
# - Type-safe functions for each query
# - Struct definitions matching your tables
# - No runtime reflection or string building
```

### Using Generated Code

```go
import "github.com/yourusername/rss-aggregator/internal/database"

// Create a new user
user, err := db.CreateUser(ctx, database.CreateUserParams{
    ID:        uuid.New(),
    Name:      "John Doe",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
})

// Get a user
user, err := db.GetUser(ctx, userID)

// List all users
users, err := db.ListUsers(ctx)
```

### sqlc Benefits

- ✅ **Type safety**: Compile-time errors for SQL mistakes
- ✅ **No reflection**: Fast runtime performance
- ✅ **No ORM complexity**: Write plain SQL
- ✅ **Auto-completion**: IDE support for generated code
- ✅ **Easy testing**: Generated code is easy to mock

## Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o rss-aggregator

# Run production binary
./rss-aggregator
```

### Production Environment Variables

Ensure these are set in production:

```env
APP_PORT=8080
DB_URL=postgresql://user:password@host:5432/dbname?sslmode=require
```

**Important**:

- Use `sslmode=require` in production
- Store `.env` securely and never commit it to version control
- Consider using environment-specific secret management (AWS Secrets Manager, HashiCorp Vault, etc.)

## Project Dependencies

### Go Packages

```bash
# Install project dependencies
go mod download
```

| Package                    | Version | Purpose                      |
| -------------------------- | ------- | ---------------------------- |
| `github.com/go-chi/chi/v5` | v5.x    | HTTP router and middleware   |
| `github.com/go-chi/cors`   | latest  | CORS support                 |
| `github.com/joho/godotenv` | latest  | Environment variable loading |
| `github.com/lib/pq`        | latest  | PostgreSQL driver            |
| `github.com/google/uuid`   | latest  | UUID generation              |

### Development Tools

| Tool  | Installation                                              | Purpose             |
| ----- | --------------------------------------------------------- | ------------------- |
| goose | `go install github.com/pressly/goose/v3/cmd/goose@latest` | Database migrations |
| sqlc  | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`     | SQL code generation |

### Installing All Dependencies

```bash
# Install Go packages
go mod tidy

# Install development tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installations
goose --version
sqlc version
```

## Docker Support

The project includes Docker Compose for the PostgreSQL database. To add full Docker support for the application:

```dockerfile
# Dockerfile example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY src/go.* ./
RUN go mod download
COPY src/ ./
RUN go build -o rss-aggregator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rss-aggregator .
EXPOSE 8080
CMD ["./rss-aggregator"]
```

## Quick Reference

### Daily Development Workflow

```bash
# 1. Start database (if not already running)
docker-compose up -d

# 2. Run migrations (if you created new ones)
make migrate-up

# 3. Generate code from SQL (after modifying queries)
sqlc generate

# 4. Run application
go build && ./rss-aggregator
```

### Adding New Features

When adding a new feature (e.g., posts, comments):

1. **Create Migration**: `make migrate-create NAME=add_posts_table`
2. **Write SQL Schema**: Edit the new migration file in `sql/schema/`
3. **Apply Migration**: `make migrate-up`
4. **Write SQL Queries**: Create query file in `sql/queries/` (e.g., `posts.sql`)
5. **Generate Code**: `sqlc generate`
6. **Create Domain Model**: Add struct to `models.go`
7. **Create Conversion Function**: Add `databaseToPost()` function
8. **Create Handler**: Create `handler_post.go` with your endpoints
9. **Register Routes**: Add routes in `main.go`

### Common Commands Cheat Sheet

```bash
# Database
docker-compose up -d              # Start database
docker-compose down               # Stop database
docker-compose logs postgres      # View logs

# Migrations (using Makefile)
make migrate-up                              # Apply migrations
make migrate-status                          # Check status
make migrate-down                            # Rollback one
make migrate-reset                           # Reset all
make migrate-create NAME=add_users_table     # Create migration

# Code Generation
sqlc generate                     # Generate Go code from SQL

# Go Version Management
goenv versions                    # List installed versions
goenv install 1.23.2             # Install specific version
goenv global 1.23.2              # Set global version
goenv local 1.22.5               # Set project-specific version

# Development
air                              # Run with hot reload
go run .                         # Run application
go build -o rss-aggregator       # Build binary
go test ./...                    # Run tests
go mod tidy                      # Clean up dependencies
```

## Code Quality & Best Practices

### Project Conventions

1. **Error Handling**: Always check and handle errors explicitly
2. **Context Usage**: Pass context through the request lifecycle for cancellation and timeouts
3. **Domain Models**: Keep domain models separate from database models
4. **API Keys**: Never log or expose API keys in responses
5. **Timestamps**: Always use `TIMESTAMP WITH TIME ZONE` for consistency
6. **UUIDs**: Use UUIDs for primary keys to avoid sequential ID enumeration
7. **Unique Constraints**: Prevent duplicate data at the database level

### Security Best Practices

- ✅ API keys are hashed with SHA256
- ✅ SQL injection prevention via sqlc parameterized queries
- ✅ CORS configured for cross-origin requests
- ✅ Foreign key constraints with CASCADE deletes
- ✅ Unique constraints on user-resource relationships
- ⚠️ **TODO**: Add rate limiting
- ⚠️ **TODO**: Add request validation middleware
- ⚠️ **TODO**: Add HTTPS support
- ⚠️ **TODO**: Add unit and integration tests

### Code Generation

This project uses code generation for:

- **sqlc**: Generates type-safe Go code from SQL queries
- **goose**: Manages database schema versions

After modifying SQL files, always run:

```bash
sqlc generate
```

## Troubleshooting

### Port 5432 Already in Use

If you get "port already in use" error:

```bash
# Check what's using port 5432
lsof -i :5432

# Stop local PostgreSQL
brew services stop postgresql@14

# Or use different port in docker-compose.yml
ports:
  - "5433:5432"
```

### Goose Connection Issues

Make sure your `DB_URL` environment variable is set:

```bash
export DB_URL="postgresql://postgres:password@localhost:5432/rss_aggregator?sslmode=disable"
```

Or source it from `.env`:

```bash
source .env
```

### Go Version Issues

If `go version` shows wrong version:

```bash
# Ensure goenv shims are first in PATH
export PATH="$HOME/.goenv/shims:$PATH"

# Reload shell configuration
source ~/.bash_profile

# Verify
which go  # Should show: /Users/yourusername/.goenv/shims/go
```

## License

This project is licensed under the MIT License.

---

**Built with ❤️ using Go and PostgreSQL**
