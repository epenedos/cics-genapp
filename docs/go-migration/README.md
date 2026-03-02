# GenApp Go Migration

This directory contains the Go + OpenTUI migration of the original CICS/COBOL General Insurance Policy Management System (GenApp).

## Overview

The migration transforms the 30-program COBOL/CICS application into a modern Go application with:

- **OpenTUI terminal interface** (tview/tcell) replacing BMS maps
- **PostgreSQL database** replacing DB2 + VSAM
- **Clean architecture** with repository/service/UI layers
- **REST API** for potential future web/mobile clients

## Project Structure

```
.
├── cmd/
│   └── genapp/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/                  # Domain models
│   ├── repository/              # Data access layer
│   ├── service/                 # Business logic layer
│   ├── handler/                 # HTTP handlers (optional API)
│   └── ui/
│       ├── views/               # Screen implementations
│       └── components/          # Reusable UI components
├── migrations/                  # Database migrations
├── scripts/                     # Utility scripts
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL 15 or later (or SQLite for development)

## Quick Start

### 1. Clone and navigate to the repository

```bash
cd /path/to/genapp
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Set up the database

Create a PostgreSQL database:

```bash
createdb genapp
```

Run migrations (once the migrations are in place):

```bash
psql -d genapp -f migrations/001_initial_schema.sql
```

### 4. Configure the application

Create a `config.yaml` file or use environment variables:

```yaml
database:
  host: localhost
  port: 5432
  user: genapp
  password: genapp
  dbname: genapp
  sslmode: disable

server:
  host: 0.0.0.0
  port: 8080

log_level: info
```

Alternatively, use environment variables with the `GENAPP_` prefix:

```bash
export GENAPP_DATABASE_HOST=localhost
export GENAPP_DATABASE_PORT=5432
export GENAPP_DATABASE_USER=genapp
export GENAPP_DATABASE_PASSWORD=genapp
export GENAPP_DATABASE_DBNAME=genapp
export GENAPP_LOG_LEVEL=debug
```

### 5. Build and run

```bash
go build -o bin/genapp ./cmd/genapp
./bin/genapp
```

Or run directly:

```bash
go run ./cmd/genapp
```

## Development

### Running tests

```bash
go test ./...
```

### Running tests with coverage

```bash
go test -cover ./...
```

### Building for production

```bash
go build -ldflags="-s -w" -o bin/genapp ./cmd/genapp
```

## Original COBOL to Go Mapping

| COBOL Layer | Go Layer | Description |
|-------------|----------|-------------|
| lgtestc1.cbl, lgtestp*.cbl | internal/ui/views/ | Terminal UI screens |
| lgacus01.cbl, lgicus01.cbl, etc. | internal/service/ | Business logic |
| lgacdb01.cbl, lgicdb01.cbl, etc. | internal/repository/ | Database access |
| lgcmarea.cpy, lgpolicy.cpy | internal/models/ | Data structures |

## Key Bindings

| Key | Action |
|-----|--------|
| Tab | Move to next field |
| Shift+Tab | Move to previous field |
| Enter | Submit form / Select option |
| Escape | Go back / Cancel |
| Ctrl+C | Exit application |

## License

This project is part of the GenApp sample and is licensed under the Eclipse Public License 2.0.
