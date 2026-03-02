# Technical Specification: CICS/COBOL to OpenTUI + Go Migration

## Task Complexity Assessment: **HARD**

This is a full-stack migration of an enterprise mainframe application (General Insurance Policy Management System - GENAPP) from CICS/COBOL to a modern architecture with OpenTUI for the terminal UI and Go for the backend. The complexity arises from:

- 30 COBOL programs with intricate business logic
- 4 BMS screen definitions with complex field mappings
- 13 copybooks defining data structures
- DB2 + VSAM data access patterns
- Inter-program communication via COMMAREA
- Named counters and transient storage queues

---

## 1. Current Application Overview

### 1.1 Technology Stack (Current)
- **Runtime**: IBM CICS Transaction Server
- **Language**: COBOL (Enterprise COBOL)
- **UI**: BMS (Basic Mapping Support) - 3270 terminal screens
- **Database**: IBM DB2 with embedded SQL
- **File System**: VSAM KSDS (Key Sequenced Data Sets)
- **IPC**: COMMAREA (32,500 bytes shared memory)
- **Counters**: CICS Named Counter Server

### 1.2 Application Functionality
The GENAPP system manages:
- **Customers**: Add, Inquiry, Update operations
- **Policies**: Motor, Endowment, House, Commercial types
- **Policy Lifecycle**: Add, Inquiry, Update, Delete
- **Claims**: Associated with commercial policies
- **Statistics**: Transaction monitoring and reporting

### 1.3 Source File Inventory

| Category | Count | Files |
|----------|-------|-------|
| COBOL Programs | 30 | lgacdb01.cbl, lgacdb02.cbl, lgacus01.cbl, lgacvs01.cbl, lgapdb01.cbl, lgapol01.cbl, lgapvs01.cbl, lgdpdb01.cbl, lgdpol01.cbl, lgdpvs01.cbl, lgicdb01.cbl, lgicus01.cbl, lgicvs01.cbl, lgipdb01.cbl, lgipol01.cbl, lgipvs01.cbl, lgstsq.cbl, lgsetup.cbl, lgtestc1.cbl, lgtestp1.cbl, lgtestp2.cbl, lgtestp3.cbl, lgtestp4.cbl, lgucdb01.cbl, lgucus01.cbl, lgucvs01.cbl, lgupdb01.cbl, lgupol01.cbl, lgupvs01.cbl, lgwebst5.cbl |
| BMS Maps | 1 | ssmap.bms (4 screen definitions) |
| Copybooks | 13 | lgcmarea.cpy, lgpolicy.cpy, pollook.cpy, polloo2.cpy, soaic01.cpy, soaipb1.cpy, soaipe1.cpy, soaiph1.cpy, soaipm1.cpy, soavcii.cpy, soavcio.cpy, soavpii.cpy, soavpio.cpy |

---

## 2. Target Architecture

### 2.1 Technology Stack (Target)
- **Frontend**: OpenTUI (Go-based terminal UI framework)
- **Backend**: Go 1.21+
- **Database**: PostgreSQL 15+ (or SQLite for development)
- **API**: REST with JSON payloads (optional gRPC for internal services)
- **Session**: In-memory or Redis for state management
- **Counters**: Atomic Go counters or Redis INCR

### 2.2 Architecture Pattern
```
┌─────────────────────────────────────────────────────────────┐
│                      OpenTUI Frontend                        │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │Customer │  │ Motor   │  │Endowment│  │ House/Commercial│ │
│  │ Screen  │  │ Screen  │  │ Screen  │  │     Screens     │ │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────────┬────────┘ │
└───────┼────────────┼────────────┼────────────────┼──────────┘
        │            │            │                │
        └────────────┴────────────┴────────────────┘
                            │
                     ┌──────┴──────┐
                     │   Go HTTP   │
                     │   Server    │
                     └──────┬──────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────┴───────┐  ┌───────┴───────┐  ┌────────┴────────┐
│   Customer    │  │    Policy     │  │   Statistics    │
│   Service     │  │   Service     │  │    Service      │
└───────┬───────┘  └───────┬───────┘  └────────┬────────┘
        │                  │                   │
        └──────────────────┼───────────────────┘
                           │
                    ┌──────┴──────┐
                    │  PostgreSQL │
                    │   Database  │
                    └─────────────┘
```

---

## 3. Data Model Migration

### 3.1 Database Schema (PostgreSQL)

```sql
-- Customer table (maps from DB2-CUSTOMER)
CREATE TABLE customers (
    id              SERIAL PRIMARY KEY,
    customer_num    VARCHAR(10) UNIQUE NOT NULL,
    first_name      VARCHAR(10),
    last_name       VARCHAR(20),
    date_of_birth   DATE,
    house_name      VARCHAR(20),
    house_number    VARCHAR(4),
    postcode        VARCHAR(8),
    phone_home      VARCHAR(20),
    phone_mobile    VARCHAR(20),
    email_address   VARCHAR(100),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Policy master table (maps from DB2-POLICY)
CREATE TABLE policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,
    customer_num    VARCHAR(10) REFERENCES customers(customer_num),
    policy_type     CHAR(1) NOT NULL CHECK (policy_type IN ('E','H','M','C')),
    issue_date      DATE,
    expiry_date     DATE,
    last_changed    TIMESTAMP,
    broker_id       VARCHAR(10),
    brokers_ref     VARCHAR(10),
    payment         DECIMAL(8,2),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Motor policy details (maps from DB2-MOTOR)
CREATE TABLE motor_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE REFERENCES policies(policy_num) ON DELETE CASCADE,
    make            VARCHAR(15),
    model           VARCHAR(15),
    value           DECIMAL(8,2),
    reg_number      VARCHAR(7),
    colour          VARCHAR(8),
    cc              INTEGER,
    manufactured    DATE,
    premium         DECIMAL(8,2),
    accidents       INTEGER DEFAULT 0
);

-- Endowment policy details (maps from DB2-ENDOWMENT)
CREATE TABLE endowment_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE REFERENCES policies(policy_num) ON DELETE CASCADE,
    with_profits    BOOLEAN DEFAULT FALSE,
    equities        BOOLEAN DEFAULT FALSE,
    managed_fund    BOOLEAN DEFAULT FALSE,
    fund_name       VARCHAR(10),
    term            INTEGER,
    sum_assured     DECIMAL(10,2),
    life_assured    VARCHAR(31)
);

-- House policy details (maps from DB2-HOUSE)
CREATE TABLE house_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE REFERENCES policies(policy_num) ON DELETE CASCADE,
    property_type   VARCHAR(15),
    bedrooms        INTEGER,
    value           DECIMAL(10,2),
    house_name      VARCHAR(20),
    house_number    VARCHAR(4),
    postcode        VARCHAR(8)
);

-- Commercial policy details (maps from DB2-COMMERCIAL)
CREATE TABLE commercial_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE REFERENCES policies(policy_num) ON DELETE CASCADE,
    address         TEXT,
    postcode        VARCHAR(8),
    latitude        VARCHAR(11),
    longitude       VARCHAR(11),
    customer        TEXT,
    property_type   TEXT,
    fire_peril      INTEGER,
    fire_premium    DECIMAL(10,2),
    crime_peril     INTEGER,
    crime_premium   DECIMAL(10,2),
    flood_peril     INTEGER,
    flood_premium   DECIMAL(10,2),
    weather_peril   INTEGER,
    weather_premium DECIMAL(10,2),
    status          INTEGER,
    reject_reason   TEXT
);

-- Claims (maps from DB2-CLAIM)
CREATE TABLE claims (
    id              SERIAL PRIMARY KEY,
    claim_num       VARCHAR(10) UNIQUE NOT NULL,
    policy_num      VARCHAR(10) REFERENCES policies(policy_num) ON DELETE CASCADE,
    claim_date      DATE,
    paid            DECIMAL(10,2),
    value           DECIMAL(10,2),
    cause           TEXT,
    observations    TEXT
);

-- Statistics counters (replaces Named Counters)
CREATE TABLE counters (
    name            VARCHAR(50) PRIMARY KEY,
    value           BIGINT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Application sequences
CREATE SEQUENCE customer_num_seq START 1000000001;
CREATE SEQUENCE policy_num_seq START 1000000001;
CREATE SEQUENCE claim_num_seq START 1000000001;
```

### 3.2 Go Domain Models

```go
// internal/models/customer.go
type Customer struct {
    ID           int64     `json:"id" db:"id"`
    CustomerNum  string    `json:"customer_num" db:"customer_num"`
    FirstName    string    `json:"first_name" db:"first_name"`
    LastName     string    `json:"last_name" db:"last_name"`
    DateOfBirth  time.Time `json:"date_of_birth" db:"date_of_birth"`
    HouseName    string    `json:"house_name" db:"house_name"`
    HouseNumber  string    `json:"house_number" db:"house_number"`
    Postcode     string    `json:"postcode" db:"postcode"`
    PhoneHome    string    `json:"phone_home" db:"phone_home"`
    PhoneMobile  string    `json:"phone_mobile" db:"phone_mobile"`
    EmailAddress string    `json:"email_address" db:"email_address"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// internal/models/policy.go
type PolicyType string
const (
    PolicyTypeEndowment  PolicyType = "E"
    PolicyTypeHouse      PolicyType = "H"
    PolicyTypeMotor      PolicyType = "M"
    PolicyTypeCommercial PolicyType = "C"
)

type Policy struct {
    ID          int64      `json:"id" db:"id"`
    PolicyNum   string     `json:"policy_num" db:"policy_num"`
    CustomerNum string     `json:"customer_num" db:"customer_num"`
    PolicyType  PolicyType `json:"policy_type" db:"policy_type"`
    IssueDate   time.Time  `json:"issue_date" db:"issue_date"`
    ExpiryDate  time.Time  `json:"expiry_date" db:"expiry_date"`
    LastChanged time.Time  `json:"last_changed" db:"last_changed"`
    BrokerID    string     `json:"broker_id" db:"broker_id"`
    BrokersRef  string     `json:"brokers_ref" db:"brokers_ref"`
    Payment     float64    `json:"payment" db:"payment"`

    // Embedded type-specific details (interface for polymorphism)
    Details     PolicyDetails `json:"details,omitempty"`
}

type PolicyDetails interface {
    Type() PolicyType
}

type MotorPolicy struct {
    PolicyNum    string    `json:"policy_num" db:"policy_num"`
    Make         string    `json:"make" db:"make"`
    Model        string    `json:"model" db:"model"`
    Value        float64   `json:"value" db:"value"`
    RegNumber    string    `json:"reg_number" db:"reg_number"`
    Colour       string    `json:"colour" db:"colour"`
    CC           int       `json:"cc" db:"cc"`
    Manufactured time.Time `json:"manufactured" db:"manufactured"`
    Premium      float64   `json:"premium" db:"premium"`
    Accidents    int       `json:"accidents" db:"accidents"`
}

type EndowmentPolicy struct {
    PolicyNum    string  `json:"policy_num" db:"policy_num"`
    WithProfits  bool    `json:"with_profits" db:"with_profits"`
    Equities     bool    `json:"equities" db:"equities"`
    ManagedFund  bool    `json:"managed_fund" db:"managed_fund"`
    FundName     string  `json:"fund_name" db:"fund_name"`
    Term         int     `json:"term" db:"term"`
    SumAssured   float64 `json:"sum_assured" db:"sum_assured"`
    LifeAssured  string  `json:"life_assured" db:"life_assured"`
}

type HousePolicy struct {
    PolicyNum    string  `json:"policy_num" db:"policy_num"`
    PropertyType string  `json:"property_type" db:"property_type"`
    Bedrooms     int     `json:"bedrooms" db:"bedrooms"`
    Value        float64 `json:"value" db:"value"`
    HouseName    string  `json:"house_name" db:"house_name"`
    HouseNumber  string  `json:"house_number" db:"house_number"`
    Postcode     string  `json:"postcode" db:"postcode"`
}

type CommercialPolicy struct {
    PolicyNum      string  `json:"policy_num" db:"policy_num"`
    Address        string  `json:"address" db:"address"`
    Postcode       string  `json:"postcode" db:"postcode"`
    Latitude       string  `json:"latitude" db:"latitude"`
    Longitude      string  `json:"longitude" db:"longitude"`
    Customer       string  `json:"customer" db:"customer"`
    PropertyType   string  `json:"property_type" db:"property_type"`
    FirePeril      int     `json:"fire_peril" db:"fire_peril"`
    FirePremium    float64 `json:"fire_premium" db:"fire_premium"`
    CrimePeril     int     `json:"crime_peril" db:"crime_peril"`
    CrimePremium   float64 `json:"crime_premium" db:"crime_premium"`
    FloodPeril     int     `json:"flood_peril" db:"flood_peril"`
    FloodPremium   float64 `json:"flood_premium" db:"flood_premium"`
    WeatherPeril   int     `json:"weather_peril" db:"weather_peril"`
    WeatherPremium float64 `json:"weather_premium" db:"weather_premium"`
    Status         int     `json:"status" db:"status"`
    RejectReason   string  `json:"reject_reason" db:"reject_reason"`
}

type Claim struct {
    ID           int64     `json:"id" db:"id"`
    ClaimNum     string    `json:"claim_num" db:"claim_num"`
    PolicyNum    string    `json:"policy_num" db:"policy_num"`
    ClaimDate    time.Time `json:"claim_date" db:"claim_date"`
    Paid         float64   `json:"paid" db:"paid"`
    Value        float64   `json:"value" db:"value"`
    Cause        string    `json:"cause" db:"cause"`
    Observations string    `json:"observations" db:"observations"`
}
```

---

## 4. Screen Migration (BMS to OpenTUI)

### 4.1 Screen Mapping

| BMS Map | Description | OpenTUI View |
|---------|-------------|--------------|
| SSMAPC1 | Customer Menu | `views/customer.go` |
| SSMAPP1 | Motor Policy | `views/motor_policy.go` |
| SSMAPP2 | Endowment Policy | `views/endowment_policy.go` |
| SSMAPP3 | House Policy | `views/house_policy.go` |

### 4.2 Field Translation Rules

| BMS Attribute | OpenTUI Equivalent |
|---------------|-------------------|
| UNPROT | Editable input field |
| PROT, ASKIP | Read-only/Label |
| BRT | Bold style |
| NUM | Numeric input validation |
| IC (Insert Cursor) | Initial focus |
| FSET (Field Set) | Default value |
| MUSTENTER | Required field validation |
| RIGHT, ZERO | Right-justify numeric with leading zeros |

### 4.3 Customer Screen Layout (Example)

```go
// views/customer.go
type CustomerView struct {
    app         *tview.Application
    form        *tview.Form
    errorText   *tview.TextView

    // Form fields
    customerNum *tview.InputField
    firstName   *tview.InputField
    lastName    *tview.InputField
    dob         *tview.InputField
    houseName   *tview.InputField
    houseNum    *tview.InputField
    postcode    *tview.InputField
    phoneHome   *tview.InputField
    phoneMobile *tview.InputField
    email       *tview.InputField
    option      *tview.DropDown
}

func (v *CustomerView) Layout() tview.Primitive {
    // 24x80 terminal layout matching SSMAPC1
    grid := tview.NewGrid().
        SetRows(1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, -1, 1, 1).
        SetColumns(30, 20, 30)

    // Row 1: Title
    grid.AddItem(tview.NewTextView().SetText("General Insurance Customer Menu"), 0, 0, 1, 3, 0, 0, false)

    // Menu options (rows 4-7)
    menu := tview.NewTextView().SetText("1. Cust Inquiry\n2. Cust Add\n3. [Reserved]\n4. Cust Update")
    grid.AddItem(menu, 2, 0, 4, 1, 0, 0, false)

    // Form fields
    v.form = tview.NewForm().
        AddInputField("Cust Number", "", 10, tview.InputFieldInteger, nil).
        AddInputField("First Name", "", 10, nil, nil).
        AddInputField("Last Name", "", 20, nil, nil).
        AddInputField("DOB (yyyy-mm-dd)", "", 10, nil, nil).
        AddInputField("House Name", "", 20, nil, nil).
        AddInputField("House Number", "", 4, nil, nil).
        AddInputField("Postcode", "", 8, nil, nil).
        AddInputField("Phone: Home", "", 20, nil, nil).
        AddInputField("Phone: Mobile", "", 20, nil, nil).
        AddInputField("Email Address", "", 27, nil, nil).
        AddDropDown("Option", []string{"1-Inquiry", "2-Add", "3-Delete", "4-Update"}, 0, nil)

    grid.AddItem(v.form, 2, 1, 10, 2, 0, 0, true)

    // Error field (row 24)
    v.errorText = tview.NewTextView().SetTextColor(tcell.ColorRed)
    grid.AddItem(v.errorText, 13, 0, 1, 3, 0, 0, false)

    return grid
}
```

---

## 5. Go Backend Structure

### 5.1 Project Layout

```
genapp/
├── cmd/
│   └── genapp/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   ├── customer.go          # Customer domain model
│   │   ├── policy.go            # Policy domain models
│   │   └── claim.go             # Claim domain model
│   ├── repository/
│   │   ├── customer_repo.go     # Customer data access
│   │   ├── policy_repo.go       # Policy data access
│   │   └── claim_repo.go        # Claim data access
│   ├── service/
│   │   ├── customer_svc.go      # Customer business logic
│   │   ├── policy_svc.go        # Policy business logic
│   │   └── counter_svc.go       # Counter/sequence service
│   ├── handler/
│   │   ├── customer_handler.go  # Customer HTTP handlers
│   │   └── policy_handler.go    # Policy HTTP handlers
│   └── ui/
│       ├── app.go               # OpenTUI application
│       ├── views/
│       │   ├── customer.go      # Customer screen
│       │   ├── motor.go         # Motor policy screen
│       │   ├── endowment.go     # Endowment policy screen
│       │   └── house.go         # House policy screen
│       └── components/
│           ├── form.go          # Reusable form component
│           └── menu.go          # Menu component
├── migrations/
│   └── 001_initial_schema.sql   # Database migrations
├── scripts/
│   └── seed.sql                 # Test data seed
├── go.mod
├── go.sum
└── README.md
```

### 5.2 Program Mapping (COBOL to Go)

| COBOL Program | Go Equivalent | Layer |
|--------------|---------------|-------|
| lgtestc1.cbl | ui/views/customer.go | UI |
| lgtestp1.cbl | ui/views/motor.go | UI |
| lgtestp2.cbl | ui/views/endowment.go | UI |
| lgtestp3.cbl | ui/views/house.go | UI |
| lgacus01.cbl | service/customer_svc.go (Add) | Service |
| lgacdb01.cbl | repository/customer_repo.go (Create) | Repository |
| lgacvs01.cbl | N/A (VSAM not needed) | N/A |
| lgicus01.cbl | service/customer_svc.go (Get) | Service |
| lgicdb01.cbl | repository/customer_repo.go (FindByNum) | Repository |
| lgucus01.cbl | service/customer_svc.go (Update) | Service |
| lgucdb01.cbl | repository/customer_repo.go (Update) | Repository |
| lgapol01.cbl | service/policy_svc.go (Add) | Service |
| lgapdb01.cbl | repository/policy_repo.go (Create) | Repository |
| lgipol01.cbl | service/policy_svc.go (Get) | Service |
| lgipdb01.cbl | repository/policy_repo.go (FindByNum) | Repository |
| lgupol01.cbl | service/policy_svc.go (Update) | Service |
| lgupdb01.cbl | repository/policy_repo.go (Update) | Repository |
| lgdpol01.cbl | service/policy_svc.go (Delete) | Service |
| lgdpdb01.cbl | repository/policy_repo.go (Delete) | Repository |
| lgsetup.cbl | migrations/ + service/counter_svc.go | Setup |
| lgwebst5.cbl | service/stats_svc.go | Service |
| lgstsq.cbl | Internal logging | Utility |

### 5.3 Dependencies

```go
// go.mod
module github.com/yourorg/genapp

go 1.21

require (
    github.com/rivo/tview v0.0.0-20240101000000-abcdef123456  // Terminal UI
    github.com/gdamore/tcell/v2 v2.7.0                         // Terminal cell library
    github.com/lib/pq v1.10.9                                  // PostgreSQL driver
    github.com/jmoiron/sqlx v1.3.5                            // SQL extensions
    github.com/spf13/viper v1.18.0                            // Configuration
    github.com/rs/zerolog v1.31.0                             // Structured logging
)
```

---

## 6. API Endpoints

### 6.1 REST API Design

```yaml
# Customer Endpoints
GET    /api/v1/customers/{customerNum}      # Inquire customer
POST   /api/v1/customers                     # Add customer
PUT    /api/v1/customers/{customerNum}       # Update customer
DELETE /api/v1/customers/{customerNum}       # Delete customer

# Policy Endpoints
GET    /api/v1/policies/{policyNum}          # Inquire policy
POST   /api/v1/policies                       # Add policy
PUT    /api/v1/policies/{policyNum}          # Update policy
DELETE /api/v1/policies/{policyNum}          # Delete policy

# Policy by Customer
GET    /api/v1/customers/{customerNum}/policies  # List customer policies

# Claim Endpoints
GET    /api/v1/claims/{claimNum}             # Inquire claim
POST   /api/v1/policies/{policyNum}/claims   # Add claim to policy
PUT    /api/v1/claims/{claimNum}             # Update claim

# Statistics
GET    /api/v1/stats                         # Get transaction statistics
```

### 6.2 Request/Response Examples

```json
// POST /api/v1/customers
{
    "first_name": "John",
    "last_name": "Smith",
    "date_of_birth": "1985-03-15",
    "house_name": "Oak House",
    "house_number": "42",
    "postcode": "AB12 3CD",
    "phone_home": "01onal2345678",
    "phone_mobile": "07123456789",
    "email_address": "john.smith@example.com"
}

// Response
{
    "customer_num": "0000000001",
    "first_name": "John",
    "last_name": "Smith",
    ...
}

// POST /api/v1/policies (Motor)
{
    "customer_num": "0000000001",
    "policy_type": "M",
    "issue_date": "2024-01-01",
    "expiry_date": "2025-01-01",
    "details": {
        "make": "Toyota",
        "model": "Camry",
        "value": 25000.00,
        "reg_number": "AB12CDE",
        "colour": "Silver",
        "cc": 2000,
        "manufactured": "2023-06-01",
        "premium": 500.00,
        "accidents": 0
    }
}
```

---

## 7. Migration Strategy

### 7.1 Phased Approach

**Phase 1: Foundation (2 steps)**
- Set up Go project structure with dependencies
- Create PostgreSQL schema and migrations
- Implement domain models

**Phase 2: Data Layer (2 steps)**
- Implement customer repository with CRUD operations
- Implement policy repository with type-specific handling
- Write repository tests

**Phase 3: Service Layer (2 steps)**
- Implement customer service with business logic
- Implement policy service with validation
- Implement counter/sequence service
- Write service tests

**Phase 4: UI Layer (3 steps)**
- Set up OpenTUI application framework
- Implement customer screen (SSMAPC1)
- Implement motor policy screen (SSMAPP1)
- Implement endowment policy screen (SSMAPP2)
- Implement house policy screen (SSMAPP3)

**Phase 5: Integration (1 step)**
- Wire up UI to services
- Add error handling and validation
- End-to-end testing

**Phase 6: Polish (1 step)**
- Add statistics/monitoring
- Documentation
- Performance testing

### 7.2 Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Data model mismatch | Validate schema against original copybooks |
| Business logic gaps | Create test cases from COBOL program flows |
| Screen layout issues | Use fixed 24x80 terminal dimensions |
| Transaction integrity | Use database transactions for multi-table ops |
| Named counter replacement | Use PostgreSQL sequences with atomic operations |

---

## 8. Verification Approach

### 8.1 Testing Strategy

1. **Unit Tests**: Test each repository and service function
2. **Integration Tests**: Test database operations with real PostgreSQL
3. **UI Tests**: Manual testing of OpenTUI screens
4. **Comparison Tests**: Compare Go output with COBOL output for same inputs

### 8.2 Test Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/repository/...
go test ./internal/service/...

# Run integration tests (requires database)
go test -tags=integration ./...

# Lint code
golangci-lint run

# Build application
go build -o bin/genapp ./cmd/genapp
```

### 8.3 Acceptance Criteria

- [ ] All CRUD operations for customers work correctly
- [ ] All CRUD operations for policies (all 4 types) work correctly
- [ ] OpenTUI screens match original BMS layout
- [ ] Field validations match original (numeric, required, length)
- [ ] Counter sequences generate unique IDs
- [ ] Error messages display correctly
- [ ] PF3/Clear key equivalents work (Escape/Ctrl+C)

---

## 9. Files to Create

| File Path | Purpose |
|-----------|---------|
| `cmd/genapp/main.go` | Application entry point |
| `internal/config/config.go` | Configuration management |
| `internal/models/customer.go` | Customer domain model |
| `internal/models/policy.go` | Policy domain models (all types) |
| `internal/models/claim.go` | Claim domain model |
| `internal/repository/db.go` | Database connection setup |
| `internal/repository/customer_repo.go` | Customer data access |
| `internal/repository/policy_repo.go` | Policy data access |
| `internal/repository/claim_repo.go` | Claim data access |
| `internal/service/customer_svc.go` | Customer business logic |
| `internal/service/policy_svc.go` | Policy business logic |
| `internal/service/counter_svc.go` | Sequence/counter service |
| `internal/ui/app.go` | OpenTUI application setup |
| `internal/ui/views/customer.go` | Customer screen |
| `internal/ui/views/motor.go` | Motor policy screen |
| `internal/ui/views/endowment.go` | Endowment policy screen |
| `internal/ui/views/house.go` | House policy screen |
| `internal/ui/components/form.go` | Reusable form component |
| `internal/ui/components/menu.go` | Menu component |
| `migrations/001_initial_schema.sql` | Database schema |
| `go.mod` | Go module definition |
| `README.md` | Project documentation |

---

## 10. Summary

This migration transforms a 30-program COBOL/CICS application into a modern Go application with:

- **~20 Go source files** replacing 30 COBOL programs + 13 copybooks
- **OpenTUI terminal interface** replacing BMS maps
- **PostgreSQL database** replacing DB2 + VSAM
- **Clean architecture** with repository/service/UI layers
- **REST API** for potential future web/mobile clients

The migration preserves all business functionality while modernizing the technology stack for easier maintenance, testing, and deployment.
