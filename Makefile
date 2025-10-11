.PHONY: help migrate-up migrate-down migrate-status migrate-reset migrate-create

# Load environment variables from .env file
include .env
export

# Database connection string
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Goose configuration
MIGRATIONS_DIR := sql/schema
DB_DRIVER := postgres

# Default target - show help
help:
	@echo "Available commands:"
	@echo "  make migrate-up       - Apply all pending migrations"
	@echo "  make migrate-down     - Rollback the last migration"
	@echo "  make migrate-status   - Show migration status"
	@echo "  make migrate-reset    - Rollback all migrations"
	@echo "  make migrate-create NAME=<name> - Create a new migration file"
	@echo ""
	@echo "Example:"
	@echo "  make migrate-create NAME=add_users_table"

# Apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" up

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" down

# Show migration status
migrate-status:
	@echo "Migration status:"
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" status

# Rollback all migrations
migrate-reset:
	@echo "Resetting all migrations..."
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" reset

# Create a new migration
# Usage: make migrate-create NAME=add_users_table
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
