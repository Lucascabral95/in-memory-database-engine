# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go REST API for an e-commerce inventory management system. Despite the repository name "in-memory-database-engine", the application uses PostgreSQL with GORM as the ORM. The project follows a clean architecture pattern with clear separation of concerns.

## Build and Run Commands

```bash
# Run the application (requires .env with DATABASE_URL)
go run cmd/api/main.go

# Build the binary
go build -o bin/api cmd/api/main.go

# Download dependencies
go mod download

# Run tests (when implemented)
go test ./...
```

## Project Structure

```
cmd/api/          # Application entry point (main.go)
internal/
  ├── config/     # Configuration loading from environment
  ├── database/   # Database initialization and migrations
  ├── handler/    # HTTP handlers (Gin framework)
  ├── model/      # Domain models and DTOs
  ├── repository/ # Data access layer (GORM)
  └── service/    # Business logic layer
pkg/
  ├── response/   # HTTP response utilities
  └── utils/      # Shared utilities (password hashing)
```

## Architecture Patterns

### Dependency Injection Flow (in main.go)
Each domain follows this constructor chain:
1. `repository.NewXRepository(db)` → returns repository interface
2. `service.NewXService(repository)` → returns service struct
3. `handler.NewXHandler(service)` → returns handler struct

### Repository Pattern
Each repository is an interface (`CategoryRepository`, `UserRepository`, etc.) implemented by a private struct (`categoryRepository`, `userRepository`). This enables easy mocking for tests.

### Base Model
All models embed `BaseModel` which provides:
- UUID primary key (auto-generated via `BeforeCreate` hook)
- `CreatedAt`, `UpdatedAt` timestamps
- Soft delete support via `DeletedAt` (excluded from JSON)

### Domain Entities
Currently implemented:
- **User**: Authentication (register/login), password hashing via bcrypt
- **Category**: CRUD operations for product categories
- **Product**: (partial implementation)
- **Order**: (partial implementation)
- **OrderItem**: (partial implementation)
- **StockMovement**: (partial implementation)

## Configuration

Configuration is loaded via `internal/config` using environment variables:
- `DATABASE_URL`: PostgreSQL connection string (required)
- `PORT`: Server port (default: 8080)
- `ENV`: Environment - determines GORM log level (default: development)

The `.env` file is loaded via `godotenv` for local development.

## Key Conventions

- **Language**: Error messages and comments are in Spanish (e.g., "JSON invalido", "Usuario creado exitosamente")
- **UUID**: All entities use UUID v4 as primary keys
- **JSON Tags**: Fields use snake_case in JSON
- **Soft Delete**: GORM soft deletes are enabled; `DeletedAt` is omitted from JSON responses
- **Validation**: UUIDs are parsed in handlers before passing to service layer
