package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cicsdev/genapp/internal/models"
)

// createTestPolicy creates a motor policy for testing claim operations.
func createTestPolicy(t *testing.T, db *DB) *models.Policy {
	t.Helper()

	customer := createTestCustomer(t, db)
	repo := NewPolicyRepository(db)
	ctx := context.Background()

	input := &models.PolicyInput{
		CustomerNum: customer.CustomerNum,
		PolicyType:  models.PolicyTypeMotor,
	}

	policy, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test policy: %v", err)
	}

	return policy
}

func TestClaimRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	input := &models.ClaimInput{
		PolicyNum:    policy.PolicyNum,
		ClaimDate:    timePtr(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)),
		Paid:         float64Ptr(1000.00),
		Value:        float64Ptr(5000.00),
		Cause:        stringPtr("Accident"),
		Observations: stringPtr("Minor damage to front bumper"),
	}

	claim, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if claim.ID == 0 {
		t.Error("Expected ID to be set")
	}
	if claim.ClaimNum == "" {
		t.Error("Expected ClaimNum to be generated")
	}
	if claim.PolicyNum != policy.PolicyNum {
		t.Errorf("Expected PolicyNum '%s', got '%s'", policy.PolicyNum, claim.PolicyNum)
	}
	if claim.GetPaid() != 1000.00 {
		t.Errorf("Expected Paid 1000.00, got %f", claim.GetPaid())
	}
	if claim.GetValue() != 5000.00 {
		t.Errorf("Expected Value 5000.00, got %f", claim.GetValue())
	}
}

func TestClaimRepository_FindByNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create a claim
	input := &models.ClaimInput{
		PolicyNum: policy.PolicyNum,
		Cause:     stringPtr("Theft"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Find by number
	found, err := repo.FindByNum(ctx, created.ClaimNum)
	if err != nil {
		t.Fatalf("FindByNum failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
	}
	if found.GetCause() != "Theft" {
		t.Errorf("Expected Cause 'Theft', got '%s'", found.GetCause())
	}
}

func TestClaimRepository_FindByNum_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewClaimRepository(db)
	ctx := context.Background()

	_, err := repo.FindByNum(ctx, "9999999999")
	if err != ErrClaimNotFound {
		t.Errorf("Expected ErrClaimNotFound, got %v", err)
	}
}

func TestClaimRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create a claim
	input := &models.ClaimInput{
		PolicyNum: policy.PolicyNum,
		Paid:      float64Ptr(0),
		Value:     float64Ptr(2000.00),
		Cause:     stringPtr("Weather damage"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update
	updateInput := &models.ClaimInput{
		Paid:         float64Ptr(1500.00),
		Observations: stringPtr("Claim settled"),
	}
	updated, err := repo.Update(ctx, created.ClaimNum, updateInput)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.GetPaid() != 1500.00 {
		t.Errorf("Expected Paid 1500.00, got %f", updated.GetPaid())
	}
	if updated.GetObservations() != "Claim settled" {
		t.Errorf("Expected Observations 'Claim settled', got '%s'", updated.GetObservations())
	}
	// Cause should be unchanged
	if updated.GetCause() != "Weather damage" {
		t.Errorf("Expected Cause 'Weather damage', got '%s'", updated.GetCause())
	}
}

func TestClaimRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create a claim
	input := &models.ClaimInput{
		PolicyNum: policy.PolicyNum,
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete
	err = repo.Delete(ctx, created.ClaimNum)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByNum(ctx, created.ClaimNum)
	if err != ErrClaimNotFound {
		t.Errorf("Expected ErrClaimNotFound after delete, got %v", err)
	}
}

func TestClaimRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewClaimRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "9999999999")
	if err != ErrClaimNotFound {
		t.Errorf("Expected ErrClaimNotFound, got %v", err)
	}
}

func TestClaimRepository_FindByPolicyNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create multiple claims for the policy
	for i := 0; i < 3; i++ {
		input := &models.ClaimInput{
			PolicyNum: policy.PolicyNum,
			Cause:     stringPtr("Claim " + string(rune('A'+i))),
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Find by policy
	claims, err := repo.FindByPolicyNum(ctx, policy.PolicyNum)
	if err != nil {
		t.Fatalf("FindByPolicyNum failed: %v", err)
	}

	if len(claims) != 3 {
		t.Errorf("Expected 3 claims, got %d", len(claims))
	}
}

func TestClaimRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Initial count
	initialCount, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	// Create claims
	for i := 0; i < 3; i++ {
		input := &models.ClaimInput{
			PolicyNum: policy.PolicyNum,
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

func TestClaimRepository_TotalPaidByPolicy(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create claims with different paid amounts
	amounts := []float64{1000.00, 2500.00, 1500.00}
	for _, amount := range amounts {
		input := &models.ClaimInput{
			PolicyNum: policy.PolicyNum,
			Paid:      float64Ptr(amount),
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Check total
	total, err := repo.TotalPaidByPolicy(ctx, policy.PolicyNum)
	if err != nil {
		t.Fatalf("TotalPaidByPolicy failed: %v", err)
	}

	expected := 5000.00
	if total != expected {
		t.Errorf("Expected total %f, got %f", expected, total)
	}
}

func TestClaimRepository_TotalValueByPolicy(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create claims with different values
	values := []float64{5000.00, 10000.00, 3000.00}
	for _, value := range values {
		input := &models.ClaimInput{
			PolicyNum: policy.PolicyNum,
			Value:     float64Ptr(value),
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Check total
	total, err := repo.TotalValueByPolicy(ctx, policy.PolicyNum)
	if err != nil {
		t.Fatalf("TotalValueByPolicy failed: %v", err)
	}

	expected := 18000.00
	if total != expected {
		t.Errorf("Expected total %f, got %f", expected, total)
	}
}

func TestClaim_Outstanding(t *testing.T) {
	claim := &models.Claim{}
	claim.Value.Valid = true
	claim.Value.Float64 = 5000.00
	claim.Paid.Valid = true
	claim.Paid.Float64 = 2000.00

	outstanding := claim.Outstanding()
	expected := 3000.00

	if outstanding != expected {
		t.Errorf("Expected outstanding %f, got %f", expected, outstanding)
	}
}

func TestClaimRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	policy := createTestPolicy(t, db)
	repo := NewClaimRepository(db)
	ctx := context.Background()

	// Create a claim
	input := &models.ClaimInput{
		PolicyNum: policy.PolicyNum,
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check exists
	exists, err := repo.Exists(ctx, created.ClaimNum)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected claim to exist")
	}

	// Check non-existent
	exists, err = repo.Exists(ctx, "9999999999")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected claim to not exist")
	}
}
