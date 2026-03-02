# GenApp Go Migration

This directory contains the Go + OpenTUI migration of the original CICS/COBOL General Insurance Policy Management System (GenApp).

## Overview

The migration transforms the 30-program COBOL/CICS application into a modern Go application with:

- **OpenTUI terminal interface** (tview/tcell) replacing BMS maps
- **PostgreSQL database** replacing DB2 + VSAM
- **Clean architecture** with repository/service/UI layers

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
│   │   ├── customer.go          # Customer model
│   │   ├── policy.go            # Policy models (all types)
│   │   ├── claim.go             # Claim model
│   │   └── counter.go           # Counter model
│   ├── repository/              # Data access layer
│   │   ├── db.go                # Database connection
│   │   ├── customer_repo.go     # Customer CRUD
│   │   ├── policy_repo.go       # Policy CRUD
│   │   ├── claim_repo.go        # Claim CRUD
│   │   └── counter_repo.go      # Counter operations
│   ├── service/                 # Business logic layer
│   │   ├── customer_svc.go      # Customer operations
│   │   ├── policy_svc.go        # Policy operations
│   │   └── counter_svc.go       # Counter/sequence service
│   └── ui/
│       ├── app.go               # OpenTUI application
│       ├── views/               # Screen implementations
│       │   ├── customer.go      # Customer screen (SSMAPC1)
│       │   ├── motor.go         # Motor policy screen (SSMAPP1)
│       │   ├── endowment.go     # Endowment policy screen (SSMAPP2)
│       │   └── house.go         # House policy screen (SSMAPP3)
│       └── components/          # Reusable UI components
│           ├── form.go          # Form component
│           ├── menu.go          # Menu component
│           └── screen.go        # Base screen layout
├── migrations/
│   └── 001_initial_schema.sql   # Database schema
├── scripts/
│   ├── seed.sql                 # Test data
│   └── test_migration.sh        # Migration test script
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL 15 or later

## Quick Start

### 1. Download dependencies

```bash
go mod download
```

### 2. Set up the database

Create a PostgreSQL database:

```bash
createdb genapp
```

Run migrations:

```bash
psql -d genapp -f migrations/001_initial_schema.sql
```

Optionally seed with test data:

```bash
psql -d genapp -f scripts/seed.sql
```

### 3. Configure the application

Create a `config.yaml` file or use environment variables:

```yaml
database:
  host: localhost
  port: 5432
  user: genapp
  password: genapp
  dbname: genapp
  sslmode: disable

log_level: info
```

Alternatively, use environment variables with the `GENAPP_` prefix:

```bash
export GENAPP_DATABASE_HOST=localhost
export GENAPP_DATABASE_PORT=5432
export GENAPP_DATABASE_USER=genapp
export GENAPP_DATABASE_PASSWORD=genapp
export GENAPP_DATABASE_DBNAME=genapp
export GENAPP_DATABASE_SSLMODE=disable
export GENAPP_LOG_LEVEL=info
```

### 4. Build and run

```bash
# Build
go build -o bin/genapp ./cmd/genapp

# Run with database
./bin/genapp

# Run in demo mode (no database required)
./bin/genapp --no-db
```

Or run directly:

```bash
go run ./cmd/genapp
```

## Command Line Options

| Flag | Description |
|------|-------------|
| `--version` | Show version information |
| `--no-db` | Run without database connection (demo mode) |

## Key Bindings

### Global Keys

| Key | Action |
|-----|--------|
| `Tab` | Move to next field |
| `Shift+Tab` | Move to previous field |
| `Enter` | Submit form / Execute selected option |
| `Escape` | Go back / Exit from main screen |
| `F3` | Go back (same as Escape) |
| `F12` | Exit application |
| `Ctrl+C` | Emergency exit |

### Screen-Specific Keys

| Key | Action |
|-----|--------|
| `F1` | Navigate to Motor Policy screen (from Customer) |
| `F2` | Navigate to Endowment Policy screen (from Customer) |
| `F4` | Navigate to House Policy screen (from Customer) |
| `F5` | Navigate to Commercial Policy screen (from Customer) |
| `F6` | Navigate back to Customer screen (from Policy screens) |

### Menu Options

On each screen, select an option number and press Enter:

| Option | Action |
|--------|--------|
| `1` | Inquiry - Retrieve record by number |
| `2` | Add - Create new record |
| `3` | Delete - Remove record (policy screens only) |
| `4` | Update - Modify existing record |

## Screens

### Customer Screen (SSMAPC1)

Fields:
- Customer Number (10 digits, auto-generated for Add)
- First Name (max 10 chars)
- Last Name (max 20 chars, required)
- Date of Birth (yyyy-mm-dd)
- House Name (max 20 chars)
- House Number (max 4 chars)
- Postcode (max 8 chars)
- Phone: Home (max 20 chars)
- Phone: Mobile (max 20 chars)
- Email Address (max 100 chars)

### Motor Policy Screen (SSMAPP1)

Fields:
- Policy Number (10 digits)
- Customer Number (10 digits, required)
- Issue Date / Expiry Date (yyyy-mm-dd)
- Car Make / Model (max 15 chars each)
- Car Value (decimal)
- Registration (max 7 chars)
- Colour (max 8 chars)
- CC (engine size, integer)
- Manufacture Date (yyyy-mm-dd)
- Accidents (integer)
- Premium (decimal)

### Endowment Policy Screen (SSMAPP2)

Fields:
- Policy Number (10 digits)
- Customer Number (10 digits, required)
- Issue Date / Expiry Date (yyyy-mm-dd)
- Fund Name (max 10 chars)
- Term (1-99 years)
- Sum Assured (decimal)
- Life Assured (max 31 chars)
- With Profits (Y/N)
- Equities (Y/N)
- Managed Funds (Y/N)

### House Policy Screen (SSMAPP3)

Fields:
- Policy Number (10 digits)
- Customer Number (10 digits, required)
- Issue Date / Expiry Date (yyyy-mm-dd)
- Property Type (max 15 chars)
- Bedrooms (integer)
- House Value (decimal)
- House Name (max 20 chars)
- House Number (max 4 chars)
- Postcode (max 8 chars)

## Development

### Running tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/repository/...
go test ./internal/service/...
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
| lgsetup.cbl | internal/service/counter_svc.go | Counter initialization |
| ssmap.bms | internal/ui/views/*.go | BMS screen definitions |

## Database Schema

The PostgreSQL schema includes:

- `customers` - Customer records
- `policies` - Base policy table with type discriminator
- `motor_policies` - Motor-specific policy details
- `endowment_policies` - Endowment-specific policy details
- `house_policies` - House-specific policy details
- `commercial_policies` - Commercial-specific policy details
- `claims` - Insurance claims
- `counters` - Application counters (replaces CICS Named Counters)

Sequences for auto-generated IDs:
- `customer_num_seq` - Starting at 1000000001
- `policy_num_seq` - Starting at 1000000001
- `claim_num_seq` - Starting at 1000000001

## License

This project is part of the GenApp sample and is licensed under the Eclipse Public License 2.0.
