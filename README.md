<div align="center">

<p style="font-size: 40px; font-weight: 700;">Tickr</p>

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
</div>


**Tickr** is a multi-tenant movie ticketing platform built for concurrency, scalability, and strict role-based access control.

## Features

- **Multi-tenant roles:** Separate, independent experiences for users, venues, and superadmins.
- **Async job pipeline:** Heavy tasks are offloaded to background workers to keep HTTP responses fast.
- **Secure by default:**
  - JWT authentication with access and refresh tokens
  - Bcrypt password hashing
  - Rate-limiting middleware against brute-force attempts and API abuse
- **Versioned migrations:** Database schema changes are managed with Goose.

## Documentation

Full server documentation is available in [`server/docs`](server/docs).

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Goose

### Clone the repository

```bash
git clone https://github.com/aprimr/tickr.git
cd tickr/server
```

### Configure environment variables

Create a `.env` file in the `server/` directory:

```env
# Server
PORT=8000
ALLOWED_ORIGINS=localhost:3000,https://www.example.com #`*` to allow all clients

# Auth
JWT_ACCESS_SECRET= #secret_access_secret
JWT_REFRESH_SECRET= #secret_refresh_secret

# Email (Brevo)
BREVO_API_KEY= #brevo_api_key
BREVO_FROM_EMAIL= #brevo_sender_email
BREVO_FROM_NAME=Tickr

# Database
DATABASE_URL_POOLED= #pooled_connection_string
GOOSE_DRIVER=postgres
GOOSE_DBSTRING= #postgres_connection_string
GOOSE_MIGRATION_DIR=db/migrations
```

### Run migrations

```bash
goose up
```

### Start the server

```bash
# Run directly
go run cmd/api/main.go

# Or with live reload (requires Air)
air
```