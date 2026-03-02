package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/cicsdev/genapp/internal/repository"
)

// Policy service errors.
var (
	ErrPolicyNotFound       = errors.New("policy not found")
	ErrInvalidPolicyNumber  = errors.New("invalid policy number")
	ErrInvalidPolicyType    = errors.New("invalid policy type")
	ErrCustomerRequired     = errors.New("customer number is required")
	ErrPolicyTypeRequired   = errors.New("policy type is required")
	ErrPolicyDetailsRequired = errors.New("policy type-specific details are required")
)

// PolicyService provides business logic for policy operations.
// Equivalent to COBOL programs: LGAPOL01 (add), LGIPOL01 (inquire),
// LGUPOL01 (update), LGDPOL01 (delete).
type PolicyService struct {
	policyRepo   *repository.PolicyRepository
	customerRepo *repository.CustomerRepository
	counterRepo  *repository.CounterRepository
}

// NewPolicyService creates a new PolicyService.
func NewPolicyService(
	policyRepo *repository.PolicyRepository,
	customerRepo *repository.CustomerRepository,
	counterRepo *repository.CounterRepository,
) *PolicyService {
	return &PolicyService{
		policyRepo:   policyRepo,
		customerRepo: customerRepo,
		counterRepo:  counterRepo,
	}
}

// AddPolicyInput represents the input for adding a new policy.
type AddPolicyInput struct {
	CustomerNum string
	PolicyType  models.PolicyType
	IssueDate   *time.Time
	ExpiryDate  *time.Time
	BrokerID    string
	BrokersRef  string
	Payment     *float64

	// Type-specific details (exactly one should be provided based on PolicyType)
	Motor      *AddMotorInput
	Endowment  *AddEndowmentInput
	House      *AddHouseInput
	Commercial *AddCommercialInput
}

// AddMotorInput represents motor policy details for add operation.
type AddMotorInput struct {
	Make         string
	Model        string
	Value        float64
	RegNumber    string
	Colour       string
	CC           int
	Manufactured *time.Time
	Premium      float64
	Accidents    int
}

// AddEndowmentInput represents endowment policy details for add operation.
type AddEndowmentInput struct {
	WithProfits bool
	Equities    bool
	ManagedFund bool
	FundName    string
	Term        int
	SumAssured  float64
	LifeAssured string
}

// AddHouseInput represents house policy details for add operation.
type AddHouseInput struct {
	PropertyType string
	Bedrooms     int
	Value        float64
	HouseName    string
	HouseNumber  string
	Postcode     string
}

// AddCommercialInput represents commercial policy details for add operation.
type AddCommercialInput struct {
	Address        string
	Postcode       string
	Latitude       string
	Longitude      string
	Customer       string
	PropertyType   string
	FirePeril      int
	FirePremium    float64
	CrimePeril     int
	CrimePremium   float64
	FloodPeril     int
	FloodPremium   float64
	WeatherPeril   int
	WeatherPremium float64
	Status         int
	RejectReason   string
}

// AddPolicyResult represents the result of adding a new policy.
type AddPolicyResult struct {
	Policy    *models.Policy
	PolicyNum string
}

// Add creates a new policy for an existing customer.
// Equivalent to LGAPOL01 COBOL program.
func (s *PolicyService) Add(ctx context.Context, input *AddPolicyInput) (*AddPolicyResult, error) {
	if input == nil {
		return nil, fmt.Errorf("%w: input cannot be nil", ErrInvalidInput)
	}

	// Validate input
	if err := s.validateAddInput(input); err != nil {
		return nil, err
	}

	// Verify customer exists
	exists, err := s.customerRepo.Exists(ctx, input.CustomerNum)
	if err != nil {
		return nil, fmt.Errorf("failed to verify customer: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("%w: customer %s", ErrCustomerNotFound, input.CustomerNum)
	}

	// Convert to repository input
	repoInput := s.buildPolicyInput(input)

	// Create the policy
	policy, err := s.policyRepo.Create(ctx, repoInput)
	if err != nil {
		return nil, fmt.Errorf("failed to add policy: %w", err)
	}

	// Increment counters
	_, _ = s.counterRepo.Increment(ctx, models.CounterPolicyAdd, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return &AddPolicyResult{
		Policy:    policy,
		PolicyNum: policy.PolicyNum,
	}, nil
}

// Get retrieves a policy by its policy number, including type-specific details.
// Equivalent to LGIPOL01 COBOL program.
func (s *PolicyService) Get(ctx context.Context, policyNum string) (*models.Policy, error) {
	if err := s.validatePolicyNum(policyNum); err != nil {
		return nil, err
	}

	policy, err := s.policyRepo.FindByNum(ctx, policyNum)
	if err != nil {
		if errors.Is(err, repository.ErrPolicyNotFound) {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	// Increment counters
	_, _ = s.counterRepo.Increment(ctx, models.CounterPolicyInquiry, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return policy, nil
}

// UpdatePolicyInput represents the input for updating a policy.
type UpdatePolicyInput struct {
	IssueDate  *time.Time
	ExpiryDate *time.Time
	BrokerID   *string
	BrokersRef *string
	Payment    *float64

	// Type-specific details (optional)
	Motor      *UpdateMotorInput
	Endowment  *UpdateEndowmentInput
	House      *UpdateHouseInput
	Commercial *UpdateCommercialInput
}

// UpdateMotorInput represents motor policy details for update operation.
type UpdateMotorInput struct {
	Make         *string
	Model        *string
	Value        *float64
	RegNumber    *string
	Colour       *string
	CC           *int
	Manufactured *time.Time
	Premium      *float64
	Accidents    *int
}

// UpdateEndowmentInput represents endowment policy details for update operation.
type UpdateEndowmentInput struct {
	WithProfits *bool
	Equities    *bool
	ManagedFund *bool
	FundName    *string
	Term        *int
	SumAssured  *float64
	LifeAssured *string
}

// UpdateHouseInput represents house policy details for update operation.
type UpdateHouseInput struct {
	PropertyType *string
	Bedrooms     *int
	Value        *float64
	HouseName    *string
	HouseNumber  *string
	Postcode     *string
}

// UpdateCommercialInput represents commercial policy details for update operation.
type UpdateCommercialInput struct {
	Address        *string
	Postcode       *string
	Latitude       *string
	Longitude      *string
	Customer       *string
	PropertyType   *string
	FirePeril      *int
	FirePremium    *float64
	CrimePeril     *int
	CrimePremium   *float64
	FloodPeril     *int
	FloodPremium   *float64
	WeatherPeril   *int
	WeatherPremium *float64
	Status         *int
	RejectReason   *string
}

// Update modifies an existing policy.
// Equivalent to LGUPOL01 COBOL program.
func (s *PolicyService) Update(ctx context.Context, policyNum string, input *UpdatePolicyInput) (*models.Policy, error) {
	if err := s.validatePolicyNum(policyNum); err != nil {
		return nil, err
	}

	if input == nil {
		return nil, fmt.Errorf("%w: input cannot be nil", ErrInvalidInput)
	}

	// Verify policy exists and get its type
	existing, err := s.policyRepo.FindByNum(ctx, policyNum)
	if err != nil {
		if errors.Is(err, repository.ErrPolicyNotFound) {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}

	// Convert to repository input
	repoInput := s.buildUpdateInput(existing.PolicyType, input)

	policy, err := s.policyRepo.Update(ctx, policyNum, repoInput)
	if err != nil {
		if errors.Is(err, repository.ErrPolicyNotFound) {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}

	// Increment counters
	_, _ = s.counterRepo.Increment(ctx, models.CounterPolicyUpdate, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return policy, nil
}

// Delete removes a policy by its policy number.
// Equivalent to LGDPOL01 COBOL program.
func (s *PolicyService) Delete(ctx context.Context, policyNum string) error {
	if err := s.validatePolicyNum(policyNum); err != nil {
		return err
	}

	err := s.policyRepo.Delete(ctx, policyNum)
	if err != nil {
		if errors.Is(err, repository.ErrPolicyNotFound) {
			return ErrPolicyNotFound
		}
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	// Increment counters
	_, _ = s.counterRepo.Increment(ctx, models.CounterPolicyDelete, 1)
	_, _ = s.counterRepo.Increment(ctx, models.CounterTotalTransactions, 1)

	return nil
}

// GetByCustomer retrieves all policies for a customer.
func (s *PolicyService) GetByCustomer(ctx context.Context, customerNum string) ([]*models.Policy, error) {
	// Validate customer number
	if customerNum == "" {
		return nil, fmt.Errorf("%w: customer number is required", ErrInvalidCustomerNumber)
	}
	if len(customerNum) != 10 {
		return nil, fmt.Errorf("%w: customer number must be 10 digits", ErrInvalidCustomerNumber)
	}

	return s.policyRepo.FindByCustomerNum(ctx, customerNum)
}

// List retrieves a paginated list of policies.
func (s *PolicyService) List(ctx context.Context, offset, limit int) ([]*models.Policy, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.policyRepo.List(ctx, offset, limit)
}

// ListByType retrieves policies filtered by type.
func (s *PolicyService) ListByType(ctx context.Context, policyType models.PolicyType, offset, limit int) ([]*models.Policy, error) {
	if !policyType.IsValid() {
		return nil, ErrInvalidPolicyType
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.policyRepo.ListByType(ctx, policyType, offset, limit)
}

// Count returns the total number of policies.
func (s *PolicyService) Count(ctx context.Context) (int64, error) {
	return s.policyRepo.Count(ctx)
}

// CountByType returns the number of policies of a specific type.
func (s *PolicyService) CountByType(ctx context.Context, policyType models.PolicyType) (int64, error) {
	if !policyType.IsValid() {
		return 0, ErrInvalidPolicyType
	}
	return s.policyRepo.CountByType(ctx, policyType)
}

// validatePolicyNum validates the policy number format.
// Policy numbers are 10-digit numeric strings.
func (s *PolicyService) validatePolicyNum(policyNum string) error {
	if policyNum == "" {
		return fmt.Errorf("%w: policy number is required", ErrInvalidPolicyNumber)
	}

	if len(policyNum) != 10 {
		return fmt.Errorf("%w: policy number must be 10 digits", ErrInvalidPolicyNumber)
	}

	for _, c := range policyNum {
		if c < '0' || c > '9' {
			return fmt.Errorf("%w: policy number must be numeric", ErrInvalidPolicyNumber)
		}
	}

	return nil
}

// validateAddInput validates the input for adding a policy.
func (s *PolicyService) validateAddInput(input *AddPolicyInput) error {
	var errs []string

	// Customer number is required
	if input.CustomerNum == "" {
		errs = append(errs, "customer number is required")
	} else if len(input.CustomerNum) != 10 {
		errs = append(errs, "customer number must be 10 digits")
	}

	// Policy type is required and must be valid
	if !input.PolicyType.IsValid() {
		errs = append(errs, "valid policy type is required (E, H, M, or C)")
	}

	// Type-specific details must match policy type
	switch input.PolicyType {
	case models.PolicyTypeMotor:
		if input.Motor == nil {
			errs = append(errs, "motor policy details are required for motor policies")
		} else {
			if input.Motor.Make == "" {
				errs = append(errs, "car make is required")
			}
			if len(input.Motor.Make) > 15 {
				errs = append(errs, "car make cannot exceed 15 characters")
			}
			if len(input.Motor.Model) > 15 {
				errs = append(errs, "car model cannot exceed 15 characters")
			}
			if len(input.Motor.RegNumber) > 7 {
				errs = append(errs, "registration number cannot exceed 7 characters")
			}
			if len(input.Motor.Colour) > 8 {
				errs = append(errs, "colour cannot exceed 8 characters")
			}
		}
	case models.PolicyTypeEndowment:
		if input.Endowment == nil {
			errs = append(errs, "endowment policy details are required for endowment policies")
		} else {
			if len(input.Endowment.FundName) > 10 {
				errs = append(errs, "fund name cannot exceed 10 characters")
			}
			if len(input.Endowment.LifeAssured) > 31 {
				errs = append(errs, "life assured cannot exceed 31 characters")
			}
			if input.Endowment.Term < 0 || input.Endowment.Term > 99 {
				errs = append(errs, "term must be between 0 and 99 years")
			}
		}
	case models.PolicyTypeHouse:
		if input.House == nil {
			errs = append(errs, "house policy details are required for house policies")
		} else {
			if len(input.House.PropertyType) > 15 {
				errs = append(errs, "property type cannot exceed 15 characters")
			}
			if len(input.House.HouseName) > 20 {
				errs = append(errs, "house name cannot exceed 20 characters")
			}
			if len(input.House.HouseNumber) > 4 {
				errs = append(errs, "house number cannot exceed 4 characters")
			}
			if len(input.House.Postcode) > 8 {
				errs = append(errs, "postcode cannot exceed 8 characters")
			}
		}
	case models.PolicyTypeCommercial:
		if input.Commercial == nil {
			errs = append(errs, "commercial policy details are required for commercial policies")
		}
	}

	// Validate field lengths for common fields
	if len(input.BrokerID) > 10 {
		errs = append(errs, "broker ID cannot exceed 10 characters")
	}
	if len(input.BrokersRef) > 10 {
		errs = append(errs, "brokers ref cannot exceed 10 characters")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrValidationFailed, strings.Join(errs, "; "))
	}

	return nil
}

// buildPolicyInput converts AddPolicyInput to repository PolicyInput.
func (s *PolicyService) buildPolicyInput(input *AddPolicyInput) *models.PolicyInput {
	pi := &models.PolicyInput{
		CustomerNum: input.CustomerNum,
		PolicyType:  input.PolicyType,
		IssueDate:   input.IssueDate,
		ExpiryDate:  input.ExpiryDate,
		BrokerID:    strPtr(input.BrokerID),
		BrokersRef:  strPtr(input.BrokersRef),
		Payment:     input.Payment,
	}

	switch input.PolicyType {
	case models.PolicyTypeMotor:
		if input.Motor != nil {
			pi.Motor = &models.MotorInput{
				Make:         strPtr(input.Motor.Make),
				Model:        strPtr(input.Motor.Model),
				Value:        float64Ptr(input.Motor.Value),
				RegNumber:    strPtr(input.Motor.RegNumber),
				Colour:       strPtr(input.Motor.Colour),
				CC:           intPtr(input.Motor.CC),
				Manufactured: input.Motor.Manufactured,
				Premium:      float64Ptr(input.Motor.Premium),
				Accidents:    intPtr(input.Motor.Accidents),
			}
		}
	case models.PolicyTypeEndowment:
		if input.Endowment != nil {
			pi.Endowment = &models.EndowmentInput{
				WithProfits: boolPtr(input.Endowment.WithProfits),
				Equities:    boolPtr(input.Endowment.Equities),
				ManagedFund: boolPtr(input.Endowment.ManagedFund),
				FundName:    strPtr(input.Endowment.FundName),
				Term:        intPtr(input.Endowment.Term),
				SumAssured:  float64Ptr(input.Endowment.SumAssured),
				LifeAssured: strPtr(input.Endowment.LifeAssured),
			}
		}
	case models.PolicyTypeHouse:
		if input.House != nil {
			pi.House = &models.HouseInput{
				PropertyType: strPtr(input.House.PropertyType),
				Bedrooms:     intPtr(input.House.Bedrooms),
				Value:        float64Ptr(input.House.Value),
				HouseName:    strPtr(input.House.HouseName),
				HouseNumber:  strPtr(input.House.HouseNumber),
				Postcode:     strPtr(input.House.Postcode),
			}
		}
	case models.PolicyTypeCommercial:
		if input.Commercial != nil {
			pi.Commercial = &models.CommercialInput{
				Address:        strPtr(input.Commercial.Address),
				Postcode:       strPtr(input.Commercial.Postcode),
				Latitude:       strPtr(input.Commercial.Latitude),
				Longitude:      strPtr(input.Commercial.Longitude),
				Customer:       strPtr(input.Commercial.Customer),
				PropertyType:   strPtr(input.Commercial.PropertyType),
				FirePeril:      intPtr(input.Commercial.FirePeril),
				FirePremium:    float64Ptr(input.Commercial.FirePremium),
				CrimePeril:     intPtr(input.Commercial.CrimePeril),
				CrimePremium:   float64Ptr(input.Commercial.CrimePremium),
				FloodPeril:     intPtr(input.Commercial.FloodPeril),
				FloodPremium:   float64Ptr(input.Commercial.FloodPremium),
				WeatherPeril:   intPtr(input.Commercial.WeatherPeril),
				WeatherPremium: float64Ptr(input.Commercial.WeatherPremium),
				Status:         intPtr(input.Commercial.Status),
				RejectReason:   strPtr(input.Commercial.RejectReason),
			}
		}
	}

	return pi
}

// buildUpdateInput converts UpdatePolicyInput to repository PolicyInput.
func (s *PolicyService) buildUpdateInput(policyType models.PolicyType, input *UpdatePolicyInput) *models.PolicyInput {
	pi := &models.PolicyInput{
		PolicyType: policyType,
		IssueDate:  input.IssueDate,
		ExpiryDate: input.ExpiryDate,
		BrokerID:   input.BrokerID,
		BrokersRef: input.BrokersRef,
		Payment:    input.Payment,
	}

	switch policyType {
	case models.PolicyTypeMotor:
		if input.Motor != nil {
			pi.Motor = &models.MotorInput{
				Make:         input.Motor.Make,
				Model:        input.Motor.Model,
				Value:        input.Motor.Value,
				RegNumber:    input.Motor.RegNumber,
				Colour:       input.Motor.Colour,
				CC:           input.Motor.CC,
				Manufactured: input.Motor.Manufactured,
				Premium:      input.Motor.Premium,
				Accidents:    input.Motor.Accidents,
			}
		}
	case models.PolicyTypeEndowment:
		if input.Endowment != nil {
			pi.Endowment = &models.EndowmentInput{
				WithProfits: input.Endowment.WithProfits,
				Equities:    input.Endowment.Equities,
				ManagedFund: input.Endowment.ManagedFund,
				FundName:    input.Endowment.FundName,
				Term:        input.Endowment.Term,
				SumAssured:  input.Endowment.SumAssured,
				LifeAssured: input.Endowment.LifeAssured,
			}
		}
	case models.PolicyTypeHouse:
		if input.House != nil {
			pi.House = &models.HouseInput{
				PropertyType: input.House.PropertyType,
				Bedrooms:     input.House.Bedrooms,
				Value:        input.House.Value,
				HouseName:    input.House.HouseName,
				HouseNumber:  input.House.HouseNumber,
				Postcode:     input.House.Postcode,
			}
		}
	case models.PolicyTypeCommercial:
		if input.Commercial != nil {
			pi.Commercial = &models.CommercialInput{
				Address:        input.Commercial.Address,
				Postcode:       input.Commercial.Postcode,
				Latitude:       input.Commercial.Latitude,
				Longitude:      input.Commercial.Longitude,
				Customer:       input.Commercial.Customer,
				PropertyType:   input.Commercial.PropertyType,
				FirePeril:      input.Commercial.FirePeril,
				FirePremium:    input.Commercial.FirePremium,
				CrimePeril:     input.Commercial.CrimePeril,
				CrimePremium:   input.Commercial.CrimePremium,
				FloodPeril:     input.Commercial.FloodPeril,
				FloodPremium:   input.Commercial.FloodPremium,
				WeatherPeril:   input.Commercial.WeatherPeril,
				WeatherPremium: input.Commercial.WeatherPremium,
				Status:         input.Commercial.Status,
				RejectReason:   input.Commercial.RejectReason,
			}
		}
	}

	return pi
}

// Helper functions for pointer conversion
func float64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
