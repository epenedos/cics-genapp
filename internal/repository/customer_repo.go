package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrCustomerNotFound is returned when a customer is not found.
var ErrCustomerNotFound = errors.New("customer not found")

// CustomerRepository provides data access operations for customers.
type CustomerRepository struct {
	db *DB
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(db *DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create inserts a new customer and returns the created customer with the generated number.
// Equivalent to LGACDB01 COBOL program.
func (r *CustomerRepository) Create(ctx context.Context, input *models.CustomerInput) (*models.Customer, error) {
	return r.CreateTx(ctx, r.db, input)
}

// CreateTx inserts a new customer within a transaction.
func (r *CustomerRepository) CreateTx(ctx context.Context, tx Transactional, input *models.CustomerInput) (*models.Customer, error) {
	query := `
		INSERT INTO customers (
			customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address
		) VALUES (
			next_customer_num(), $1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
	`

	var customer models.Customer
	err := sqlx.GetContext(ctx, tx, &customer, query,
		toNullString(input.FirstName),
		toNullString(input.LastName),
		toNullTime(input.DateOfBirth),
		toNullString(input.HouseName),
		toNullString(input.HouseNumber),
		toNullString(input.Postcode),
		toNullString(input.PhoneHome),
		toNullString(input.PhoneMobile),
		toNullString(input.EmailAddress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &customer, nil
}

// FindByNum retrieves a customer by their customer number.
// Equivalent to LGICDB01 COBOL program.
func (r *CustomerRepository) FindByNum(ctx context.Context, customerNum string) (*models.Customer, error) {
	return r.FindByNumTx(ctx, r.db, customerNum)
}

// FindByNumTx retrieves a customer by number within a transaction.
func (r *CustomerRepository) FindByNumTx(ctx context.Context, tx Transactional, customerNum string) (*models.Customer, error) {
	query := `
		SELECT id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
		FROM customers
		WHERE customer_num = $1
	`

	var customer models.Customer
	err := sqlx.GetContext(ctx, tx, &customer, query, customerNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return &customer, nil
}

// FindByID retrieves a customer by their internal ID.
func (r *CustomerRepository) FindByID(ctx context.Context, id int64) (*models.Customer, error) {
	query := `
		SELECT id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	var customer models.Customer
	err := r.db.GetContext(ctx, &customer, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return &customer, nil
}

// Update updates an existing customer.
// Equivalent to LGUCDB01 COBOL program.
func (r *CustomerRepository) Update(ctx context.Context, customerNum string, input *models.CustomerInput) (*models.Customer, error) {
	return r.UpdateTx(ctx, r.db, customerNum, input)
}

// UpdateTx updates a customer within a transaction.
func (r *CustomerRepository) UpdateTx(ctx context.Context, tx Transactional, customerNum string, input *models.CustomerInput) (*models.Customer, error) {
	query := `
		UPDATE customers SET
			first_name = COALESCE($2, first_name),
			last_name = COALESCE($3, last_name),
			date_of_birth = COALESCE($4, date_of_birth),
			house_name = COALESCE($5, house_name),
			house_number = COALESCE($6, house_number),
			postcode = COALESCE($7, postcode),
			phone_home = COALESCE($8, phone_home),
			phone_mobile = COALESCE($9, phone_mobile),
			email_address = COALESCE($10, email_address)
		WHERE customer_num = $1
		RETURNING id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
	`

	var customer models.Customer
	err := sqlx.GetContext(ctx, tx, &customer, query,
		customerNum,
		toNullString(input.FirstName),
		toNullString(input.LastName),
		toNullTime(input.DateOfBirth),
		toNullString(input.HouseName),
		toNullString(input.HouseNumber),
		toNullString(input.Postcode),
		toNullString(input.PhoneHome),
		toNullString(input.PhoneMobile),
		toNullString(input.EmailAddress),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return &customer, nil
}

// Delete removes a customer by their customer number.
// Note: This will cascade delete associated policies.
func (r *CustomerRepository) Delete(ctx context.Context, customerNum string) error {
	return r.DeleteTx(ctx, r.db, customerNum)
}

// DeleteTx removes a customer within a transaction.
func (r *CustomerRepository) DeleteTx(ctx context.Context, tx Transactional, customerNum string) error {
	query := `DELETE FROM customers WHERE customer_num = $1`

	result, err := tx.ExecContext(ctx, query, customerNum)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCustomerNotFound
	}

	return nil
}

// List retrieves a paginated list of customers.
func (r *CustomerRepository) List(ctx context.Context, offset, limit int) ([]*models.Customer, error) {
	query := `
		SELECT id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
		FROM customers
		ORDER BY customer_num
		LIMIT $1 OFFSET $2
	`

	var customers []*models.Customer
	err := r.db.SelectContext(ctx, &customers, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	return customers, nil
}

// Count returns the total number of customers.
func (r *CustomerRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM customers`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}

	return count, nil
}

// SearchByLastName finds customers by last name (partial match).
func (r *CustomerRepository) SearchByLastName(ctx context.Context, lastName string, limit int) ([]*models.Customer, error) {
	query := `
		SELECT id, customer_num, first_name, last_name, date_of_birth,
			house_name, house_number, postcode,
			phone_home, phone_mobile, email_address,
			created_at, updated_at
		FROM customers
		WHERE last_name ILIKE $1
		ORDER BY last_name, first_name
		LIMIT $2
	`

	var customers []*models.Customer
	err := r.db.SelectContext(ctx, &customers, query, "%"+lastName+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, nil
}

// Exists checks if a customer exists by their customer number.
func (r *CustomerRepository) Exists(ctx context.Context, customerNum string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE customer_num = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, customerNum)
	if err != nil {
		return false, fmt.Errorf("failed to check customer existence: %w", err)
	}

	return exists, nil
}
