package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

// Counter service errors.
var (
	ErrCounterNotFound = errors.New("counter not found")
	ErrInvalidCounter  = errors.New("invalid counter name")
)

// CounterService provides business logic for counter operations.
// Equivalent to LGSETUP COBOL program which manages Named Counter Server operations.
type CounterService struct {
	counterRepo *repository.CounterRepository
}

// NewCounterService creates a new CounterService.
func NewCounterService(counterRepo *repository.CounterRepository) *CounterService {
	return &CounterService{
		counterRepo: counterRepo,
	}
}

// NextCustomerNumber generates the next customer number.
// This is equivalent to the LGSETUP COBOL program's function of generating
// unique customer numbers using the CICS Named Counter Server.
func (s *CounterService) NextCustomerNumber(ctx context.Context) (string, error) {
	// Increment the customer number counter
	// The counter starts at 1000000001 as per the original COBOL application
	value, err := s.counterRepo.Increment(ctx, "customer_num", 1)
	if err != nil {
		return "", fmt.Errorf("failed to generate customer number: %w", err)
	}

	// Format as 10-digit string with leading zeros
	return fmt.Sprintf("%010d", value), nil
}

// NextPolicyNumber generates the next policy number.
// This is equivalent to the LGSETUP COBOL program's function of generating
// unique policy numbers using the CICS Named Counter Server.
func (s *CounterService) NextPolicyNumber(ctx context.Context) (string, error) {
	// Increment the policy number counter
	value, err := s.counterRepo.Increment(ctx, "policy_num", 1)
	if err != nil {
		return "", fmt.Errorf("failed to generate policy number: %w", err)
	}

	// Format as 10-digit string with leading zeros
	return fmt.Sprintf("%010d", value), nil
}

// NextClaimNumber generates the next claim number.
func (s *CounterService) NextClaimNumber(ctx context.Context) (string, error) {
	value, err := s.counterRepo.Increment(ctx, "claim_num", 1)
	if err != nil {
		return "", fmt.Errorf("failed to generate claim number: %w", err)
	}

	return fmt.Sprintf("%010d", value), nil
}

// GetCounter retrieves a counter by name.
func (s *CounterService) GetCounter(ctx context.Context, name string) (*models.Counter, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidCounter)
	}

	counter, err := s.counterRepo.Get(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrCounterNotFound) {
			return nil, ErrCounterNotFound
		}
		return nil, fmt.Errorf("failed to get counter: %w", err)
	}

	return counter, nil
}

// GetCounterValue retrieves just the value of a counter.
func (s *CounterService) GetCounterValue(ctx context.Context, name string) (int64, error) {
	counter, err := s.GetCounter(ctx, name)
	if err != nil {
		return 0, err
	}
	return counter.Value, nil
}

// IncrementCounter increments a counter by the specified delta and returns the new value.
func (s *CounterService) IncrementCounter(ctx context.Context, name string, delta int64) (int64, error) {
	if name == "" {
		return 0, fmt.Errorf("%w: name is required", ErrInvalidCounter)
	}

	value, err := s.counterRepo.Increment(ctx, name, delta)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return value, nil
}

// SetCounter sets a counter to a specific value.
func (s *CounterService) SetCounter(ctx context.Context, name string, value int64) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidCounter)
	}

	return s.counterRepo.Set(ctx, name, value)
}

// ResetCounter resets a counter to zero.
func (s *CounterService) ResetCounter(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidCounter)
	}

	return s.counterRepo.Reset(ctx, name)
}

// ListCounters retrieves all counters.
func (s *CounterService) ListCounters(ctx context.Context) ([]*models.Counter, error) {
	return s.counterRepo.List(ctx)
}

// Statistics represents aggregate statistics for the application.
type Statistics struct {
	CustomerAddCount    int64 `json:"customer_add_count"`
	CustomerInqCount    int64 `json:"customer_inq_count"`
	CustomerUpdCount    int64 `json:"customer_upd_count"`
	PolicyAddCount      int64 `json:"policy_add_count"`
	PolicyInqCount      int64 `json:"policy_inq_count"`
	PolicyUpdCount      int64 `json:"policy_upd_count"`
	PolicyDelCount      int64 `json:"policy_del_count"`
	ClaimAddCount       int64 `json:"claim_add_count"`
	TotalTransactions   int64 `json:"total_transactions"`
}

// GetStatistics retrieves all application statistics.
func (s *CounterService) GetStatistics(ctx context.Context) (*Statistics, error) {
	counters, err := s.counterRepo.GetMultiple(ctx, []string{
		models.CounterCustomerAdd,
		models.CounterCustomerInquiry,
		models.CounterCustomerUpdate,
		models.CounterPolicyAdd,
		models.CounterPolicyInquiry,
		models.CounterPolicyUpdate,
		models.CounterPolicyDelete,
		models.CounterClaimAdd,
		models.CounterTotalTransactions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	stats := &Statistics{}
	for _, counter := range counters {
		switch counter.Name {
		case models.CounterCustomerAdd:
			stats.CustomerAddCount = counter.Value
		case models.CounterCustomerInquiry:
			stats.CustomerInqCount = counter.Value
		case models.CounterCustomerUpdate:
			stats.CustomerUpdCount = counter.Value
		case models.CounterPolicyAdd:
			stats.PolicyAddCount = counter.Value
		case models.CounterPolicyInquiry:
			stats.PolicyInqCount = counter.Value
		case models.CounterPolicyUpdate:
			stats.PolicyUpdCount = counter.Value
		case models.CounterPolicyDelete:
			stats.PolicyDelCount = counter.Value
		case models.CounterClaimAdd:
			stats.ClaimAddCount = counter.Value
		case models.CounterTotalTransactions:
			stats.TotalTransactions = counter.Value
		}
	}

	return stats, nil
}

// InitializeCounters initializes all application counters with default values.
// This should be called during application startup if counters don't exist.
func (s *CounterService) InitializeCounters(ctx context.Context) error {
	// Initialize transaction counters to 0
	countersToInit := []string{
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

	for _, name := range countersToInit {
		// Only set if counter doesn't exist (uses ON CONFLICT DO UPDATE)
		// This ensures we don't overwrite existing values
		exists, err := s.counterRepo.Exists(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to check counter %s: %w", name, err)
		}
		if !exists {
			if err := s.counterRepo.Set(ctx, name, 0); err != nil {
				return fmt.Errorf("failed to initialize counter %s: %w", name, err)
			}
		}
	}

	return nil
}
