package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrCounterNotFound is returned when a counter is not found.
var ErrCounterNotFound = errors.New("counter not found")

// CounterRepository provides data access operations for named counters.
// This replaces CICS Named Counter Server functionality.
type CounterRepository struct {
	db *DB
}

// NewCounterRepository creates a new CounterRepository.
func NewCounterRepository(db *DB) *CounterRepository {
	return &CounterRepository{db: db}
}

// Get retrieves a counter by name.
func (r *CounterRepository) Get(ctx context.Context, name string) (*models.Counter, error) {
	return r.GetTx(ctx, r.db, name)
}

// GetTx retrieves a counter within a transaction.
func (r *CounterRepository) GetTx(ctx context.Context, tx Transactional, name string) (*models.Counter, error) {
	query := `SELECT name, value, updated_at FROM counters WHERE name = $1`

	var counter models.Counter
	err := sqlx.GetContext(ctx, tx, &counter, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCounterNotFound
		}
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	return &counter, nil
}

// GetValue retrieves just the value of a counter.
func (r *CounterRepository) GetValue(ctx context.Context, name string) (int64, error) {
	counter, err := r.Get(ctx, name)
	if err != nil {
		return 0, err
	}
	return counter.Value, nil
}

// Increment atomically increments a counter and returns the new value.
// This is equivalent to CICS GET COUNTER with UPDATE option.
func (r *CounterRepository) Increment(ctx context.Context, name string, delta int64) (int64, error) {
	return r.IncrementTx(ctx, r.db, name, delta)
}

// IncrementTx increments a counter within a transaction.
func (r *CounterRepository) IncrementTx(ctx context.Context, tx Transactional, name string, delta int64) (int64, error) {
	// Use the PostgreSQL function for atomic increment
	query := `SELECT increment_counter($1, $2)`

	var newValue int64
	err := sqlx.GetContext(ctx, tx, &newValue, query, name, delta)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return newValue, nil
}

// Set sets a counter to a specific value.
func (r *CounterRepository) Set(ctx context.Context, name string, value int64) error {
	return r.SetTx(ctx, r.db, name, value)
}

// SetTx sets a counter within a transaction.
func (r *CounterRepository) SetTx(ctx context.Context, tx Transactional, name string, value int64) error {
	query := `
		INSERT INTO counters (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE
		SET value = $2, updated_at = CURRENT_TIMESTAMP
	`

	_, err := tx.ExecContext(ctx, query, name, value)
	if err != nil {
		return fmt.Errorf("failed to set counter: %w", err)
	}

	return nil
}

// Reset resets a counter to zero.
func (r *CounterRepository) Reset(ctx context.Context, name string) error {
	return r.Set(ctx, name, 0)
}

// Delete removes a counter.
func (r *CounterRepository) Delete(ctx context.Context, name string) error {
	return r.DeleteTx(ctx, r.db, name)
}

// DeleteTx removes a counter within a transaction.
func (r *CounterRepository) DeleteTx(ctx context.Context, tx Transactional, name string) error {
	query := `DELETE FROM counters WHERE name = $1`

	result, err := tx.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete counter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCounterNotFound
	}

	return nil
}

// List retrieves all counters.
func (r *CounterRepository) List(ctx context.Context) ([]*models.Counter, error) {
	query := `SELECT name, value, updated_at FROM counters ORDER BY name`

	var counters []*models.Counter
	err := r.db.SelectContext(ctx, &counters, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list counters: %w", err)
	}

	return counters, nil
}

// Exists checks if a counter exists by name.
func (r *CounterRepository) Exists(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM counters WHERE name = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, name)
	if err != nil {
		return false, fmt.Errorf("failed to check counter existence: %w", err)
	}

	return exists, nil
}

// GetMultiple retrieves multiple counters by name.
func (r *CounterRepository) GetMultiple(ctx context.Context, names []string) ([]*models.Counter, error) {
	if len(names) == 0 {
		return []*models.Counter{}, nil
	}

	query, args, err := sqlx.In(`SELECT name, value, updated_at FROM counters WHERE name IN (?) ORDER BY name`, names)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = r.db.Rebind(query)

	var counters []*models.Counter
	err = r.db.SelectContext(ctx, &counters, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters: %w", err)
	}

	return counters, nil
}

// IncrementAll increments multiple counters atomically.
func (r *CounterRepository) IncrementAll(ctx context.Context, increments map[string]int64) error {
	return r.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		for name, delta := range increments {
			_, err := r.IncrementTx(ctx, tx, name, delta)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
