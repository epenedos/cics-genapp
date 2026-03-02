// Package models defines the domain models for the GENAPP insurance application.
// These models map to PostgreSQL tables and are derived from the original
// COBOL copybook structures (lgcmarea.cpy, lgpolicy.cpy).
package models

import (
	"database/sql"
	"time"
)

// Customer represents a customer record.
// Maps from DB2-CUSTOMER structure in lgpolicy.cpy.
type Customer struct {
	ID           int64          `json:"id" db:"id"`
	CustomerNum  string         `json:"customer_num" db:"customer_num"`   // PIC 9(10)
	FirstName    sql.NullString `json:"first_name" db:"first_name"`       // PIC X(10)
	LastName     sql.NullString `json:"last_name" db:"last_name"`         // PIC X(20)
	DateOfBirth  sql.NullTime   `json:"date_of_birth" db:"date_of_birth"` // PIC X(10)
	HouseName    sql.NullString `json:"house_name" db:"house_name"`       // PIC X(20)
	HouseNumber  sql.NullString `json:"house_number" db:"house_number"`   // PIC X(4)
	Postcode     sql.NullString `json:"postcode" db:"postcode"`           // PIC X(8)
	PhoneHome    sql.NullString `json:"phone_home" db:"phone_home"`       // PIC X(20)
	PhoneMobile  sql.NullString `json:"phone_mobile" db:"phone_mobile"`   // PIC X(20)
	EmailAddress sql.NullString `json:"email_address" db:"email_address"` // PIC X(100)
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

// CustomerInput represents input data for creating or updating a customer.
// Uses string pointers to distinguish between null and empty values.
type CustomerInput struct {
	FirstName    *string    `json:"first_name"`
	LastName     *string    `json:"last_name"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	HouseName    *string    `json:"house_name"`
	HouseNumber  *string    `json:"house_number"`
	Postcode     *string    `json:"postcode"`
	PhoneHome    *string    `json:"phone_home"`
	PhoneMobile  *string    `json:"phone_mobile"`
	EmailAddress *string    `json:"email_address"`
}

// GetFirstName returns the first name or empty string if null.
func (c *Customer) GetFirstName() string {
	if c.FirstName.Valid {
		return c.FirstName.String
	}
	return ""
}

// GetLastName returns the last name or empty string if null.
func (c *Customer) GetLastName() string {
	if c.LastName.Valid {
		return c.LastName.String
	}
	return ""
}

// GetDateOfBirth returns the date of birth or zero time if null.
func (c *Customer) GetDateOfBirth() time.Time {
	if c.DateOfBirth.Valid {
		return c.DateOfBirth.Time
	}
	return time.Time{}
}

// GetHouseName returns the house name or empty string if null.
func (c *Customer) GetHouseName() string {
	if c.HouseName.Valid {
		return c.HouseName.String
	}
	return ""
}

// GetHouseNumber returns the house number or empty string if null.
func (c *Customer) GetHouseNumber() string {
	if c.HouseNumber.Valid {
		return c.HouseNumber.String
	}
	return ""
}

// GetPostcode returns the postcode or empty string if null.
func (c *Customer) GetPostcode() string {
	if c.Postcode.Valid {
		return c.Postcode.String
	}
	return ""
}

// GetPhoneHome returns the home phone or empty string if null.
func (c *Customer) GetPhoneHome() string {
	if c.PhoneHome.Valid {
		return c.PhoneHome.String
	}
	return ""
}

// GetPhoneMobile returns the mobile phone or empty string if null.
func (c *Customer) GetPhoneMobile() string {
	if c.PhoneMobile.Valid {
		return c.PhoneMobile.String
	}
	return ""
}

// GetEmailAddress returns the email address or empty string if null.
func (c *Customer) GetEmailAddress() string {
	if c.EmailAddress.Valid {
		return c.EmailAddress.String
	}
	return ""
}

// FullName returns the customer's full name.
func (c *Customer) FullName() string {
	first := c.GetFirstName()
	last := c.GetLastName()
	if first == "" && last == "" {
		return ""
	}
	if first == "" {
		return last
	}
	if last == "" {
		return first
	}
	return first + " " + last
}
