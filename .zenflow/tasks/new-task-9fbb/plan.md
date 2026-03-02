# Spec and build

## Configuration
- **Artifacts Path**: {@artifacts_path} → `.zenflow/tasks/{task_id}`

---

## Agent Instructions

Ask the user questions when anything is unclear or needs their input. This includes:
- Ambiguous or incomplete requirements
- Technical decisions that affect architecture or user experience
- Trade-offs that require business context

Do not make assumptions on important decisions — get clarification first.

If you are blocked and need user clarification, mark the current step with `[!]` in plan.md before stopping.

---

## Workflow Steps

### [x] Step: Technical Specification
<!-- chat-id: 47dcfef4-15a8-44d3-9ce4-9ef09fc5609e -->

**Complexity: HARD** - Full-stack migration of 30 COBOL programs to Go + OpenTUI

Technical specification saved to `.zenflow/tasks/new-task-9fbb/spec.md` containing:
- Current CICS/COBOL architecture analysis (30 programs, 4 BMS screens, 13 copybooks)
- Target Go + OpenTUI architecture design
- PostgreSQL database schema (replacing DB2 + VSAM)
- Domain models and repository patterns
- OpenTUI screen layouts matching BMS maps
- REST API endpoint design
- Phased migration strategy

---

### [x] Step: Project Foundation Setup
<!-- chat-id: 17911c83-5cba-4fde-8fa0-b28fac6f3e02 -->

Set up the Go project structure and dependencies.

- [x] Initialize Go module (`go mod init`)
- [x] Create project directory structure (`cmd/`, `internal/`, `migrations/`)
- [x] Add dependencies: tview (OpenTUI), tcell, pq, sqlx, viper, zerolog
- [x] Create `go.mod` and `go.sum`
- [x] Set up basic configuration management (`internal/config/config.go`)
- [x] Create initial `README.md` with setup instructions

**Verification**: `go mod tidy && go build ./...` - PASSED

**Files created:**
- `go.mod` - Go module definition with all required dependencies
- `go.sum` - Dependency checksums
- `cmd/genapp/main.go` - Application entry point
- `internal/config/config.go` - Configuration management with Viper
- `docs/go-migration/README.md` - Setup and development instructions
- Directory structure: `internal/{models,repository,service,handler,ui/{views,components}}`, `migrations/`, `scripts/`
- Updated `.gitignore` with Go-specific patterns

---

### [ ] Step: Database Schema and Migrations

Create PostgreSQL schema matching the COBOL data structures.

- [ ] Create `migrations/001_initial_schema.sql` with all tables:
  - customers (from DB2-CUSTOMER)
  - policies (from DB2-POLICY)
  - motor_policies (from DB2-MOTOR)
  - endowment_policies (from DB2-ENDOWMENT)
  - house_policies (from DB2-HOUSE)
  - commercial_policies (from DB2-COMMERCIAL)
  - claims (from DB2-CLAIM)
  - counters (replacing Named Counters)
- [ ] Create sequences for customer_num, policy_num, claim_num
- [ ] Add `scripts/seed.sql` with test data
- [ ] Test migration on local PostgreSQL

**Verification**: Run migrations and verify tables exist

---

### [ ] Step: Domain Models and Repository Layer

Implement Go domain models and data access layer.

- [ ] Create `internal/models/customer.go` (Customer struct)
- [ ] Create `internal/models/policy.go` (Policy, MotorPolicy, EndowmentPolicy, HousePolicy, CommercialPolicy)
- [ ] Create `internal/models/claim.go` (Claim struct)
- [ ] Create `internal/repository/db.go` (database connection pool)
- [ ] Create `internal/repository/customer_repo.go` with CRUD operations
- [ ] Create `internal/repository/policy_repo.go` with type-specific handling
- [ ] Create `internal/repository/claim_repo.go`
- [ ] Write unit tests for each repository

**Verification**: `go test ./internal/repository/...`

---

### [ ] Step: Service Layer Implementation

Implement business logic services (replacing COBOL US programs).

- [ ] Create `internal/service/customer_svc.go`:
  - Add customer (LGACUS01 equivalent)
  - Get customer (LGICUS01 equivalent)
  - Update customer (LGUCUS01 equivalent)
- [ ] Create `internal/service/policy_svc.go`:
  - Add policy with type-specific details (LGAPOL01 equivalent)
  - Get policy (LGIPOL01 equivalent)
  - Update policy (LGUPOL01 equivalent)
  - Delete policy (LGDPOL01 equivalent)
- [ ] Create `internal/service/counter_svc.go` (LGSETUP equivalent):
  - Generate customer numbers
  - Generate policy numbers
  - Atomic counter increments
- [ ] Write unit tests for all services

**Verification**: `go test ./internal/service/...`

---

### [ ] Step: OpenTUI Application Framework

Set up the terminal UI application structure.

- [ ] Create `internal/ui/app.go` (main TUI application)
- [ ] Create `internal/ui/components/form.go` (reusable form component)
- [ ] Create `internal/ui/components/menu.go` (menu component with option selection)
- [ ] Implement key bindings (Enter, Escape/PF3, Tab navigation)
- [ ] Set up 24x80 fixed terminal dimensions
- [ ] Create basic navigation between screens

**Verification**: Build and run basic UI skeleton

---

### [ ] Step: Customer Screen Implementation

Implement the customer menu screen (SSMAPC1 equivalent).

- [ ] Create `internal/ui/views/customer.go`
- [ ] Layout: Title, menu options, form fields, error display
- [ ] Fields: Customer Number, First/Last Name, DOB, Address, Phone, Email
- [ ] Option dropdown: 1-Inquiry, 2-Add, 3-Delete, 4-Update
- [ ] Field validation (numeric for customer number, date format for DOB)
- [ ] Connect to customer service for operations
- [ ] Display error messages in error field

**Verification**: Manual testing of customer CRUD operations

---

### [ ] Step: Policy Screens Implementation

Implement the policy screens (SSMAPP1, SSMAPP2, SSMAPP3 equivalents).

- [ ] Create `internal/ui/views/motor.go` (Motor policy - SSMAPP1):
  - Policy/Customer numbers, dates, car details, premium, accidents
- [ ] Create `internal/ui/views/endowment.go` (Endowment policy - SSMAPP2):
  - Policy/Customer numbers, dates, fund details, Y/N checkboxes
- [ ] Create `internal/ui/views/house.go` (House policy - SSMAPP3):
  - Policy/Customer numbers, property details
- [ ] Connect each screen to policy service
- [ ] Add navigation between customer and policy screens

**Verification**: Manual testing of each policy type CRUD

---

### [ ] Step: Integration and Error Handling

Wire up all components and add comprehensive error handling.

- [ ] Create `cmd/genapp/main.go` (entry point)
- [ ] Initialize database connection
- [ ] Create service instances with repositories
- [ ] Pass services to UI views
- [ ] Add transaction handling for multi-table operations
- [ ] Implement comprehensive error display
- [ ] Add graceful shutdown handling

**Verification**: End-to-end test of complete workflow

---

### [ ] Step: Documentation and Final Testing

Complete the migration with documentation and testing.

- [ ] Update `README.md` with:
  - Setup instructions
  - Database configuration
  - Running the application
  - Key mappings
- [ ] Create comparison test: same operations on COBOL vs Go
- [ ] Performance testing
- [ ] Write `{@artifacts_path}/report.md` with:
  - What was implemented
  - How the solution was tested
  - Challenges encountered

**Verification**: `go test ./... && go build -o bin/genapp ./cmd/genapp`
