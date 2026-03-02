package models

import "time"

// Counter represents a named counter for tracking statistics.
// Replaces CICS Named Counter Server functionality.
type Counter struct {
	Name      string    `json:"name" db:"name"`
	Value     int64     `json:"value" db:"value"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Predefined counter names matching the original COBOL application.
const (
	CounterCustomerAdd      = "customer_add_count"
	CounterCustomerInquiry  = "customer_inq_count"
	CounterCustomerUpdate   = "customer_upd_count"
	CounterPolicyAdd        = "policy_add_count"
	CounterPolicyInquiry    = "policy_inq_count"
	CounterPolicyUpdate     = "policy_upd_count"
	CounterPolicyDelete     = "policy_del_count"
	CounterClaimAdd         = "claim_add_count"
	CounterTotalTransactions = "total_transactions"
)
