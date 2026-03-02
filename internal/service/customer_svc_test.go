package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

// MockCustomerRepository is a mock implementation of customer repository for testing.
type MockCustomerRepository struct {
	CreateFunc      func(ctx context.Context, input *models.CustomerInput) (*models.Customer, error)
	FindByNumFunc   func(ctx context.Context, customerNum string) (*models.Customer, error)
	UpdateFunc      func(ctx context.Context, customerNum string, input *models.CustomerInput) (*models.Customer, error)
	DeleteFunc      func(ctx context.Context, customerNum string) error
	ExistsFunc      func(ctx context.Context, customerNum string) (bool, error)
	ListFunc        func(ctx context.Context, offset, limit int) ([]*models.Customer, error)
	SearchFunc      func(ctx context.Context, lastName string, limit int) ([]*models.Customer, error)
	CountFunc       func(ctx context.Context) (int64, error)
}

func (m *MockCustomerRepository) Create(ctx context.Context, input *models.CustomerInput) (*models.Customer, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCustomerRepository) FindByNum(ctx context.Context, customerNum string) (*models.Customer, error) {
	if m.FindByNumFunc != nil {
		return m.FindByNumFunc(ctx, customerNum)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCustomerRepository) Update(ctx context.Context, customerNum string, input *models.CustomerInput) (*models.Customer, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, customerNum, input)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCustomerRepository) Delete(ctx context.Context, customerNum string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, customerNum)
	}
	return errors.New("not implemented")
}

func (m *MockCustomerRepository) Exists(ctx context.Context, customerNum string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, customerNum)
	}
	return false, errors.New("not implemented")
}

func (m *MockCustomerRepository) List(ctx context.Context, offset, limit int) ([]*models.Customer, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, offset, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCustomerRepository) SearchByLastName(ctx context.Context, lastName string, limit int) ([]*models.Customer, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, lastName, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCustomerRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, errors.New("not implemented")
}

// MockCounterRepository is a mock implementation of counter repository for testing.
type MockCounterRepository struct {
	IncrementFunc func(ctx context.Context, name string, delta int64) (int64, error)
	GetFunc       func(ctx context.Context, name string) (*models.Counter, error)
	SetFunc       func(ctx context.Context, name string, value int64) error
	ExistsFunc    func(ctx context.Context, name string) (bool, error)
	ListFunc      func(ctx context.Context) ([]*models.Counter, error)
	ResetFunc     func(ctx context.Context, name string) error
	GetMultipleFunc func(ctx context.Context, names []string) ([]*models.Counter, error)
}

func (m *MockCounterRepository) Increment(ctx context.Context, name string, delta int64) (int64, error) {
	if m.IncrementFunc != nil {
		return m.IncrementFunc(ctx, name, delta)
	}
	return 0, nil // Default: silent success for counter increments
}

func (m *MockCounterRepository) Get(ctx context.Context, name string) (*models.Counter, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCounterRepository) Set(ctx context.Context, name string, value int64) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, name, value)
	}
	return nil
}

func (m *MockCounterRepository) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}

func (m *MockCounterRepository) List(ctx context.Context) ([]*models.Counter, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCounterRepository) Reset(ctx context.Context, name string) error {
	if m.ResetFunc != nil {
		return m.ResetFunc(ctx, name)
	}
	return nil
}

func (m *MockCounterRepository) GetMultiple(ctx context.Context, names []string) ([]*models.Counter, error) {
	if m.GetMultipleFunc != nil {
		return m.GetMultipleFunc(ctx, names)
	}
	return nil, errors.New("not implemented")
}

// createTestCustomerService creates a service with mock repositories for testing.
func createTestCustomerService(custRepo *MockCustomerRepository, counterRepo *MockCounterRepository) *CustomerService {
	// The service uses concrete repository types, so we need to adapt
	// For now, we'll test via integration tests or use interfaces in production code
	// This demonstrates the testing approach
	return &CustomerService{
		customerRepo: nil, // Would use interface in production
		counterRepo:  nil,
	}
}

func TestCustomerService_ValidateCustomerNum(t *testing.T) {
	svc := &CustomerService{}

	tests := []struct {
		name        string
		customerNum string
		wantErr     bool
		errContains string
	}{
		{"valid number", "1234567890", false, ""},
		{"empty", "", true, "required"},
		{"too short", "123456789", true, "10 digits"},
		{"too long", "12345678901", true, "10 digits"},
		{"non-numeric", "123456789a", true, "numeric"},
		{"letters only", "abcdefghij", true, "numeric"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validateCustomerNum(tt.customerNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCustomerNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("validateCustomerNum() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestCustomerService_ValidateAddInput(t *testing.T) {
	svc := &CustomerService{}

	tests := []struct {
		name        string
		input       *AddCustomerInput
		wantErr     bool
		errContains string
	}{
		{
			name: "valid input",
			input: &AddCustomerInput{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing last name",
			input: &AddCustomerInput{
				FirstName: "John",
			},
			wantErr:     true,
			errContains: "last name is required",
		},
		{
			name: "first name too long",
			input: &AddCustomerInput{
				FirstName: "JohnJohnJohn", // > 10 chars
				LastName:  "Doe",
			},
			wantErr:     true,
			errContains: "first name cannot exceed 10 characters",
		},
		{
			name: "last name too long",
			input: &AddCustomerInput{
				FirstName: "John",
				LastName:  "DoeDoeDoeDoeDoeDoeDoeDoe", // > 20 chars
			},
			wantErr:     true,
			errContains: "last name cannot exceed 20 characters",
		},
		{
			name: "invalid email",
			input: &AddCustomerInput{
				FirstName:    "John",
				LastName:     "Doe",
				EmailAddress: "invalid-email",
			},
			wantErr:     true,
			errContains: "invalid email address format",
		},
		{
			name: "postcode too long",
			input: &AddCustomerInput{
				FirstName: "John",
				LastName:  "Doe",
				Postcode:  "AB12 3CDE", // > 8 chars
			},
			wantErr:     true,
			errContains: "postcode cannot exceed 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validateAddInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAddInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("validateAddInput() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestCustomerService_ValidateUpdateInput(t *testing.T) {
	svc := &CustomerService{}

	tests := []struct {
		name        string
		input       *UpdateCustomerInput
		wantErr     bool
		errContains string
	}{
		{
			name:    "nil values are ok",
			input:   &UpdateCustomerInput{},
			wantErr: false,
		},
		{
			name: "valid partial update",
			input: &UpdateCustomerInput{
				FirstName: strPtr("Jane"),
			},
			wantErr: false,
		},
		{
			name: "first name too long",
			input: &UpdateCustomerInput{
				FirstName: strPtr("JohnJohnJohn"),
			},
			wantErr:     true,
			errContains: "first name cannot exceed 10 characters",
		},
		{
			name: "invalid email",
			input: &UpdateCustomerInput{
				EmailAddress: strPtr("not-an-email"),
			},
			wantErr:     true,
			errContains: "invalid email address format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.validateUpdateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUpdateInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("validateUpdateInput() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.org", true},
		{"invalid", false},
		{"@nodomain.com", false},
		{"noat.com", false},
		{"", false},
		{"user@", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}

func TestStrPtr(t *testing.T) {
	// Test empty string returns nil
	result := strPtr("")
	if result != nil {
		t.Error("strPtr(\"\") should return nil")
	}

	// Test non-empty string returns pointer
	result = strPtr("hello")
	if result == nil {
		t.Error("strPtr(\"hello\") should not return nil")
	}
	if *result != "hello" {
		t.Errorf("strPtr(\"hello\") = %q, want \"hello\"", *result)
	}
}

// Integration test for CustomerService (requires database)
func TestCustomerService_Integration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db)
	counterRepo := repository.NewCounterRepository(db)
	svc := NewCustomerService(customerRepo, counterRepo)
	ctx := context.Background()

	// Test Add
	t.Run("Add", func(t *testing.T) {
		dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
		input := &AddCustomerInput{
			FirstName:    "John",
			LastName:     "Doe",
			DateOfBirth:  &dob,
			HouseName:    "Oak Villa",
			HouseNumber:  "42",
			Postcode:     "AB12 3CD",
			PhoneHome:    "01onal2345678",
			PhoneMobile:  "07123456789",
			EmailAddress: "john.doe@example.com",
		}

		result, err := svc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		if result.CustomerNum == "" {
			t.Error("Expected CustomerNum to be generated")
		}
		if result.Customer.GetFirstName() != "John" {
			t.Errorf("Expected FirstName 'John', got '%s'", result.Customer.GetFirstName())
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		// First add a customer
		input := &AddCustomerInput{
			FirstName: "Jane",
			LastName:  "Smith",
		}
		addResult, err := svc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Get the customer
		customer, err := svc.Get(ctx, addResult.CustomerNum)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if customer.GetFirstName() != "Jane" {
			t.Errorf("Expected FirstName 'Jane', got '%s'", customer.GetFirstName())
		}
	})

	// Test Get not found
	t.Run("Get_NotFound", func(t *testing.T) {
		_, err := svc.Get(ctx, "9999999999")
		if !errors.Is(err, ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound, got %v", err)
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		// First add a customer
		input := &AddCustomerInput{
			FirstName: "Bob",
			LastName:  "Wilson",
		}
		addResult, err := svc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Update
		updateInput := &UpdateCustomerInput{
			LastName:    strPtr("Johnson"),
			PhoneMobile: strPtr("07999888777"),
		}
		customer, err := svc.Update(ctx, addResult.CustomerNum, updateInput)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if customer.GetLastName() != "Johnson" {
			t.Errorf("Expected LastName 'Johnson', got '%s'", customer.GetLastName())
		}
		if customer.GetPhoneMobile() != "07999888777" {
			t.Errorf("Expected PhoneMobile '07999888777', got '%s'", customer.GetPhoneMobile())
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		// First add a customer
		input := &AddCustomerInput{
			FirstName: "Delete",
			LastName:  "Me",
		}
		addResult, err := svc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Delete
		err = svc.Delete(ctx, addResult.CustomerNum)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify deleted
		_, err = svc.Get(ctx, addResult.CustomerNum)
		if !errors.Is(err, ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound after delete, got %v", err)
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// setupTestDB creates a database connection for integration tests.
func setupTestDB(t *testing.T) *repository.DB {
	t.Helper()

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=genapp_test sslmode=disable"

	db, err := repository.NewDBFromDSN(dsn)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return nil
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		t.Skipf("Skipping integration test: %v", err)
		return nil
	}

	// Clean up test data
	ctx := context.Background()
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
		_, _ = db.ExecContext(ctx, "DELETE FROM "+table)
	}

	// Reset sequences
	sequences := []string{
		"customer_num_seq",
		"policy_num_seq",
		"claim_num_seq",
	}
	for _, seq := range sequences {
		_, _ = db.ExecContext(ctx, "ALTER SEQUENCE "+seq+" RESTART WITH 1000000001")
	}

	return db
}

// createTestCustomer creates a customer model for testing
func createTestCustomer(num, firstName, lastName string) *models.Customer {
	c := &models.Customer{
		ID:          1,
		CustomerNum: num,
	}
	if firstName != "" {
		c.FirstName = sql.NullString{String: firstName, Valid: true}
	}
	if lastName != "" {
		c.LastName = sql.NullString{String: lastName, Valid: true}
	}
	return c
}
