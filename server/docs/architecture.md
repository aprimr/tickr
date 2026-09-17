# Tickr Architecture

This document specifies the architectural patterns and layout of the **Tickr** backend.

Tickr follows a clean feature-based architecture.

## Tools and Technologies
* **Core:** Go, Chi Router, PostgreSQL 
* **Tools:** Goose


## Folder Structure

Following the standard Go project layout, the `server/` directory is structured as follows:

```text
server/
├── cmd/
│   └── api/
│       └── main.go            # Application entrypoint
├── internal/
│   ├── auth/                  # Authentication & RBAC feature module
│   │   ├── handler.go         # HTTP handlers and route endpoints
│   │   ├── service.go         # Core business logic
│   │   ├── repository.go      # Database query execution layer
│   │   ├── models.go          # Auth domain models and structs
│   │   ├── errors.go          # Feature-specific error definitions
│   │   └── dto.go             # Request/Response data transfer objects
│   ├── database/              # PostgreSQL connection pool initialization
│   ├── middleware/            # Auth, RBAC, logging, and recovery middleware
│   ├── models/                # Global domain structs and enum mappings
├── db/
│   └── migrations/            # Versioned Goose SQL migration files
├── docs/                      # Technical specifications and database schema specs
├── .env  
├── .gitignore
├── go.mod
└── go.sum