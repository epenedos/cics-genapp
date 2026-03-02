package models

import (
	"database/sql"
	"time"
)

// Claim represents an insurance claim record.
// Maps from DB2-CLAIM structure in lgpolicy.cpy.
type Claim struct {
	ID           int64           `json:"id" db:"id"`
	ClaimNum     string          `json:"claim_num" db:"claim_num"`       // PIC 9(10)
	PolicyNum    string          `json:"policy_num" db:"policy_num"`     // PIC 9(10)
	ClaimDate    sql.NullTime    `json:"claim_date" db:"claim_date"`     // PIC X(10)
	Paid         sql.NullFloat64 `json:"paid" db:"paid"`                 // PIC 9(8)
	Value        sql.NullFloat64 `json:"value" db:"value"`               // PIC 9(8)
	Cause        sql.NullString  `json:"cause" db:"cause"`               // PIC X(255)
	Observations sql.NullString  `json:"observations" db:"observations"` // PIC X(255)
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// ClaimInput represents input data for creating or updating a claim.
type ClaimInput struct {
	PolicyNum    string     `json:"policy_num"`
	ClaimDate    *time.Time `json:"claim_date"`
	Paid         *float64   `json:"paid"`
	Value        *float64   `json:"value"`
	Cause        *string    `json:"cause"`
	Observations *string    `json:"observations"`
}

// GetClaimDate returns the claim date or zero time if null.
func (c *Claim) GetClaimDate() time.Time {
	if c.ClaimDate.Valid {
		return c.ClaimDate.Time
	}
	return time.Time{}
}

// GetPaid returns the paid amount or 0 if null.
func (c *Claim) GetPaid() float64 {
	if c.Paid.Valid {
		return c.Paid.Float64
	}
	return 0
}

// GetValue returns the claim value or 0 if null.
func (c *Claim) GetValue() float64 {
	if c.Value.Valid {
		return c.Value.Float64
	}
	return 0
}

// GetCause returns the claim cause or empty string if null.
func (c *Claim) GetCause() string {
	if c.Cause.Valid {
		return c.Cause.String
	}
	return ""
}

// GetObservations returns the observations or empty string if null.
func (c *Claim) GetObservations() string {
	if c.Observations.Valid {
		return c.Observations.String
	}
	return ""
}

// Outstanding returns the outstanding amount (value - paid).
func (c *Claim) Outstanding() float64 {
	return c.GetValue() - c.GetPaid()
}
