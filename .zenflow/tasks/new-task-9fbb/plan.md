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

### [x] Step: Database Schema and Migrations
<!-- chat-id: 83a66761-2b7c-4058-87f5-963c610c4646 -->

Create PostgreSQL schema matching the COBOL data structures.

- [x] Create `migrations/001_initial_schema.sql` with all tables:
  - customers (from DB2-CUSTOMER)
  - policies (from DB2-POLICY)
  - motor_policies (from DB2-MOTOR)
  - endowment_policies (from DB2-ENDOWMENT)
  - house_policies (from DB2-HOUSE)
  - commercial_policies (from DB2-COMMERCIAL)
  - claims (from DB2-CLAIM)
  - counters (replacing Named Counters)
- [x] Create sequences for customer_num, policy_num, claim_num
- [x] Add `scripts/seed.sql` with test data
- [x] Test migration on local PostgreSQL

**Verification**: Run migrations and verify tables exist - PASSED

**Files created:**
- `migrations/001_initial_schema.sql` - Full PostgreSQL schema with:
  - 8 tables: customers, policies, motor_policies, endowment_policies, house_policies, commercial_policies, claims, counters
  - 3 sequences: customer_num_seq, policy_num_seq, claim_num_seq (starting at 1000000001)
  - Helper functions: next_customer_num(), next_policy_num(), next_claim_num(), increment_counter()
  - Auto-update triggers for updated_at timestamps
  - Foreign key constraints and indexes
  - Field sizes matching COBOL PIC definitions from lgcmarea.cpy and lgpolicy.cpy
- `scripts/seed.sql` - Test data with:
  - 10 sample customers
  - 10 policies (3 motor, 2 endowment, 3 house, 2 commercial)
  - 3 claims
  - Initial counter values
- `scripts/test_migration.sh` - Test script that runs PostgreSQL in Docker and validates migrations

---

### [x] Step: Domain Models and Repository Layer
<!-- chat-id: 976c2a8c-8e97-4c9d-9323-3bd0eff13924 -->

Implement Go domain models and data access layer.

- [x] Create `internal/models/customer.go` (Customer struct)
- [x] Create `internal/models/policy.go` (Policy, MotorPolicy, EndowmentPolicy, HousePolicy, CommercialPolicy)
- [x] Create `internal/models/claim.go` (Claim struct)
- [x] Create `internal/repository/db.go` (database connection pool)
- [x] Create `internal/repository/customer_repo.go` with CRUD operations
- [x] Create `internal/repository/policy_repo.go` with type-specific handling
- [x] Create `internal/repository/claim_repo.go`
- [x] Write unit tests for each repository

**Verification**: `go test ./internal/repository/...` - PASSED (unit tests pass, integration tests skip when DB unavailable)

**Files created:**
- `internal/models/customer.go` - Customer domain model with null-safe getters
- `internal/models/policy.go` - Policy, MotorPolicy, EndowmentPolicy, HousePolicy, CommercialPolicy models with PolicyType enum
- `internal/models/claim.go` - Claim domain model
- `internal/models/counter.go` - Counter model with predefined counter names
- `internal/repository/db.go` - Database connection pool with transaction support
- `internal/repository/helpers.go` - Null type conversion helpers
- `internal/repository/customer_repo.go` - Full CRUD for customers (equivalent to LGACDB01, LGICDB01, LGUCDB01)
- `internal/repository/policy_repo.go` - Full CRUD for all policy types with type-specific handling (equivalent to LGAPDB01, LGIPDB01, LGUPDB01, LGDPDB01)
- `internal/repository/claim_repo.go` - Full CRUD for claims
- `internal/repository/counter_repo.go` - Counter operations (atomic increment, get, set)
- `internal/repository/*_test.go` - Comprehensive test suites for all repositories

---

### [x] Step: Service Layer Implementation
<!-- chat-id: 129dc343-bb9c-4a20-a352-05579bc675a4 -->

Implement business logic services (replacing COBOL US programs).

- [x] Create `internal/service/customer_svc.go`:
  - Add customer (LGACUS01 equivalent)
  - Get customer (LGICUS01 equivalent)
  - Update customer (LGUCUS01 equivalent)
- [x] Create `internal/service/policy_svc.go`:
  - Add policy with type-specific details (LGAPOL01 equivalent)
  - Get policy (LGIPOL01 equivalent)
  - Update policy (LGUPOL01 equivalent)
  - Delete policy (LGDPOL01 equivalent)
- [x] Create `internal/service/counter_svc.go` (LGSETUP equivalent):
  - Generate customer numbers
  - Generate policy numbers
  - Atomic counter increments
- [x] Write unit tests for all services

**Verification**: `go test ./internal/service/...` - PASSED (unit tests pass, integration tests skip when DB unavailable)

**Files created:**
- `internal/service/customer_svc.go` - Customer service with:
  - Add (LGACUS01 equivalent) with validation and counter tracking
  - Get (LGICUS01 equivalent) with counter tracking
  - Update (LGUCUS01 equivalent) with validation and counter tracking
  - Delete, List, Search, Count operations
  - Input validation matching COBOL PIC definitions
  - Email format validation
- `internal/service/policy_svc.go` - Policy service with:
  - Add (LGAPOL01 equivalent) with type-specific validation
  - Get (LGIPOL01 equivalent) with type-specific details
  - Update (LGUPOL01 equivalent) with type-specific handling
  - Delete (LGDPOL01 equivalent) with counter tracking
  - GetByCustomer, List, ListByType operations
  - Support for all 4 policy types (Motor, Endowment, House, Commercial)
- `internal/service/counter_svc.go` - Counter service (LGSETUP equivalent) with:
  - NextCustomerNumber, NextPolicyNumber, NextClaimNumber generation
  - Counter get/set/increment/reset operations
  - GetStatistics for application metrics
  - InitializeCounters for startup
- `internal/service/customer_svc_test.go` - Customer service tests
- `internal/service/policy_svc_test.go` - Policy service tests
- `internal/service/counter_svc_test.go` - Counter service tests

---

### [x] Step: OpenTUI Application Framework
<!-- chat-id: 261496ed-e528-4539-af19-66e2a9fe2c3a -->

Set up the terminal UI application structure.

- [x] Create `internal/ui/app.go` (main TUI application)
- [x] Create `internal/ui/components/form.go` (reusable form component)
- [x] Create `internal/ui/components/menu.go` (menu component with option selection)
- [x] Implement key bindings (Enter, Escape/PF3, Tab navigation)
- [x] Set up 24x80 fixed terminal dimensions
- [x] Create basic navigation between screens

**Verification**: `go mod tidy && go build ./...` - PASSED, `go test ./...` - PASSED

**Files created:**
- `internal/ui/app.go` - Main TUI application with:
  - Screen type enumeration (Customer, Motor, Endowment, House, Commercial, Claim)
  - View interface for all screens
  - Services container for dependency injection
  - Page-based navigation with SwitchTo()
  - Global key bindings (Escape/F3 for back, Ctrl+C/F12 for exit)
  - Tab navigation between fields
- `internal/ui/components/form.go` - Reusable form component with:
  - Field types: Text, Numeric, Date, YesNo, Decimal
  - Input validation and acceptance functions
  - Tab/Backtab navigation between fields
  - Required field validation
  - Right-justify and zero-pad formatting
  - Helper functions: FormatCustomerNum(), FormatPolicyNum()
- `internal/ui/components/menu.go` - Menu component with:
  - Numbered option selection (1-9)
  - Enable/disable options
  - Pre-built menus: CustomerMenu(), PolicyMenu(), CommercialPolicyMenu(), ClaimMenu()
  - OperationType enum (Inquiry, Add, Delete, Update)
- `internal/ui/components/screen.go` - Base screen layout matching BMS 24x80 format:
  - Row 1: Screen ID + Title
  - Rows 4-7: Menu area
  - Rows 4-18: Form area
  - Row 22: Option selection
  - Row 24: Error/status message
- `internal/ui/views/base.go` - Base view implementing common functionality
- `internal/ui/views/main_menu.go` - Main navigation menu
- `internal/ui/views/customer_placeholder.go` - Customer screen (SSMAPC1) with:
  - All 10 form fields matching BMS definition
  - Menu options: Inquiry, Add, Update
  - Navigation to policy screens via F-keys
- `internal/ui/views/policy_placeholders.go` - Policy and Claim screens:
  - MotorPolicyView (SSMAPP1) - 13 fields
  - EndowmentPolicyView (SSMAPP2) - 11 fields
  - HousePolicyView (SSMAPP3) - 10 fields
  - CommercialPolicyView (SSMAPP4) - 10 fields
  - ClaimView (SSMAPP5) - 8 fields
- Updated `cmd/genapp/main.go` with UI initialization and navigation wiring

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
