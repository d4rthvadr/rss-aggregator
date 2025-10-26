# RSS Aggregator

A Go-based RSS feed aggregator built with Chi router and PostgreSQL. This application allows users to manage RSS feeds, follow feeds, and aggregate content from multiple sources.

## Features

- ✅ User management with auto-generated API keys
- ✅ API key-based authentication
- ✅ Create and manage RSS feeds
- ✅ Follow/unfollow RSS feeds
- ✅ **RSS Feed Scraping**: Background worker that automatically fetches and parses RSS feeds
- ✅ **Post Storage**: Store individual RSS posts/articles from feeds
- ✅ **Concurrent Processing**: Multi-threaded feed scraping with configurable concurrency
- ✅ **Smart Feed Rotation**: Fetches feeds based on last update time for fair distribution
- ✅ **Post Retrieval**: Get posts for users based on their followed feeds
- ✅ Thread-safe request handling with context-based authentication
- ✅ Type-safe database queries with sqlc
- ✅ Database migrations with goose
- ✅ Clean architecture with domain/database model separation
- ✅ CORS support for web clients
- ✅ PostgreSQL with UUID primary keys and timezone-aware timestamps

## Installation

# Install Go 1.23.2 with Go Version Manager

goenv install 1.23.2
goenv global 1.23.2

````

### 2. Install Development Tools

```bash
# Install goose for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc for type-safe SQL code generation
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installations
goose --version
sqlc version
````

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

Create a `.env` file from the example:

```bash
cp .env.example .env
```

Update the values in `.env` if needed (default values should work for local development).

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

## API Endpoints

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
| GET    | `/v1/posts`        | Get posts from followed feeds  | `?limit=20` (optional)                 | Array of Post objects       |

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

#### Get Posts

```bash
curl -X GET http://localhost:8080/v1/posts?limit=10 \
  -H "Authorization: ApiKey a1b2c3d4e5f6..."
```

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

**Note**: For Makefile migration commands to work, ensure `DB_URL` is exported in your shell:

```bash
export DB_URL="postgresql://postgres:password@localhost:5432/rss_aggregator?sslmode=disable"
```

Or source your `.env` file:

```bash
set -a
source .env
set +a
```

## Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o rss-aggregator

# Run production binary
./rss-aggregator
```

### Development Tools

| Tool  | Installation                                              | Purpose             |
| ----- | --------------------------------------------------------- | ------------------- |
| goose | `go install github.com/pressly/goose/v3/cmd/goose@latest` | Database migrations |
| sqlc  | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`     | SQL code generation |

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


# Development
go run .                         # Run application
go build -o rss-aggregator       # Build binary
go mod tidy                      # Clean up dependencies

# Run application
go build && ./rss-aggregator
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
