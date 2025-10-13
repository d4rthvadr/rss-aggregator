# RSS Aggregator

A Go-based RSS feed aggregator built with Chi router and PostgreSQL.

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
├── docker-compose.yml      # PostgreSQL database setup
├── src/                    # Go source code
│   ├── main.go            # Main application entry point
│   ├── json.go            # JSON response utilities
│   ├── handler_readiness.go
│   ├── handler_error.go
│   ├── go.mod             # Go module dependencies
│   └── go.sum
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

### 3. Install Go Dependencies

```bash
cd src
go mod tidy
```

### 4. Run the Application

```bash
# From the src directory
go run .

# Or build and run
go build -o rss-aggregator
./rss-aggregator
```

## Database Connection

The application connects to PostgreSQL with these default settings:

- **Host**: localhost
- **Port**: 5432
- **Database**: rss_aggregator
- **User**: postgres
- **Password**: password

### Environment Variables

Create a `.env` file in the `src` directory for database configuration:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rss_aggregator
DB_SSLMODE=disable
```

## API Endpoints

The application provides the following endpoints:

- `GET /v1/healthz` - Health check endpoint
- `GET /v1/err` - Error testing endpoint

## Development

### Hot Reload with Air

For development with hot reload:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Initialize Air config
air init

# Run with hot reload
air
```

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
# Build binary
cd src
go build -o ../rss-aggregator

# Run production binary
cd ..
./rss-aggregator
```

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
| air   | `go install github.com/cosmtrek/air@latest`               | Hot reload          |

### Installing All Dependencies

```bash
# Install Go packages
go mod tidy

# Install development tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/cosmtrek/air@latest

# Verify installations
goose --version
sqlc version
air -v
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
# 1. Start database
docker-compose up -d

# 2. Run migrations
make migrate-up

# 3. Generate code from SQL
sqlc generate

# 4. Run application with hot reload
air

# Or run without hot reload
go run .
```

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
