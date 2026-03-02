package repository

import (
	"context"
	"os"
	"testing"
)

// setupTestDB creates a database connection for testing.
// It uses the DATABASE_URL environment variable if set, otherwise
// uses a default local PostgreSQL connection.
//
// Tests that use this function are integration tests and require
// a running PostgreSQL instance with the migrations applied.
func setupTestDB(t *testing.T) *DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=genapp_test sslmode=disable"
	}

	db, err := NewDBFromDSN(dsn)
	if err != nil {
		t.Skipf("Skipping test: could not connect to database: %v", err)
	}

	// Verify connection
	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		t.Skipf("Skipping test: database ping failed: %v", err)
	}

	// Clean up test data before each test
	cleanupTestData(t, db)

	return db
}

// cleanupTestData removes all test data from the database.
// This ensures tests start with a clean slate.
func cleanupTestData(t *testing.T, db *DB) {
	t.Helper()

	ctx := context.Background()

	// Delete in order respecting foreign keys
	tables := []string{
		"claims",
		"motor_policies",
		"endowment_policies",
		"house_policies",
		"commercial_policies",
		"policies",
		"customers",
	}

	for _, table := range tables {
		_, err := db.ExecContext(ctx, "DELETE FROM "+table)
		if err != nil {
			t.Logf("Warning: could not clean table %s: %v", table, err)
		}
	}

	// Reset sequences
	sequences := []string{
		"customer_num_seq",
		"policy_num_seq",
		"claim_num_seq",
	}

	for _, seq := range sequences {
		_, err := db.ExecContext(ctx, "ALTER SEQUENCE "+seq+" RESTART WITH 1000000001")
		if err != nil {
			t.Logf("Warning: could not reset sequence %s: %v", seq, err)
		}
	}
}

// runWithTestDB executes a test function with a database connection.
// This is useful for running subtests with a shared database.
func runWithTestDB(t *testing.T, name string, fn func(t *testing.T, db *DB)) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run(name, func(t *testing.T) {
		fn(t, db)
	})
}
