package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cicsdev/genapp/internal/models"
)

// createTestCustomer creates a customer for testing policy operations.
func createTestCustomer(t *testing.T, db *DB) *models.Customer {
	t.Helper()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	input := &models.CustomerInput{
		FirstName: stringPtr("Test"),
		LastName:  stringPtr("Customer"),
	}

	customer, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test customer: %v", err)
	}

	return customer
}

func TestPolicyRepository_CreateMotor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeMotor,
		IssueDate:   timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		ExpiryDate:  timePtr(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		Payment:     float64Ptr(500.00),
		Motor: &models.MotorInput{
			Make:      stringPtr("Toyota"),
			Model:     stringPtr("Camry"),
			Value:     float64Ptr(25000.00),
			RegNumber: stringPtr("AB12CDE"),
			Colour:    stringPtr("Silver"),
			CC:        intPtr(2000),
			Premium:   float64Ptr(500.00),
			Accidents: intPtr(0),
		},
	}

	policy, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if policy.ID == 0 {
		t.Error("Expected ID to be set")
	}
	if policy.PolicyNum == "" {
		t.Error("Expected PolicyNum to be generated")
	}
	if policy.PolicyType != models.PolicyTypeMotor {
		t.Errorf("Expected PolicyType 'M', got '%s'", policy.PolicyType)
	}
	if policy.Motor == nil {
		t.Fatal("Expected Motor details to be set")
	}
	if policy.Motor.GetMake() != "Toyota" {
		t.Errorf("Expected Make 'Toyota', got '%s'", policy.Motor.GetMake())
	}
}

func TestPolicyRepository_CreateEndowment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeEndowment,
		Endowment: &models.EndowmentInput{
			WithProfits: boolPtr(true),
			Equities:    boolPtr(false),
			ManagedFund: boolPtr(true),
			FundName:    stringPtr("GrowthFund"),
			Term:        intPtr(25),
			SumAssured:  float64Ptr(100000.00),
			LifeAssured: stringPtr("John Doe"),
		},
	}

	policy, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if policy.Endowment == nil {
		t.Fatal("Expected Endowment details to be set")
	}
	if !policy.Endowment.GetWithProfits() {
		t.Error("Expected WithProfits to be true")
	}
	if policy.Endowment.GetTerm() != 25 {
		t.Errorf("Expected Term 25, got %d", policy.Endowment.GetTerm())
	}
}

func TestPolicyRepository_CreateHouse(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeHouse,
		House: &models.HouseInput{
			PropertyType: stringPtr("Detached"),
			Bedrooms:     intPtr(4),
			Value:        float64Ptr(350000.00),
			HouseName:    stringPtr("Rose Cottage"),
			HouseNumber:  stringPtr("15"),
			Postcode:     stringPtr("SW1A 1AA"),
		},
	}

	policy, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if policy.House == nil {
		t.Fatal("Expected House details to be set")
	}
	if policy.House.GetPropertyType() != "Detached" {
		t.Errorf("Expected PropertyType 'Detached', got '%s'", policy.House.GetPropertyType())
	}
	if policy.House.GetBedrooms() != 4 {
		t.Errorf("Expected Bedrooms 4, got %d", policy.House.GetBedrooms())
	}
}

func TestPolicyRepository_CreateCommercial(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeCommercial,
		Commercial: &models.CommercialInput{
			Address:      stringPtr("123 Business Street"),
			Postcode:     stringPtr("EC1A 1BB"),
			PropertyType: stringPtr("Office"),
			FirePeril:    intPtr(1),
			FirePremium:  float64Ptr(1000.00),
			CrimePeril:   intPtr(2),
			CrimePremium: float64Ptr(500.00),
			Status:       intPtr(1),
		},
	}

	policy, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if policy.Commercial == nil {
		t.Fatal("Expected Commercial details to be set")
	}
	if policy.Commercial.GetAddress() != "123 Business Street" {
		t.Errorf("Expected Address '123 Business Street', got '%s'", policy.Commercial.GetAddress())
	}
	if policy.Commercial.TotalPremium() != 1500.00 {
		t.Errorf("Expected TotalPremium 1500.00, got %f", policy.Commercial.TotalPremium())
	}
}

func TestPolicyRepository_CreateInvalidType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  PolicyType("X"), // Invalid type
	}

	_, err := repo.Create(ctx, input)
	if err != ErrInvalidPolicyType {
		t.Errorf("Expected ErrInvalidPolicyType, got %v", err)
	}
}

func TestPolicyRepository_FindByNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	// Create a motor policy
	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeMotor,
		Motor: &models.MotorInput{
			Make:  stringPtr("Honda"),
			Model: stringPtr("Civic"),
		},
	}

	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Find by number
	found, err := repo.FindByNum(ctx, created.PolicyNum)
	if err != nil {
		t.Fatalf("FindByNum failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
	}
	if found.Motor == nil {
		t.Fatal("Expected Motor details to be loaded")
	}
	if found.Motor.GetMake() != "Honda" {
		t.Errorf("Expected Make 'Honda', got '%s'", found.Motor.GetMake())
	}
}

func TestPolicyRepository_FindByNum_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPolicyRepository(db)
	ctx := context.Background()

	_, err := repo.FindByNum(ctx, "9999999999")
	if err != ErrPolicyNotFound {
		t.Errorf("Expected ErrPolicyNotFound, got %v", err)
	}
}

func TestPolicyRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	// Create a motor policy
	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeMotor,
		Payment:     float64Ptr(500.00),
		Motor: &models.MotorInput{
			Make:      stringPtr("Ford"),
			Model:     stringPtr("Focus"),
			Premium:   float64Ptr(400.00),
			Accidents: intPtr(0),
		},
	}

	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update
	updateInput := &models.PolicyInput{
		Payment: float64Ptr(550.00),
		Motor: &models.MotorInput{
			Premium:   float64Ptr(450.00),
			Accidents: intPtr(1),
		},
	}

	updated, err := repo.Update(ctx, created.PolicyNum, updateInput)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.GetPayment() != 550.00 {
		t.Errorf("Expected Payment 550.00, got %f", updated.GetPayment())
	}
	if updated.Motor.GetPremium() != 450.00 {
		t.Errorf("Expected Premium 450.00, got %f", updated.Motor.GetPremium())
	}
	if updated.Motor.GetAccidents() != 1 {
		t.Errorf("Expected Accidents 1, got %d", updated.Motor.GetAccidents())
	}
	// Make should be unchanged
	if updated.Motor.GetMake() != "Ford" {
		t.Errorf("Expected Make 'Ford', got '%s'", updated.Motor.GetMake())
	}
}

func TestPolicyRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	// Create a policy
	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeMotor,
	}

	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete
	err = repo.Delete(ctx, created.PolicyNum)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByNum(ctx, created.PolicyNum)
	if err != ErrPolicyNotFound {
		t.Errorf("Expected ErrPolicyNotFound after delete, got %v", err)
	}
}

func TestPolicyRepository_FindByCustomerNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	// Create multiple policies for the customer
	for _, policyType := range []models.PolicyType{models.PolicyTypeMotor, models.PolicyTypeHouse, models.PolicyTypeEndowment} {
		input := &models.PolicyInput{
			CustomerNum: customer.CustomerNum,
			PolicyType:  policyType,
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Find by customer
	policies, err := repo.FindByCustomerNum(ctx, customer.CustomerNum)
	if err != nil {
		t.Fatalf("FindByCustomerNum failed: %v", err)
	}

	if len(policies) != 3 {
		t.Errorf("Expected 3 policies, got %d", len(policies))
	}
}

func TestPolicyRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	// Initial count
	initialCount, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	// Create policies
	for i := 0; i < 3; i++ {
		input := &models.PolicyInput{
			CustomerNum: customer.CustomerNum,
			PolicyType:  models.PolicyTypeMotor,
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Verify count
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != initialCount+3 {
		t.Errorf("Expected count %d, got %d", initialCount+3, count)
	}
}

func TestPolicyType_IsValid(t *testing.T) {
	tests := []struct {
		policyType models.PolicyType
		expected   bool
	}{
		{models.PolicyTypeMotor, true},
		{models.PolicyTypeEndowment, true},
		{models.PolicyTypeHouse, true},
		{models.PolicyTypeCommercial, true},
		{models.PolicyType("X"), false},
		{models.PolicyType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.policyType), func(t *testing.T) {
			result := tt.policyType.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %v for type '%s', got %v", tt.expected, tt.policyType, result)
			}
		})
	}
}

func TestPolicyType_Description(t *testing.T) {
	tests := []struct {
		policyType  models.PolicyType
		description string
	}{
		{models.PolicyTypeMotor, "Motor"},
		{models.PolicyTypeEndowment, "Endowment"},
		{models.PolicyTypeHouse, "House"},
		{models.PolicyTypeCommercial, "Commercial"},
		{models.PolicyType("X"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.policyType), func(t *testing.T) {
			result := tt.policyType.Description()
			if result != tt.description {
				t.Errorf("Expected '%s', got '%s'", tt.description, result)
			}
		})
	}
}

// PolicyType alias for use in tests
type PolicyType = models.PolicyType
