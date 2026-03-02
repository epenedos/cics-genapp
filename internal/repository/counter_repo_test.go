package repository

import (
	"context"
	"testing"

	"github.com/cicsdev/genapp/internal/models"
)

func TestCounterRepository_Increment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "test_counter"

	// First increment (creates counter)
	value, err := repo.Increment(ctx, counterName, 1)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if value != 1 {
		t.Errorf("Expected value 1, got %d", value)
	}

	// Second increment
	value, err = repo.Increment(ctx, counterName, 5)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if value != 6 {
		t.Errorf("Expected value 6, got %d", value)
	}
}

func TestCounterRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "get_test_counter"

	// Set a value first
	err := repo.Set(ctx, counterName, 42)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the counter
	counter, err := repo.Get(ctx, counterName)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if counter.Name != counterName {
		t.Errorf("Expected name '%s', got '%s'", counterName, counter.Name)
	}
	if counter.Value != 42 {
		t.Errorf("Expected value 42, got %d", counter.Value)
	}
}

func TestCounterRepository_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, "non_existent_counter")
	if err != ErrCounterNotFound {
		t.Errorf("Expected ErrCounterNotFound, got %v", err)
	}
}

func TestCounterRepository_GetValue(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "value_test_counter"

	// Set a value
	err := repo.Set(ctx, counterName, 100)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get value only
	value, err := repo.GetValue(ctx, counterName)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}

	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}
}

func TestCounterRepository_Set(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "set_test_counter"

	// Set initial value
	err := repo.Set(ctx, counterName, 50)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify
	value, err := repo.GetValue(ctx, counterName)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if value != 50 {
		t.Errorf("Expected value 50, got %d", value)
	}

	// Update value
	err = repo.Set(ctx, counterName, 200)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify update
	value, err = repo.GetValue(ctx, counterName)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if value != 200 {
		t.Errorf("Expected value 200, got %d", value)
	}
}

func TestCounterRepository_Reset(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "reset_test_counter"

	// Set a value
	err := repo.Set(ctx, counterName, 999)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Reset
	err = repo.Reset(ctx, counterName)
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	// Verify reset
	value, err := repo.GetValue(ctx, counterName)
	if err != nil {
		t.Fatalf("GetValue failed: %v", err)
	}
	if value != 0 {
		t.Errorf("Expected value 0 after reset, got %d", value)
	}
}

func TestCounterRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "delete_test_counter"

	// Create counter
	err := repo.Set(ctx, counterName, 10)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete
	err = repo.Delete(ctx, counterName)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.Get(ctx, counterName)
	if err != ErrCounterNotFound {
		t.Errorf("Expected ErrCounterNotFound after delete, got %v", err)
	}
}

func TestCounterRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "non_existent_counter")
	if err != ErrCounterNotFound {
		t.Errorf("Expected ErrCounterNotFound, got %v", err)
	}
}

func TestCounterRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	// Note: There are pre-existing counters from the migration seed data
	// Get initial count
	initial, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	initialCount := len(initial)

	// Create some test counters
	testCounters := []string{"list_test_a", "list_test_b", "list_test_c"}
	for i, name := range testCounters {
		err := repo.Set(ctx, name, int64(i+1))
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
	}

	// List
	counters, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	expectedCount := initialCount + len(testCounters)
	if len(counters) != expectedCount {
		t.Errorf("Expected %d counters, got %d", expectedCount, len(counters))
	}
}

func TestCounterRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counterName := "exists_test_counter"

	// Check non-existent
	exists, err := repo.Exists(ctx, counterName)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected counter to not exist")
	}

	// Create counter
	err = repo.Set(ctx, counterName, 1)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Check exists
	exists, err = repo.Exists(ctx, counterName)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected counter to exist")
	}
}

func TestCounterRepository_GetMultiple(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	// Create counters
	testCounters := map[string]int64{
		"multi_test_a": 10,
		"multi_test_b": 20,
		"multi_test_c": 30,
	}

	for name, value := range testCounters {
		err := repo.Set(ctx, name, value)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
	}

	// Get multiple
	names := []string{"multi_test_a", "multi_test_c"}
	counters, err := repo.GetMultiple(ctx, names)
	if err != nil {
		t.Fatalf("GetMultiple failed: %v", err)
	}

	if len(counters) != 2 {
		t.Errorf("Expected 2 counters, got %d", len(counters))
	}

	// Verify values
	counterMap := make(map[string]*models.Counter)
	for _, c := range counters {
		counterMap[c.Name] = c
	}

	if c, ok := counterMap["multi_test_a"]; ok {
		if c.Value != 10 {
			t.Errorf("Expected multi_test_a value 10, got %d", c.Value)
		}
	} else {
		t.Error("Expected multi_test_a in results")
	}

	if c, ok := counterMap["multi_test_c"]; ok {
		if c.Value != 30 {
			t.Errorf("Expected multi_test_c value 30, got %d", c.Value)
		}
	} else {
		t.Error("Expected multi_test_c in results")
	}
}

func TestCounterRepository_GetMultiple_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCounterRepository(db)
	ctx := context.Background()

	counters, err := repo.GetMultiple(ctx, []string{})
	if err != nil {
		t.Fatalf("GetMultiple failed: %v", err)
	}

	if len(counters) != 0 {
		t.Errorf("Expected 0 counters for empty input, got %d", len(counters))
	}
}

func TestPredefinedCounterNames(t *testing.T) {
	// Verify predefined counter names exist
	expectedCounters := []string{
		models.CounterCustomerAdd,
		models.CounterCustomerInquiry,
		models.CounterCustomerUpdate,
		models.CounterPolicyAdd,
		models.CounterPolicyInquiry,
		models.CounterPolicyUpdate,
		models.CounterPolicyDelete,
		models.CounterClaimAdd,
		models.CounterTotalTransactions,
	}

	for _, name := range expectedCounters {
		if name == "" {
			t.Error("Found empty predefined counter name")
		}
	}

	// Verify uniqueness
	seen := make(map[string]bool)
	for _, name := range expectedCounters {
		if seen[name] {
			t.Errorf("Duplicate counter name: %s", name)
		}
		seen[name] = true
	}
}
