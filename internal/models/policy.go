package models

import (
	"database/sql"
	"time"
)

// PolicyType represents the type of insurance policy.
type PolicyType string

const (
	PolicyTypeEndowment  PolicyType = "E"
	PolicyTypeHouse      PolicyType = "H"
	PolicyTypeMotor      PolicyType = "M"
	PolicyTypeCommercial PolicyType = "C"
)

// String returns the policy type as a string.
func (pt PolicyType) String() string {
	return string(pt)
}

// Description returns a human-readable description of the policy type.
func (pt PolicyType) Description() string {
	switch pt {
	case PolicyTypeEndowment:
		return "Endowment"
	case PolicyTypeHouse:
		return "House"
	case PolicyTypeMotor:
		return "Motor"
	case PolicyTypeCommercial:
		return "Commercial"
	default:
		return "Unknown"
	}
}

// IsValid checks if the policy type is one of the allowed values.
func (pt PolicyType) IsValid() bool {
	switch pt {
	case PolicyTypeEndowment, PolicyTypeHouse, PolicyTypeMotor, PolicyTypeCommercial:
		return true
	}
	return false
}

// Policy represents a master policy record.
// Maps from DB2-POLICY structure in lgpolicy.cpy.
type Policy struct {
	ID          int64           `json:"id" db:"id"`
	PolicyNum   string          `json:"policy_num" db:"policy_num"`     // PIC 9(10)
	CustomerNum string          `json:"customer_num" db:"customer_num"` // PIC 9(10)
	PolicyType  PolicyType      `json:"policy_type" db:"policy_type"`   // PIC X
	IssueDate   sql.NullTime    `json:"issue_date" db:"issue_date"`     // PIC X(10)
	ExpiryDate  sql.NullTime    `json:"expiry_date" db:"expiry_date"`   // PIC X(10)
	LastChanged sql.NullTime    `json:"last_changed" db:"last_changed"` // PIC X(26)
	BrokerID    sql.NullString  `json:"broker_id" db:"broker_id"`       // PIC 9(10)
	BrokersRef  sql.NullString  `json:"brokers_ref" db:"brokers_ref"`   // PIC X(10)
	Payment     sql.NullFloat64 `json:"payment" db:"payment"`           // PIC 9(6)
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`

	// Type-specific details (populated when fetching with details)
	Motor      *MotorPolicy      `json:"motor,omitempty" db:"-"`
	Endowment  *EndowmentPolicy  `json:"endowment,omitempty" db:"-"`
	House      *HousePolicy      `json:"house,omitempty" db:"-"`
	Commercial *CommercialPolicy `json:"commercial,omitempty" db:"-"`
}

// PolicyInput represents input data for creating or updating a policy.
type PolicyInput struct {
	CustomerNum string      `json:"customer_num"`
	PolicyType  PolicyType  `json:"policy_type"`
	IssueDate   *time.Time  `json:"issue_date"`
	ExpiryDate  *time.Time  `json:"expiry_date"`
	BrokerID    *string     `json:"broker_id"`
	BrokersRef  *string     `json:"brokers_ref"`
	Payment     *float64    `json:"payment"`
	Motor       *MotorInput `json:"motor,omitempty"`
	Endowment   *EndowmentInput `json:"endowment,omitempty"`
	House       *HouseInput `json:"house,omitempty"`
	Commercial  *CommercialInput `json:"commercial,omitempty"`
}

// GetIssueDate returns the issue date or zero time if null.
func (p *Policy) GetIssueDate() time.Time {
	if p.IssueDate.Valid {
		return p.IssueDate.Time
	}
	return time.Time{}
}

// GetExpiryDate returns the expiry date or zero time if null.
func (p *Policy) GetExpiryDate() time.Time {
	if p.ExpiryDate.Valid {
		return p.ExpiryDate.Time
	}
	return time.Time{}
}

// GetPayment returns the payment amount or 0 if null.
func (p *Policy) GetPayment() float64 {
	if p.Payment.Valid {
		return p.Payment.Float64
	}
	return 0
}

// MotorPolicy represents motor-specific policy details.
// Maps from DB2-MOTOR structure in lgpolicy.cpy.
type MotorPolicy struct {
	ID           int64           `json:"id" db:"id"`
	PolicyNum    string          `json:"policy_num" db:"policy_num"`       // PIC 9(10)
	Make         sql.NullString  `json:"make" db:"make"`                   // PIC X(15)
	Model        sql.NullString  `json:"model" db:"model"`                 // PIC X(15)
	Value        sql.NullFloat64 `json:"value" db:"value"`                 // PIC 9(6)
	RegNumber    sql.NullString  `json:"reg_number" db:"reg_number"`       // PIC X(7)
	Colour       sql.NullString  `json:"colour" db:"colour"`               // PIC X(8)
	CC           sql.NullInt64   `json:"cc" db:"cc"`                       // PIC 9(4)
	Manufactured sql.NullTime    `json:"manufactured" db:"manufactured"`   // PIC X(10)
	Premium      sql.NullFloat64 `json:"premium" db:"premium"`             // PIC 9(6)
	Accidents    sql.NullInt64   `json:"accidents" db:"accidents"`         // PIC 9(6)
}

// MotorInput represents input data for creating or updating motor policy details.
type MotorInput struct {
	Make         *string    `json:"make"`
	Model        *string    `json:"model"`
	Value        *float64   `json:"value"`
	RegNumber    *string    `json:"reg_number"`
	Colour       *string    `json:"colour"`
	CC           *int       `json:"cc"`
	Manufactured *time.Time `json:"manufactured"`
	Premium      *float64   `json:"premium"`
	Accidents    *int       `json:"accidents"`
}

// GetMake returns the car make or empty string if null.
func (m *MotorPolicy) GetMake() string {
	if m.Make.Valid {
		return m.Make.String
	}
	return ""
}

// GetModel returns the car model or empty string if null.
func (m *MotorPolicy) GetModel() string {
	if m.Model.Valid {
		return m.Model.String
	}
	return ""
}

// GetValue returns the car value or 0 if null.
func (m *MotorPolicy) GetValue() float64 {
	if m.Value.Valid {
		return m.Value.Float64
	}
	return 0
}

// GetRegNumber returns the registration number or empty string if null.
func (m *MotorPolicy) GetRegNumber() string {
	if m.RegNumber.Valid {
		return m.RegNumber.String
	}
	return ""
}

// GetColour returns the car colour or empty string if null.
func (m *MotorPolicy) GetColour() string {
	if m.Colour.Valid {
		return m.Colour.String
	}
	return ""
}

// GetCC returns the engine cc or 0 if null.
func (m *MotorPolicy) GetCC() int {
	if m.CC.Valid {
		return int(m.CC.Int64)
	}
	return 0
}

// GetManufactured returns the manufacture date or zero time if null.
func (m *MotorPolicy) GetManufactured() time.Time {
	if m.Manufactured.Valid {
		return m.Manufactured.Time
	}
	return time.Time{}
}

// GetPremium returns the premium or 0 if null.
func (m *MotorPolicy) GetPremium() float64 {
	if m.Premium.Valid {
		return m.Premium.Float64
	}
	return 0
}

// GetAccidents returns the accident count or 0 if null.
func (m *MotorPolicy) GetAccidents() int {
	if m.Accidents.Valid {
		return int(m.Accidents.Int64)
	}
	return 0
}

// EndowmentPolicy represents endowment-specific policy details.
// Maps from DB2-ENDOWMENT structure in lgpolicy.cpy.
type EndowmentPolicy struct {
	ID          int64           `json:"id" db:"id"`
	PolicyNum   string          `json:"policy_num" db:"policy_num"`     // PIC 9(10)
	WithProfits sql.NullBool    `json:"with_profits" db:"with_profits"` // PIC X (Y/N)
	Equities    sql.NullBool    `json:"equities" db:"equities"`         // PIC X (Y/N)
	ManagedFund sql.NullBool    `json:"managed_fund" db:"managed_fund"` // PIC X (Y/N)
	FundName    sql.NullString  `json:"fund_name" db:"fund_name"`       // PIC X(10)
	Term        sql.NullInt64   `json:"term" db:"term"`                 // PIC 9(2)
	SumAssured  sql.NullFloat64 `json:"sum_assured" db:"sum_assured"`   // PIC 9(6)
	LifeAssured sql.NullString  `json:"life_assured" db:"life_assured"` // PIC X(31)
}

// EndowmentInput represents input data for creating or updating endowment policy details.
type EndowmentInput struct {
	WithProfits *bool    `json:"with_profits"`
	Equities    *bool    `json:"equities"`
	ManagedFund *bool    `json:"managed_fund"`
	FundName    *string  `json:"fund_name"`
	Term        *int     `json:"term"`
	SumAssured  *float64 `json:"sum_assured"`
	LifeAssured *string  `json:"life_assured"`
}

// GetWithProfits returns whether with profits is enabled.
func (e *EndowmentPolicy) GetWithProfits() bool {
	if e.WithProfits.Valid {
		return e.WithProfits.Bool
	}
	return false
}

// GetEquities returns whether equities are enabled.
func (e *EndowmentPolicy) GetEquities() bool {
	if e.Equities.Valid {
		return e.Equities.Bool
	}
	return false
}

// GetManagedFund returns whether managed fund is enabled.
func (e *EndowmentPolicy) GetManagedFund() bool {
	if e.ManagedFund.Valid {
		return e.ManagedFund.Bool
	}
	return false
}

// GetFundName returns the fund name or empty string if null.
func (e *EndowmentPolicy) GetFundName() string {
	if e.FundName.Valid {
		return e.FundName.String
	}
	return ""
}

// GetTerm returns the term in years or 0 if null.
func (e *EndowmentPolicy) GetTerm() int {
	if e.Term.Valid {
		return int(e.Term.Int64)
	}
	return 0
}

// GetSumAssured returns the sum assured or 0 if null.
func (e *EndowmentPolicy) GetSumAssured() float64 {
	if e.SumAssured.Valid {
		return e.SumAssured.Float64
	}
	return 0
}

// GetLifeAssured returns the life assured name or empty string if null.
func (e *EndowmentPolicy) GetLifeAssured() string {
	if e.LifeAssured.Valid {
		return e.LifeAssured.String
	}
	return ""
}

// HousePolicy represents house-specific policy details.
// Maps from DB2-HOUSE structure in lgpolicy.cpy.
type HousePolicy struct {
	ID           int64           `json:"id" db:"id"`
	PolicyNum    string          `json:"policy_num" db:"policy_num"`       // PIC 9(10)
	PropertyType sql.NullString  `json:"property_type" db:"property_type"` // PIC X(15)
	Bedrooms     sql.NullInt64   `json:"bedrooms" db:"bedrooms"`           // PIC 9(3)
	Value        sql.NullFloat64 `json:"value" db:"value"`                 // PIC 9(8)
	HouseName    sql.NullString  `json:"house_name" db:"house_name"`       // PIC X(20)
	HouseNumber  sql.NullString  `json:"house_number" db:"house_number"`   // PIC X(4)
	Postcode     sql.NullString  `json:"postcode" db:"postcode"`           // PIC X(8)
}

// HouseInput represents input data for creating or updating house policy details.
type HouseInput struct {
	PropertyType *string  `json:"property_type"`
	Bedrooms     *int     `json:"bedrooms"`
	Value        *float64 `json:"value"`
	HouseName    *string  `json:"house_name"`
	HouseNumber  *string  `json:"house_number"`
	Postcode     *string  `json:"postcode"`
}

// GetPropertyType returns the property type or empty string if null.
func (h *HousePolicy) GetPropertyType() string {
	if h.PropertyType.Valid {
		return h.PropertyType.String
	}
	return ""
}

// GetBedrooms returns the number of bedrooms or 0 if null.
func (h *HousePolicy) GetBedrooms() int {
	if h.Bedrooms.Valid {
		return int(h.Bedrooms.Int64)
	}
	return 0
}

// GetValue returns the property value or 0 if null.
func (h *HousePolicy) GetValue() float64 {
	if h.Value.Valid {
		return h.Value.Float64
	}
	return 0
}

// GetHouseName returns the house name or empty string if null.
func (h *HousePolicy) GetHouseName() string {
	if h.HouseName.Valid {
		return h.HouseName.String
	}
	return ""
}

// GetHouseNumber returns the house number or empty string if null.
func (h *HousePolicy) GetHouseNumber() string {
	if h.HouseNumber.Valid {
		return h.HouseNumber.String
	}
	return ""
}

// GetPostcode returns the postcode or empty string if null.
func (h *HousePolicy) GetPostcode() string {
	if h.Postcode.Valid {
		return h.Postcode.String
	}
	return ""
}

// CommercialPolicy represents commercial-specific policy details.
// Maps from DB2-COMMERCIAL structure in lgpolicy.cpy.
type CommercialPolicy struct {
	ID             int64           `json:"id" db:"id"`
	PolicyNum      string          `json:"policy_num" db:"policy_num"`           // PIC 9(10)
	Address        sql.NullString  `json:"address" db:"address"`                 // PIC X(255)
	Postcode       sql.NullString  `json:"postcode" db:"postcode"`               // PIC X(8)
	Latitude       sql.NullString  `json:"latitude" db:"latitude"`               // PIC X(11)
	Longitude      sql.NullString  `json:"longitude" db:"longitude"`             // PIC X(11)
	Customer       sql.NullString  `json:"customer" db:"customer"`               // PIC X(255)
	PropertyType   sql.NullString  `json:"property_type" db:"property_type"`     // PIC X(255)
	FirePeril      sql.NullInt64   `json:"fire_peril" db:"fire_peril"`           // PIC 9(4)
	FirePremium    sql.NullFloat64 `json:"fire_premium" db:"fire_premium"`       // PIC 9(8)
	CrimePeril     sql.NullInt64   `json:"crime_peril" db:"crime_peril"`         // PIC 9(4)
	CrimePremium   sql.NullFloat64 `json:"crime_premium" db:"crime_premium"`     // PIC 9(8)
	FloodPeril     sql.NullInt64   `json:"flood_peril" db:"flood_peril"`         // PIC 9(4)
	FloodPremium   sql.NullFloat64 `json:"flood_premium" db:"flood_premium"`     // PIC 9(8)
	WeatherPeril   sql.NullInt64   `json:"weather_peril" db:"weather_peril"`     // PIC 9(4)
	WeatherPremium sql.NullFloat64 `json:"weather_premium" db:"weather_premium"` // PIC 9(8)
	Status         sql.NullInt64   `json:"status" db:"status"`                   // PIC 9(4)
	RejectReason   sql.NullString  `json:"reject_reason" db:"reject_reason"`     // PIC X(255)
}

// CommercialInput represents input data for creating or updating commercial policy details.
type CommercialInput struct {
	Address        *string  `json:"address"`
	Postcode       *string  `json:"postcode"`
	Latitude       *string  `json:"latitude"`
	Longitude      *string  `json:"longitude"`
	Customer       *string  `json:"customer"`
	PropertyType   *string  `json:"property_type"`
	FirePeril      *int     `json:"fire_peril"`
	FirePremium    *float64 `json:"fire_premium"`
	CrimePeril     *int     `json:"crime_peril"`
	CrimePremium   *float64 `json:"crime_premium"`
	FloodPeril     *int     `json:"flood_peril"`
	FloodPremium   *float64 `json:"flood_premium"`
	WeatherPeril   *int     `json:"weather_peril"`
	WeatherPremium *float64 `json:"weather_premium"`
	Status         *int     `json:"status"`
	RejectReason   *string  `json:"reject_reason"`
}

// GetAddress returns the address or empty string if null.
func (c *CommercialPolicy) GetAddress() string {
	if c.Address.Valid {
		return c.Address.String
	}
	return ""
}

// GetPostcode returns the postcode or empty string if null.
func (c *CommercialPolicy) GetPostcode() string {
	if c.Postcode.Valid {
		return c.Postcode.String
	}
	return ""
}

// GetLatitude returns the latitude or empty string if null.
func (c *CommercialPolicy) GetLatitude() string {
	if c.Latitude.Valid {
		return c.Latitude.String
	}
	return ""
}

// GetLongitude returns the longitude or empty string if null.
func (c *CommercialPolicy) GetLongitude() string {
	if c.Longitude.Valid {
		return c.Longitude.String
	}
	return ""
}

// GetCustomer returns the customer or empty string if null.
func (c *CommercialPolicy) GetCustomer() string {
	if c.Customer.Valid {
		return c.Customer.String
	}
	return ""
}

// GetPropertyType returns the property type or empty string if null.
func (c *CommercialPolicy) GetPropertyType() string {
	if c.PropertyType.Valid {
		return c.PropertyType.String
	}
	return ""
}

// GetFirePeril returns the fire peril level or 0 if null.
func (c *CommercialPolicy) GetFirePeril() int {
	if c.FirePeril.Valid {
		return int(c.FirePeril.Int64)
	}
	return 0
}

// GetFirePremium returns the fire premium or 0 if null.
func (c *CommercialPolicy) GetFirePremium() float64 {
	if c.FirePremium.Valid {
		return c.FirePremium.Float64
	}
	return 0
}

// GetCrimePeril returns the crime peril level or 0 if null.
func (c *CommercialPolicy) GetCrimePeril() int {
	if c.CrimePeril.Valid {
		return int(c.CrimePeril.Int64)
	}
	return 0
}

// GetCrimePremium returns the crime premium or 0 if null.
func (c *CommercialPolicy) GetCrimePremium() float64 {
	if c.CrimePremium.Valid {
		return c.CrimePremium.Float64
	}
	return 0
}

// GetFloodPeril returns the flood peril level or 0 if null.
func (c *CommercialPolicy) GetFloodPeril() int {
	if c.FloodPeril.Valid {
		return int(c.FloodPeril.Int64)
	}
	return 0
}

// GetFloodPremium returns the flood premium or 0 if null.
func (c *CommercialPolicy) GetFloodPremium() float64 {
	if c.FloodPremium.Valid {
		return c.FloodPremium.Float64
	}
	return 0
}

// GetWeatherPeril returns the weather peril level or 0 if null.
func (c *CommercialPolicy) GetWeatherPeril() int {
	if c.WeatherPeril.Valid {
		return int(c.WeatherPeril.Int64)
	}
	return 0
}

// GetWeatherPremium returns the weather premium or 0 if null.
func (c *CommercialPolicy) GetWeatherPremium() float64 {
	if c.WeatherPremium.Valid {
		return c.WeatherPremium.Float64
	}
	return 0
}

// GetStatus returns the status or 0 if null.
func (c *CommercialPolicy) GetStatus() int {
	if c.Status.Valid {
		return int(c.Status.Int64)
	}
	return 0
}

// GetRejectReason returns the reject reason or empty string if null.
func (c *CommercialPolicy) GetRejectReason() string {
	if c.RejectReason.Valid {
		return c.RejectReason.String
	}
	return ""
}

// TotalPremium calculates the total premium for a commercial policy.
func (c *CommercialPolicy) TotalPremium() float64 {
	return c.GetFirePremium() + c.GetCrimePremium() + c.GetFloodPremium() + c.GetWeatherPremium()
}
