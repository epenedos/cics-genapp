package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

func TestCounterService_GetCounter(t *testing.T) {
	svc := &CounterService{}

	// Test empty name validation
	_, err := svc.GetCounter(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if !errors.Is(err, ErrInvalidCounter) {
		t.Errorf("Expected ErrInvalidCounter, got %v", err)
	}
}

func TestCounterService_IncrementCounter(t *testing.T) {
	svc := &CounterService{}

	// Test empty name validation
	_, err := svc.IncrementCounter(context.Background(), "", 1)
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if !errors.Is(err, ErrInvalidCounter) {
		t.Errorf("Expected ErrInvalidCounter, got %v", err)
	}
}

func TestCounterService_SetCounter(t *testing.T) {
	svc := &CounterService{}

	// Test empty name validation
	err := svc.SetCounter(context.Background(), "", 100)
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if !errors.Is(err, ErrInvalidCounter) {
		t.Errorf("Expected ErrInvalidCounter, got %v", err)
	}
}

func TestCounterService_ResetCounter(t *testing.T) {
	svc := &CounterService{}

	// Test empty name validation
	err := svc.ResetCounter(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty name")
	}
	if !errors.Is(err, ErrInvalidCounter) {
		t.Errorf("Expected ErrInvalidCounter, got %v", err)
	}
}

func TestStatistics(t *testing.T) {
	stats := &Statistics{
		CustomerAddCount:  10,
		CustomerInqCount:  20,
		CustomerUpdCount:  5,
		PolicyAddCount:    15,
		PolicyInqCount:    25,
		PolicyUpdCount:    8,
		PolicyDelCount:    2,
		ClaimAddCount:     3,
		TotalTransactions: 88,
	}

	if stats.CustomerAddCount != 10 {
		t.Errorf("Expected CustomerAddCount 10, got %d", stats.CustomerAddCount)
	}
	if stats.TotalTransactions != 88 {
		t.Errorf("Expected TotalTransactions 88, got %d", stats.TotalTransactions)
	}
}

// Integration test for CounterService (requires database)
func TestCounterService_Integration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer db.Close()

	counterRepo := repository.NewCounterRepository(db)
	svc := NewCounterService(counterRepo)
	ctx := context.Background()

	// Initialize counters first
	err := svc.InitializeCounters(ctx)
	if err != nil {
		t.Fatalf("InitializeCounters failed: %v", err)
	}

	// Test SetCounter and GetCounter
	t.Run("SetAndGet", func(t *testing.T) {
		err := svc.SetCounter(ctx, "test_counter", 100)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		counter, err := svc.GetCounter(ctx, "test_counter")
		if err != nil {
			t.Fatalf("GetCounter failed: %v", err)
		}

		if counter.Value != 100 {
			t.Errorf("Expected value 100, got %d", counter.Value)
		}
	})

	// Test IncrementCounter
	t.Run("Increment", func(t *testing.T) {
		err := svc.SetCounter(ctx, "inc_counter", 50)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		newValue, err := svc.IncrementCounter(ctx, "inc_counter", 10)
		if err != nil {
			t.Fatalf("IncrementCounter failed: %v", err)
		}

		if newValue != 60 {
			t.Errorf("Expected value 60, got %d", newValue)
		}
	})

	// Test ResetCounter
	t.Run("Reset", func(t *testing.T) {
		err := svc.SetCounter(ctx, "reset_counter", 100)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		err = svc.ResetCounter(ctx, "reset_counter")
		if err != nil {
			t.Fatalf("ResetCounter failed: %v", err)
		}

		value, err := svc.GetCounterValue(ctx, "reset_counter")
		if err != nil {
			t.Fatalf("GetCounterValue failed: %v", err)
		}

		if value != 0 {
			t.Errorf("Expected value 0 after reset, got %d", value)
		}
	})

	// Test GetStatistics
	t.Run("GetStatistics", func(t *testing.T) {
		// Set some counter values
		_ = svc.SetCounter(ctx, models.CounterCustomerAdd, 10)
		_ = svc.SetCounter(ctx, models.CounterPolicyAdd, 20)
		_ = svc.SetCounter(ctx, models.CounterTotalTransactions, 100)

		stats, err := svc.GetStatistics(ctx)
		if err != nil {
			t.Fatalf("GetStatistics failed: %v", err)
		}

		if stats.CustomerAddCount != 10 {
			t.Errorf("Expected CustomerAddCount 10, got %d", stats.CustomerAddCount)
		}
		if stats.PolicyAddCount != 20 {
			t.Errorf("Expected PolicyAddCount 20, got %d", stats.PolicyAddCount)
		}
		if stats.TotalTransactions != 100 {
			t.Errorf("Expected TotalTransactions 100, got %d", stats.TotalTransactions)
		}
	})

	// Test ListCounters
	t.Run("ListCounters", func(t *testing.T) {
		counters, err := svc.ListCounters(ctx)
		if err != nil {
			t.Fatalf("ListCounters failed: %v", err)
		}

		if len(counters) == 0 {
			t.Error("Expected at least one counter")
		}
	})

	// Test NextCustomerNumber
	t.Run("NextCustomerNumber", func(t *testing.T) {
		// Set the counter to a known value
		err := svc.SetCounter(ctx, "customer_num", 1000000000)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		num, err := svc.NextCustomerNumber(ctx)
		if err != nil {
			t.Fatalf("NextCustomerNumber failed: %v", err)
		}

		if len(num) != 10 {
			t.Errorf("Expected 10-digit number, got %s (length %d)", num, len(num))
		}
		if num != "1000000001" {
			t.Errorf("Expected '1000000001', got '%s'", num)
		}
	})

	// Test NextPolicyNumber
	t.Run("NextPolicyNumber", func(t *testing.T) {
		err := svc.SetCounter(ctx, "policy_num", 1000000000)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		num, err := svc.NextPolicyNumber(ctx)
		if err != nil {
			t.Fatalf("NextPolicyNumber failed: %v", err)
		}

		if len(num) != 10 {
			t.Errorf("Expected 10-digit number, got %s (length %d)", num, len(num))
		}
	})

	// Test NextClaimNumber
	t.Run("NextClaimNumber", func(t *testing.T) {
		err := svc.SetCounter(ctx, "claim_num", 1000000000)
		if err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		num, err := svc.NextClaimNumber(ctx)
		if err != nil {
			t.Fatalf("NextClaimNumber failed: %v", err)
		}

		if len(num) != 10 {
			t.Errorf("Expected 10-digit number, got %s (length %d)", num, len(num))
		}
	})
}
