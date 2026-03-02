package repository

import (
	"context"
	"testing"
	"time"

	"github.com/cicsdev/genapp/internal/models"
)

func TestCustomerRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	input := &models.CustomerInput{
		FirstName:    stringPtr("John"),
		LastName:     stringPtr("Doe"),
		DateOfBirth:  timePtr(time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)),
		HouseName:    stringPtr("Oak Villa"),
		HouseNumber:  stringPtr("42"),
		Postcode:     stringPtr("AB12 3CD"),
		PhoneHome:    stringPtr("01onal2345678"),
		PhoneMobile:  stringPtr("07123456789"),
		EmailAddress: stringPtr("john.doe@example.com"),
	}

	customer, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if customer.ID == 0 {
		t.Error("Expected ID to be set")
	}
	if customer.CustomerNum == "" {
		t.Error("Expected CustomerNum to be generated")
	}
	if customer.GetFirstName() != "John" {
		t.Errorf("Expected FirstName 'John', got '%s'", customer.GetFirstName())
	}
	if customer.GetLastName() != "Doe" {
		t.Errorf("Expected LastName 'Doe', got '%s'", customer.GetLastName())
	}
}

func TestCustomerRepository_FindByNum(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Create a customer first
	input := &models.CustomerInput{
		FirstName: stringPtr("Jane"),
		LastName:  stringPtr("Smith"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Find by number
	found, err := repo.FindByNum(ctx, created.CustomerNum)
	if err != nil {
		t.Fatalf("FindByNum failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
	}
	if found.GetFirstName() != "Jane" {
		t.Errorf("Expected FirstName 'Jane', got '%s'", found.GetFirstName())
	}
}

func TestCustomerRepository_FindByNum_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	_, err := repo.FindByNum(ctx, "9999999999")
	if err != ErrCustomerNotFound {
		t.Errorf("Expected ErrCustomerNotFound, got %v", err)
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Create a customer
	input := &models.CustomerInput{
		FirstName: stringPtr("Bob"),
		LastName:  stringPtr("Wilson"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update
	updateInput := &models.CustomerInput{
		LastName:    stringPtr("Johnson"),
		PhoneMobile: stringPtr("07999888777"),
	}
	updated, err := repo.Update(ctx, created.CustomerNum, updateInput)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	if updated.GetFirstName() != "Bob" {
		t.Errorf("FirstName should remain 'Bob', got '%s'", updated.GetFirstName())
	}
	if updated.GetLastName() != "Johnson" {
		t.Errorf("Expected LastName 'Johnson', got '%s'", updated.GetLastName())
	}
	if updated.GetPhoneMobile() != "07999888777" {
		t.Errorf("Expected PhoneMobile '07999888777', got '%s'", updated.GetPhoneMobile())
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Create a customer
	input := &models.CustomerInput{
		FirstName: stringPtr("Delete"),
		LastName:  stringPtr("Me"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete
	err = repo.Delete(ctx, created.CustomerNum)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByNum(ctx, created.CustomerNum)
	if err != ErrCustomerNotFound {
		t.Errorf("Expected ErrCustomerNotFound after delete, got %v", err)
	}
}

func TestCustomerRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "9999999999")
	if err != ErrCustomerNotFound {
		t.Errorf("Expected ErrCustomerNotFound, got %v", err)
	}
}

func TestCustomerRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Create some customers
	for i := 0; i < 5; i++ {
		name := stringPtr("Customer")
		input := &models.CustomerInput{
			FirstName: name,
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// List with pagination
	customers, err := repo.List(ctx, 0, 3)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(customers) != 3 {
		t.Errorf("Expected 3 customers, got %d", len(customers))
	}

	// List second page
	customers, err = repo.List(ctx, 3, 3)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(customers) != 2 {
		t.Errorf("Expected 2 customers on second page, got %d", len(customers))
	}
}

func TestCustomerRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Initial count
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	initialCount := count

	// Create customers
	for i := 0; i < 3; i++ {
		input := &models.CustomerInput{
			FirstName: stringPtr("Test"),
		}
		_, err := repo.Create(ctx, input)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Verify count
	count, err = repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != initialCount+3 {
		t.Errorf("Expected count %d, got %d", initialCount+3, count)
	}
}

func TestCustomerRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepository(db)
	ctx := context.Background()

	// Create a customer
	input := &models.CustomerInput{
		FirstName: stringPtr("Exists"),
	}
	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check exists
	exists, err := repo.Exists(ctx, created.CustomerNum)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Expected customer to exist")
	}

	// Check non-existent
	exists, err = repo.Exists(ctx, "9999999999")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Expected customer to not exist")
	}
}

func TestCustomer_FullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName *string
		lastName  *string
		expected  string
	}{
		{"Both names", stringPtr("John"), stringPtr("Doe"), "John Doe"},
		{"First only", stringPtr("John"), nil, "John"},
		{"Last only", nil, stringPtr("Doe"), "Doe"},
		{"Neither", nil, nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := &models.Customer{}
			if tt.firstName != nil {
				customer.FirstName.Valid = true
				customer.FirstName.String = *tt.firstName
			}
			if tt.lastName != nil {
				customer.LastName.Valid = true
				customer.LastName.String = *tt.lastName
			}

			result := customer.FullName()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
