package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

func TestPolicyService_ValidatePolicyNum(t *testing.T) {
	svc := &PolicyService{}

	tests := []struct {
		name       string
		policyNum  string
		wantErr    bool
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
			err := svc.validatePolicyNum(tt.policyNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePolicyNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("validatePolicyNum() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestPolicyService_ValidateAddInput(t *testing.T) {
	svc := &PolicyService{}

	tests := []struct {
		name        string
		input       *AddPolicyInput
		wantErr     bool
		errContains string
	}{
		{
			name: "valid motor policy",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeMotor,
				Motor: &AddMotorInput{
					Make:  "Toyota",
					Model: "Camry",
				},
			},
			wantErr: false,
		},
		{
			name: "valid endowment policy",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeEndowment,
				Endowment: &AddEndowmentInput{
					WithProfits: true,
					Term:        10,
				},
			},
			wantErr: false,
		},
		{
			name: "valid house policy",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeHouse,
				House: &AddHouseInput{
					PropertyType: "Detached",
					Bedrooms:     3,
				},
			},
			wantErr: false,
		},
		{
			name: "valid commercial policy",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeCommercial,
				Commercial:  &AddCommercialInput{},
			},
			wantErr: false,
		},
		{
			name: "missing customer number",
			input: &AddPolicyInput{
				PolicyType: models.PolicyTypeMotor,
				Motor: &AddMotorInput{
					Make: "Toyota",
				},
			},
			wantErr:     true,
			errContains: "customer number is required",
		},
		{
			name: "invalid customer number length",
			input: &AddPolicyInput{
				CustomerNum: "123",
				PolicyType:  models.PolicyTypeMotor,
				Motor: &AddMotorInput{
					Make: "Toyota",
				},
			},
			wantErr:     true,
			errContains: "customer number must be 10 digits",
		},
		{
			name: "invalid policy type",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  "X",
			},
			wantErr:     true,
			errContains: "valid policy type is required",
		},
		{
			name: "motor policy missing details",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeMotor,
			},
			wantErr:     true,
			errContains: "motor policy details are required",
		},
		{
			name: "motor policy missing make",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeMotor,
				Motor:       &AddMotorInput{},
			},
			wantErr:     true,
			errContains: "car make is required",
		},
		{
			name: "motor policy make too long",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeMotor,
				Motor: &AddMotorInput{
					Make: "ToyotaVeryLongName", // > 15 chars
				},
			},
			wantErr:     true,
			errContains: "car make cannot exceed 15 characters",
		},
		{
			name: "endowment policy missing details",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeEndowment,
			},
			wantErr:     true,
			errContains: "endowment policy details are required",
		},
		{
			name: "endowment policy term out of range",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeEndowment,
				Endowment: &AddEndowmentInput{
					Term: 100, // > 99
				},
			},
			wantErr:     true,
			errContains: "term must be between 0 and 99",
		},
		{
			name: "house policy missing details",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeHouse,
			},
			wantErr:     true,
			errContains: "house policy details are required",
		},
		{
			name: "commercial policy missing details",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeCommercial,
			},
			wantErr:     true,
			errContains: "commercial policy details are required",
		},
		{
			name: "broker ID too long",
			input: &AddPolicyInput{
				CustomerNum: "1234567890",
				PolicyType:  models.PolicyTypeMotor,
				BrokerID:    "12345678901", // > 10 chars
				Motor: &AddMotorInput{
					Make: "Toyota",
				},
			},
			wantErr:     true,
			errContains: "broker ID cannot exceed 10 characters",
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

func TestPolicyService_BuildPolicyInput(t *testing.T) {
	svc := &PolicyService{}

	issueDate := time.Now()
	expiryDate := issueDate.AddDate(1, 0, 0)
	payment := 100.50

	t.Run("motor policy", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: "1234567890",
			PolicyType:  models.PolicyTypeMotor,
			IssueDate:   &issueDate,
			ExpiryDate:  &expiryDate,
			BrokerID:    "BR123",
			BrokersRef:  "REF456",
			Payment:     &payment,
			Motor: &AddMotorInput{
				Make:      "Toyota",
				Model:     "Camry",
				Value:     25000,
				RegNumber: "ABC123",
				Colour:    "Blue",
				CC:        2000,
				Premium:   500,
				Accidents: 0,
			},
		}

		result := svc.buildPolicyInput(input)

		if result.CustomerNum != "1234567890" {
			t.Errorf("Expected CustomerNum '1234567890', got '%s'", result.CustomerNum)
		}
		if result.PolicyType != models.PolicyTypeMotor {
			t.Errorf("Expected PolicyType Motor, got '%s'", result.PolicyType)
		}
		if result.Motor == nil {
			t.Fatal("Expected Motor details to be set")
		}
		if result.Motor.Make == nil || *result.Motor.Make != "Toyota" {
			t.Error("Expected Motor.Make to be 'Toyota'")
		}
	})

	t.Run("endowment policy", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: "1234567890",
			PolicyType:  models.PolicyTypeEndowment,
			Endowment: &AddEndowmentInput{
				WithProfits: true,
				Equities:    false,
				ManagedFund: true,
				FundName:    "Growth",
				Term:        10,
				SumAssured:  100000,
				LifeAssured: "John Doe",
			},
		}

		result := svc.buildPolicyInput(input)

		if result.Endowment == nil {
			t.Fatal("Expected Endowment details to be set")
		}
		if result.Endowment.WithProfits == nil || !*result.Endowment.WithProfits {
			t.Error("Expected Endowment.WithProfits to be true")
		}
	})

	t.Run("house policy", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: "1234567890",
			PolicyType:  models.PolicyTypeHouse,
			House: &AddHouseInput{
				PropertyType: "Detached",
				Bedrooms:     3,
				Value:        350000,
				HouseName:    "Rose Cottage",
				HouseNumber:  "1",
				Postcode:     "AB1 2CD",
			},
		}

		result := svc.buildPolicyInput(input)

		if result.House == nil {
			t.Fatal("Expected House details to be set")
		}
		if result.House.PropertyType == nil || *result.House.PropertyType != "Detached" {
			t.Error("Expected House.PropertyType to be 'Detached'")
		}
	})

	t.Run("commercial policy", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: "1234567890",
			PolicyType:  models.PolicyTypeCommercial,
			Commercial: &AddCommercialInput{
				Address:     "123 Business St",
				Postcode:    "CD3 4EF",
				FirePeril:   1,
				FirePremium: 1000,
			},
		}

		result := svc.buildPolicyInput(input)

		if result.Commercial == nil {
			t.Fatal("Expected Commercial details to be set")
		}
		if result.Commercial.Address == nil || *result.Commercial.Address != "123 Business St" {
			t.Error("Expected Commercial.Address to be '123 Business St'")
		}
	})
}

func TestPolicyService_BuildUpdateInput(t *testing.T) {
	svc := &PolicyService{}

	newValue := 30000.0
	newMake := "Honda"

	t.Run("motor policy update", func(t *testing.T) {
		input := &UpdatePolicyInput{
			Motor: &UpdateMotorInput{
				Make:  &newMake,
				Value: &newValue,
			},
		}

		result := svc.buildUpdateInput(models.PolicyTypeMotor, input)

		if result.Motor == nil {
			t.Fatal("Expected Motor details to be set")
		}
		if result.Motor.Make == nil || *result.Motor.Make != "Honda" {
			t.Error("Expected Motor.Make to be 'Honda'")
		}
		if result.Motor.Value == nil || *result.Motor.Value != 30000.0 {
			t.Error("Expected Motor.Value to be 30000.0")
		}
	})
}

func TestFloat64Ptr(t *testing.T) {
	// Test zero returns nil
	result := float64Ptr(0)
	if result != nil {
		t.Error("float64Ptr(0) should return nil")
	}

	// Test non-zero returns pointer
	result = float64Ptr(100.5)
	if result == nil {
		t.Error("float64Ptr(100.5) should not return nil")
	}
	if *result != 100.5 {
		t.Errorf("float64Ptr(100.5) = %f, want 100.5", *result)
	}
}

func TestIntPtr(t *testing.T) {
	// Test zero returns nil
	result := intPtr(0)
	if result != nil {
		t.Error("intPtr(0) should return nil")
	}

	// Test non-zero returns pointer
	result = intPtr(42)
	if result == nil {
		t.Error("intPtr(42) should not return nil")
	}
	if *result != 42 {
		t.Errorf("intPtr(42) = %d, want 42", *result)
	}
}

func TestBoolPtr(t *testing.T) {
	// Test false returns pointer to false (not nil)
	result := boolPtr(false)
	if result == nil {
		t.Error("boolPtr(false) should not return nil")
	}
	if *result != false {
		t.Error("boolPtr(false) should return pointer to false")
	}

	// Test true returns pointer to true
	result = boolPtr(true)
	if result == nil {
		t.Error("boolPtr(true) should not return nil")
	}
	if *result != true {
		t.Error("boolPtr(true) should return pointer to true")
	}
}

// Integration test for PolicyService (requires database)
func TestPolicyService_Integration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Database not available")
	}
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	counterRepo := repository.NewCounterRepository(db)

	customerSvc := NewCustomerService(customerRepo, counterRepo)
	policySvc := NewPolicyService(policyRepo, customerRepo, counterRepo)
	ctx := context.Background()

	// First, create a customer
	customerInput := &AddCustomerInput{
		FirstName: "Policy",
		LastName:  "Test",
	}
	customerResult, err := customerSvc.Add(ctx, customerInput)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}
	customerNum := customerResult.CustomerNum

	// Test Add Motor Policy
	t.Run("Add_Motor", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: customerNum,
			PolicyType:  models.PolicyTypeMotor,
			Motor: &AddMotorInput{
				Make:      "Toyota",
				Model:     "Camry",
				Value:     25000,
				RegNumber: "ABC123",
				Colour:    "Blue",
				CC:        2000,
				Premium:   500,
			},
		}

		result, err := policySvc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		if result.PolicyNum == "" {
			t.Error("Expected PolicyNum to be generated")
		}
		if result.Policy.PolicyType != models.PolicyTypeMotor {
			t.Errorf("Expected PolicyType Motor, got '%s'", result.Policy.PolicyType)
		}
		if result.Policy.Motor == nil {
			t.Error("Expected Motor details to be populated")
		}
	})

	// Test Add with non-existent customer
	t.Run("Add_CustomerNotFound", func(t *testing.T) {
		input := &AddPolicyInput{
			CustomerNum: "9999999999",
			PolicyType:  models.PolicyTypeMotor,
			Motor: &AddMotorInput{
				Make: "Toyota",
			},
		}

		_, err := policySvc.Add(ctx, input)
		if err == nil {
			t.Error("Expected error for non-existent customer")
		}
		if !errors.Is(err, ErrCustomerNotFound) {
			// The error wraps ErrCustomerNotFound with customer number
			if !contains(err.Error(), "customer") {
				t.Errorf("Expected customer not found error, got: %v", err)
			}
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		// First add a policy
		input := &AddPolicyInput{
			CustomerNum: customerNum,
			PolicyType:  models.PolicyTypeEndowment,
			Endowment: &AddEndowmentInput{
				WithProfits: true,
				Term:        10,
				SumAssured:  100000,
			},
		}
		addResult, err := policySvc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Get the policy
		policy, err := policySvc.Get(ctx, addResult.PolicyNum)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if policy.PolicyType != models.PolicyTypeEndowment {
			t.Errorf("Expected PolicyType Endowment, got '%s'", policy.PolicyType)
		}
		if policy.Endowment == nil {
			t.Error("Expected Endowment details to be populated")
		}
	})

	// Test Get not found
	t.Run("Get_NotFound", func(t *testing.T) {
		_, err := policySvc.Get(ctx, "9999999999")
		if !errors.Is(err, ErrPolicyNotFound) {
			t.Errorf("Expected ErrPolicyNotFound, got %v", err)
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		// First add a policy
		input := &AddPolicyInput{
			CustomerNum: customerNum,
			PolicyType:  models.PolicyTypeHouse,
			House: &AddHouseInput{
				PropertyType: "Detached",
				Bedrooms:     3,
				Value:        350000,
			},
		}
		addResult, err := policySvc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Update
		newValue := 400000.0
		newBedrooms := 4
		updateInput := &UpdatePolicyInput{
			House: &UpdateHouseInput{
				Value:    &newValue,
				Bedrooms: &newBedrooms,
			},
		}
		policy, err := policySvc.Update(ctx, addResult.PolicyNum, updateInput)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if policy.House == nil {
			t.Fatal("Expected House details to be populated")
		}
		if policy.House.GetValue() != 400000 {
			t.Errorf("Expected Value 400000, got %f", policy.House.GetValue())
		}
		if policy.House.GetBedrooms() != 4 {
			t.Errorf("Expected Bedrooms 4, got %d", policy.House.GetBedrooms())
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		// First add a policy
		input := &AddPolicyInput{
			CustomerNum: customerNum,
			PolicyType:  models.PolicyTypeCommercial,
			Commercial:  &AddCommercialInput{},
		}
		addResult, err := policySvc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		// Delete
		err = policySvc.Delete(ctx, addResult.PolicyNum)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify deleted
		_, err = policySvc.Get(ctx, addResult.PolicyNum)
		if !errors.Is(err, ErrPolicyNotFound) {
			t.Errorf("Expected ErrPolicyNotFound after delete, got %v", err)
		}
	})

	// Test GetByCustomer
	t.Run("GetByCustomer", func(t *testing.T) {
		// Add another policy for the same customer
		input := &AddPolicyInput{
			CustomerNum: customerNum,
			PolicyType:  models.PolicyTypeMotor,
			Motor: &AddMotorInput{
				Make: "Honda",
			},
		}
		_, err := policySvc.Add(ctx, input)
		if err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		policies, err := policySvc.GetByCustomer(ctx, customerNum)
		if err != nil {
			t.Fatalf("GetByCustomer failed: %v", err)
		}

		if len(policies) < 1 {
			t.Error("Expected at least one policy for customer")
		}
	})
}
