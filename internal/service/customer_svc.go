// Package service implements the business logic layer for the GENAPP application.
// Services correspond to the COBOL "US" (user service) programs that orchestrate
// business operations by calling data access programs.
package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

// Common service errors.
var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrInvalidCustomerNumber = errors.New("invalid customer number")
	ErrInvalidInput          = errors.New("invalid input")
	ErrValidationFailed      = errors.New("validation failed")
)

// CustomerService provides business logic for customer operations.
// Equivalent to COBOL programs: LGACUS01 (add), LGICUS01 (inquire), LGUCUS01 (update).
type CustomerService struct {
	customerRepo *repository.CustomerRepository
	counterRepo  *repository.CounterRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(
	customerRepo *repository.CustomerRepository,
	counterRepo *repository.CounterRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		counterRepo:  counterRepo,
	}
}

// AddCustomerInput represents the input for adding a new customer.
type AddCustomerInput struct {
	FirstName    string
	LastName     string
	DateOfBirth  *time.Time
	HouseName    string
	HouseNumber  string
	Postcode     string
	PhoneHome    string
	PhoneMobile  string
	EmailAddress string
}

// AddCustomerResult represents the result of adding a new customer.
type AddCustomerResult struct {
	Customer    *models.Customer
	CustomerNum string
}

// Add creates a new customer and returns the generated customer number.
// Equivalent to LGACUS01 COBOL program.
func (s *CustomerService) Add(ctx context.Context, input *AddCustomerInput) (*AddCustomerResult, error) {
	if input == nil {
		return nil, fmt.Errorf("%w: input cannot be nil", ErrInvalidInput)
	}

	// Validate required fields
	if err := s.validateAddInput(input); err != nil {
		return nil, err
	}

	// Convert to repository input
	repoInput := &models.CustomerInput{
		FirstName:    strPtr(input.FirstName),
		LastName:     strPtr(input.LastName),
		DateOfBirth:  input.DateOfBirth,
		HouseName:    strPtr(input.HouseName),
		HouseNumber:  strPtr(input.HouseNumber),
		Postcode:     strPtr(input.Postcode),
		PhoneHome:    strPtr(input.PhoneHome),
		PhoneMobile:  strPtr(input.PhoneMobile),
		EmailAddress: strPtr(input.EmailAddress),
	}

	// Create the customer (repository generates the customer number)
	customer, err := s.customerRepo.Create(ctx, repoInput)
	if err != nil {
		return nil, fmt.Errorf("failed to add customer: %w", err)
	}

	// Increment the customer add counter
	_, _ = s.counterRepo.Increment(ctx, models.CounterCustomerAdd, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return &AddCustomerResult{
		Customer:    customer,
		CustomerNum: customer.CustomerNum,
	}, nil
}

// Get retrieves a customer by their customer number.
// Equivalent to LGICUS01 COBOL program.
func (s *CustomerService) Get(ctx context.Context, customerNum string) (*models.Customer, error) {
	if err := s.validateCustomerNum(customerNum); err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.FindByNum(ctx, customerNum)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Increment the inquiry counter
	_, _ = s.counterRepo.Increment(ctx, models.CounterCustomerInquiry, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return customer, nil
}

// UpdateCustomerInput represents the input for updating a customer.
type UpdateCustomerInput struct {
	FirstName    *string
	LastName     *string
	DateOfBirth  *time.Time
	HouseName    *string
	HouseNumber  *string
	Postcode     *string
	PhoneHome    *string
	PhoneMobile  *string
	EmailAddress *string
}

// Update modifies an existing customer.
// Equivalent to LGUCUS01 COBOL program.
func (s *CustomerService) Update(ctx context.Context, customerNum string, input *UpdateCustomerInput) (*models.Customer, error) {
	if err := s.validateCustomerNum(customerNum); err != nil {
		return nil, err
	}

	if input == nil {
		return nil, fmt.Errorf("%w: input cannot be nil", ErrInvalidInput)
	}

	// Validate input fields
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Check if customer exists
	exists, err := s.customerRepo.Exists(ctx, customerNum)
	if err != nil {
		return nil, fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !exists {
		return nil, ErrCustomerNotFound
	}

	// Convert to repository input
	repoInput := &models.CustomerInput{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		DateOfBirth:  input.DateOfBirth,
		HouseName:    input.HouseName,
		HouseNumber:  input.HouseNumber,
		Postcode:     input.Postcode,
		PhoneHome:    input.PhoneHome,
		PhoneMobile:  input.PhoneMobile,
		EmailAddress: input.EmailAddress,
	}

	customer, err := s.customerRepo.Update(ctx, customerNum, repoInput)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	// Increment the update counter
	_, _ = s.counterRepo.Increment(ctx, models.CounterCustomerUpdate, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return customer, nil
}

// Delete removes a customer by their customer number.
// Note: The original COBOL application did not have a delete customer function,
// but we include it for completeness.
func (s *CustomerService) Delete(ctx context.Context, customerNum string) error {
	if err := s.validateCustomerNum(customerNum); err != nil {
		return err
	}

	err := s.customerRepo.Delete(ctx, customerNum)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return ErrCustomerNotFound
		}
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// List retrieves a paginated list of customers.
func (s *CustomerService) List(ctx context.Context, offset, limit int) ([]*models.Customer, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.customerRepo.List(ctx, offset, limit)
}

// Search finds customers by last name (partial match).
func (s *CustomerService) Search(ctx context.Context, lastName string, limit int) ([]*models.Customer, error) {
	if lastName == "" {
		return nil, fmt.Errorf("%w: last name is required for search", ErrInvalidInput)
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.customerRepo.SearchByLastName(ctx, lastName, limit)
}

// Count returns the total number of customers.
func (s *CustomerService) Count(ctx context.Context) (int64, error) {
	return s.customerRepo.Count(ctx)
}

// validateCustomerNum validates the customer number format.
// Customer numbers are 10-digit numeric strings.
func (s *CustomerService) validateCustomerNum(customerNum string) error {
	if customerNum == "" {
		return fmt.Errorf("%w: customer number is required", ErrInvalidCustomerNumber)
	}

	// Customer number should be 10 digits
	if len(customerNum) != 10 {
		return fmt.Errorf("%w: customer number must be 10 digits", ErrInvalidCustomerNumber)
	}

	// Check if it's numeric
	for _, c := range customerNum {
		if c < '0' || c > '9' {
			return fmt.Errorf("%w: customer number must be numeric", ErrInvalidCustomerNumber)
		}
	}

	return nil
}

// validateAddInput validates the input for adding a customer.
func (s *CustomerService) validateAddInput(input *AddCustomerInput) error {
	var errs []string

	// Last name is required (matching COBOL validation)
	if strings.TrimSpace(input.LastName) == "" {
		errs = append(errs, "last name is required")
	}

	// Validate field lengths (matching COBOL PIC definitions)
	if len(input.FirstName) > 10 {
		errs = append(errs, "first name cannot exceed 10 characters")
	}
	if len(input.LastName) > 20 {
		errs = append(errs, "last name cannot exceed 20 characters")
	}
	if len(input.HouseName) > 20 {
		errs = append(errs, "house name cannot exceed 20 characters")
	}
	if len(input.HouseNumber) > 4 {
		errs = append(errs, "house number cannot exceed 4 characters")
	}
	if len(input.Postcode) > 8 {
		errs = append(errs, "postcode cannot exceed 8 characters")
	}
	if len(input.PhoneHome) > 20 {
		errs = append(errs, "home phone cannot exceed 20 characters")
	}
	if len(input.PhoneMobile) > 20 {
		errs = append(errs, "mobile phone cannot exceed 20 characters")
	}
	if len(input.EmailAddress) > 100 {
		errs = append(errs, "email address cannot exceed 100 characters")
	}

	// Validate email format if provided
	if input.EmailAddress != "" && !isValidEmail(input.EmailAddress) {
		errs = append(errs, "invalid email address format")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidationFailed, strings.Join(errs, "; "))
	}

	return nil
}

// validateUpdateInput validates the input for updating a customer.
func (s *CustomerService) validateUpdateInput(input *UpdateCustomerInput) error {
	var errs []string

	// Validate field lengths if provided
	if input.FirstName != nil && len(*input.FirstName) > 10 {
		errs = append(errs, "first name cannot exceed 10 characters")
	}
	if input.LastName != nil && len(*input.LastName) > 20 {
		errs = append(errs, "last name cannot exceed 20 characters")
	}
	if input.HouseName != nil && len(*input.HouseName) > 20 {
		errs = append(errs, "house name cannot exceed 20 characters")
	}
	if input.HouseNumber != nil && len(*input.HouseNumber) > 4 {
		errs = append(errs, "house number cannot exceed 4 characters")
	}
	if input.Postcode != nil && len(*input.Postcode) > 8 {
		errs = append(errs, "postcode cannot exceed 8 characters")
	}
	if input.PhoneHome != nil && len(*input.PhoneHome) > 20 {
		errs = append(errs, "home phone cannot exceed 20 characters")
	}
	if input.PhoneMobile != nil && len(*input.PhoneMobile) > 20 {
		errs = append(errs, "mobile phone cannot exceed 20 characters")
	}
	if input.EmailAddress != nil && len(*input.EmailAddress) > 100 {
		errs = append(errs, "email address cannot exceed 100 characters")
	}

	// Validate email format if provided
	if input.EmailAddress != nil && *input.EmailAddress != "" && !isValidEmail(*input.EmailAddress) {
		errs = append(errs, "invalid email address format")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidationFailed, strings.Join(errs, "; "))
	}

	return nil
}

// strPtr converts a string to a pointer, returning nil for empty strings.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// isValidEmail performs basic email format validation.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
