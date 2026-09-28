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
│       └── main.go                  # Application entrypoint
├── internal/
│   ├── auth/                         
│   │   ├── errors.go                # Feature-specific error definitions
│   │   ├── handler.go               # HTTP handlers and route endpoints
│   │   ├── models.go                # Auth domain models and structs
│   │   ├── repository.go            # Database query execution layer
│   │   ├── routes.go                # Feature-specific route declarations
│   │   ├── service.go               # Core business logic
│   │   └── utils.go                 # Feature-specific helper utilities
│   ├── db/                          # PostgreSQL connection pool initialization
│   ├── domain/                      # Global domain structs and enums
│   ├── email/                       # Email dispatchers and HTML templates
│   ├── middleware/                  # Middleware for Auth, RBAC, rate-limiting and so on
│   ├── utils/                       # Global utility modules
│   └── pkg/                         # Shared internal packages
├── 
├── db/                              
│   └── migrations/                  # SQL migration files
├── docs/                            # Server Architecture and API documentation
├── .air.toml                        
├── .env                             
├── .gitignore                       
├── go.mod                           
└── go.sum                           