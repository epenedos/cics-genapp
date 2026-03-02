package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrPolicyNotFound is returned when a policy is not found.
var ErrPolicyNotFound = errors.New("policy not found")

// ErrInvalidPolicyType is returned when an invalid policy type is provided.
var ErrInvalidPolicyType = errors.New("invalid policy type")

// PolicyRepository provides data access operations for policies.
type PolicyRepository struct {
	db *DB
}

// NewPolicyRepository creates a new PolicyRepository.
func NewPolicyRepository(db *DB) *PolicyRepository {
	return &PolicyRepository{db: db}
}

// Create inserts a new policy with type-specific details.
// Equivalent to LGAPDB01 COBOL program.
func (r *PolicyRepository) Create(ctx context.Context, input *models.PolicyInput) (*models.Policy, error) {
	var policy *models.Policy
	err := r.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		var err error
		policy, err = r.CreateTx(ctx, tx, input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// CreateTx inserts a new policy within a transaction.
func (r *PolicyRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, input *models.PolicyInput) (*models.Policy, error) {
	if !input.PolicyType.IsValid() {
		return nil, ErrInvalidPolicyType
	}

	// Insert master policy record
	query := `
		INSERT INTO policies (
			policy_num, customer_num, policy_type,
			issue_date, expiry_date, broker_id, brokers_ref, payment
		) VALUES (
			next_policy_num(), $1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
	`

	var policy models.Policy
	err := tx.GetContext(ctx, &policy, query,
		input.CustomerNum,
		input.PolicyType,
		toNullTime(input.IssueDate),
		toNullTime(input.ExpiryDate),
		toNullString(input.BrokerID),
		toNullString(input.BrokersRef),
		toNullFloat64(input.Payment),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	// Insert type-specific details
	switch input.PolicyType {
	case models.PolicyTypeMotor:
		if input.Motor != nil {
			motor, err := r.createMotorDetailsTx(ctx, tx, policy.PolicyNum, input.Motor)
			if err != nil {
				return nil, err
			}
			policy.Motor = motor
		}
	case models.PolicyTypeEndowment:
		if input.Endowment != nil {
			endowment, err := r.createEndowmentDetailsTx(ctx, tx, policy.PolicyNum, input.Endowment)
			if err != nil {
				return nil, err
			}
			policy.Endowment = endowment
		}
	case models.PolicyTypeHouse:
		if input.House != nil {
			house, err := r.createHouseDetailsTx(ctx, tx, policy.PolicyNum, input.House)
			if err != nil {
				return nil, err
			}
			policy.House = house
		}
	case models.PolicyTypeCommercial:
		if input.Commercial != nil {
			commercial, err := r.createCommercialDetailsTx(ctx, tx, policy.PolicyNum, input.Commercial)
			if err != nil {
				return nil, err
			}
			policy.Commercial = commercial
		}
	}

	return &policy, nil
}

// createMotorDetailsTx inserts motor policy details.
func (r *PolicyRepository) createMotorDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.MotorInput) (*models.MotorPolicy, error) {
	query := `
		INSERT INTO motor_policies (
			policy_num, make, model, value, reg_number,
			colour, cc, manufactured, premium, accidents
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, policy_num, make, model, value, reg_number,
			colour, cc, manufactured, premium, accidents
	`

	var motor models.MotorPolicy
	err := tx.GetContext(ctx, &motor, query,
		policyNum,
		toNullString(input.Make),
		toNullString(input.Model),
		toNullFloat64(input.Value),
		toNullString(input.RegNumber),
		toNullString(input.Colour),
		toNullInt64(input.CC),
		toNullTime(input.Manufactured),
		toNullFloat64(input.Premium),
		toNullInt64(input.Accidents),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create motor details: %w", err)
	}

	return &motor, nil
}

// createEndowmentDetailsTx inserts endowment policy details.
func (r *PolicyRepository) createEndowmentDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.EndowmentInput) (*models.EndowmentPolicy, error) {
	query := `
		INSERT INTO endowment_policies (
			policy_num, with_profits, equities, managed_fund,
			fund_name, term, sum_assured, life_assured
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, policy_num, with_profits, equities, managed_fund,
			fund_name, term, sum_assured, life_assured
	`

	var endowment models.EndowmentPolicy
	err := tx.GetContext(ctx, &endowment, query,
		policyNum,
		toNullBool(input.WithProfits),
		toNullBool(input.Equities),
		toNullBool(input.ManagedFund),
		toNullString(input.FundName),
		toNullInt64(input.Term),
		toNullFloat64(input.SumAssured),
		toNullString(input.LifeAssured),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create endowment details: %w", err)
	}

	return &endowment, nil
}

// createHouseDetailsTx inserts house policy details.
func (r *PolicyRepository) createHouseDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.HouseInput) (*models.HousePolicy, error) {
	query := `
		INSERT INTO house_policies (
			policy_num, property_type, bedrooms, value,
			house_name, house_number, postcode
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, policy_num, property_type, bedrooms, value,
			house_name, house_number, postcode
	`

	var house models.HousePolicy
	err := tx.GetContext(ctx, &house, query,
		policyNum,
		toNullString(input.PropertyType),
		toNullInt64(input.Bedrooms),
		toNullFloat64(input.Value),
		toNullString(input.HouseName),
		toNullString(input.HouseNumber),
		toNullString(input.Postcode),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create house details: %w", err)
	}

	return &house, nil
}

// createCommercialDetailsTx inserts commercial policy details.
func (r *PolicyRepository) createCommercialDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.CommercialInput) (*models.CommercialPolicy, error) {
	query := `
		INSERT INTO commercial_policies (
			policy_num, address, postcode, latitude, longitude,
			customer, property_type,
			fire_peril, fire_premium, crime_peril, crime_premium,
			flood_peril, flood_premium, weather_peril, weather_premium,
			status, reject_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, policy_num, address, postcode, latitude, longitude,
			customer, property_type,
			fire_peril, fire_premium, crime_peril, crime_premium,
			flood_peril, flood_premium, weather_peril, weather_premium,
			status, reject_reason
	`

	var commercial models.CommercialPolicy
	err := tx.GetContext(ctx, &commercial, query,
		policyNum,
		toNullString(input.Address),
		toNullString(input.Postcode),
		toNullString(input.Latitude),
		toNullString(input.Longitude),
		toNullString(input.Customer),
		toNullString(input.PropertyType),
		toNullInt64(input.FirePeril),
		toNullFloat64(input.FirePremium),
		toNullInt64(input.CrimePeril),
		toNullFloat64(input.CrimePremium),
		toNullInt64(input.FloodPeril),
		toNullFloat64(input.FloodPremium),
		toNullInt64(input.WeatherPeril),
		toNullFloat64(input.WeatherPremium),
		toNullInt64(input.Status),
		toNullString(input.RejectReason),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create commercial details: %w", err)
	}

	return &commercial, nil
}

// FindByNum retrieves a policy by its policy number, including type-specific details.
// Equivalent to LGIPDB01 COBOL program.
func (r *PolicyRepository) FindByNum(ctx context.Context, policyNum string) (*models.Policy, error) {
	return r.FindByNumTx(ctx, r.db, policyNum)
}

// FindByNumTx retrieves a policy within a transaction.
func (r *PolicyRepository) FindByNumTx(ctx context.Context, tx Transactional, policyNum string) (*models.Policy, error) {
	query := `
		SELECT id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
		FROM policies
		WHERE policy_num = $1
	`

	var policy models.Policy
	err := sqlx.GetContext(ctx, tx, &policy, query, policyNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}

	// Load type-specific details
	if err := r.loadPolicyDetails(ctx, tx, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// loadPolicyDetails loads type-specific details for a policy.
func (r *PolicyRepository) loadPolicyDetails(ctx context.Context, tx Transactional, policy *models.Policy) error {
	switch policy.PolicyType {
	case models.PolicyTypeMotor:
		motor, err := r.findMotorDetails(ctx, tx, policy.PolicyNum)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		policy.Motor = motor
	case models.PolicyTypeEndowment:
		endowment, err := r.findEndowmentDetails(ctx, tx, policy.PolicyNum)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		policy.Endowment = endowment
	case models.PolicyTypeHouse:
		house, err := r.findHouseDetails(ctx, tx, policy.PolicyNum)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		policy.House = house
	case models.PolicyTypeCommercial:
		commercial, err := r.findCommercialDetails(ctx, tx, policy.PolicyNum)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		policy.Commercial = commercial
	}
	return nil
}

// findMotorDetails retrieves motor policy details.
func (r *PolicyRepository) findMotorDetails(ctx context.Context, tx Transactional, policyNum string) (*models.MotorPolicy, error) {
	query := `
		SELECT id, policy_num, make, model, value, reg_number,
			colour, cc, manufactured, premium, accidents
		FROM motor_policies
		WHERE policy_num = $1
	`

	var motor models.MotorPolicy
	err := sqlx.GetContext(ctx, tx, &motor, query, policyNum)
	if err != nil {
		return nil, err
	}
	return &motor, nil
}

// findEndowmentDetails retrieves endowment policy details.
func (r *PolicyRepository) findEndowmentDetails(ctx context.Context, tx Transactional, policyNum string) (*models.EndowmentPolicy, error) {
	query := `
		SELECT id, policy_num, with_profits, equities, managed_fund,
			fund_name, term, sum_assured, life_assured
		FROM endowment_policies
		WHERE policy_num = $1
	`

	var endowment models.EndowmentPolicy
	err := sqlx.GetContext(ctx, tx, &endowment, query, policyNum)
	if err != nil {
		return nil, err
	}
	return &endowment, nil
}

// findHouseDetails retrieves house policy details.
func (r *PolicyRepository) findHouseDetails(ctx context.Context, tx Transactional, policyNum string) (*models.HousePolicy, error) {
	query := `
		SELECT id, policy_num, property_type, bedrooms, value,
			house_name, house_number, postcode
		FROM house_policies
		WHERE policy_num = $1
	`

	var house models.HousePolicy
	err := sqlx.GetContext(ctx, tx, &house, query, policyNum)
	if err != nil {
		return nil, err
	}
	return &house, nil
}

// findCommercialDetails retrieves commercial policy details.
func (r *PolicyRepository) findCommercialDetails(ctx context.Context, tx Transactional, policyNum string) (*models.CommercialPolicy, error) {
	query := `
		SELECT id, policy_num, address, postcode, latitude, longitude,
			customer, property_type,
			fire_peril, fire_premium, crime_peril, crime_premium,
			flood_peril, flood_premium, weather_peril, weather_premium,
			status, reject_reason
		FROM commercial_policies
		WHERE policy_num = $1
	`

	var commercial models.CommercialPolicy
	err := sqlx.GetContext(ctx, tx, &commercial, query, policyNum)
	if err != nil {
		return nil, err
	}
	return &commercial, nil
}

// Update updates an existing policy and its type-specific details.
// Equivalent to LGUPDB01 COBOL program.
func (r *PolicyRepository) Update(ctx context.Context, policyNum string, input *models.PolicyInput) (*models.Policy, error) {
	var policy *models.Policy
	err := r.db.WithTx(ctx, func(tx *sqlx.Tx) error {
		var err error
		policy, err = r.UpdateTx(ctx, tx, policyNum, input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// UpdateTx updates a policy within a transaction.
func (r *PolicyRepository) UpdateTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.PolicyInput) (*models.Policy, error) {
	query := `
		UPDATE policies SET
			issue_date = COALESCE($2, issue_date),
			expiry_date = COALESCE($3, expiry_date),
			broker_id = COALESCE($4, broker_id),
			brokers_ref = COALESCE($5, brokers_ref),
			payment = COALESCE($6, payment),
			last_changed = CURRENT_TIMESTAMP
		WHERE policy_num = $1
		RETURNING id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
	`

	var policy models.Policy
	err := tx.GetContext(ctx, &policy, query,
		policyNum,
		toNullTime(input.IssueDate),
		toNullTime(input.ExpiryDate),
		toNullString(input.BrokerID),
		toNullString(input.BrokersRef),
		toNullFloat64(input.Payment),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}

	// Update type-specific details
	switch policy.PolicyType {
	case models.PolicyTypeMotor:
		if input.Motor != nil {
			motor, err := r.updateMotorDetailsTx(ctx, tx, policyNum, input.Motor)
			if err != nil {
				return nil, err
			}
			policy.Motor = motor
		}
	case models.PolicyTypeEndowment:
		if input.Endowment != nil {
			endowment, err := r.updateEndowmentDetailsTx(ctx, tx, policyNum, input.Endowment)
			if err != nil {
				return nil, err
			}
			policy.Endowment = endowment
		}
	case models.PolicyTypeHouse:
		if input.House != nil {
			house, err := r.updateHouseDetailsTx(ctx, tx, policyNum, input.House)
			if err != nil {
				return nil, err
			}
			policy.House = house
		}
	case models.PolicyTypeCommercial:
		if input.Commercial != nil {
			commercial, err := r.updateCommercialDetailsTx(ctx, tx, policyNum, input.Commercial)
			if err != nil {
				return nil, err
			}
			policy.Commercial = commercial
		}
	}

	return &policy, nil
}

// updateMotorDetailsTx updates motor policy details.
func (r *PolicyRepository) updateMotorDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.MotorInput) (*models.MotorPolicy, error) {
	query := `
		UPDATE motor_policies SET
			make = COALESCE($2, make),
			model = COALESCE($3, model),
			value = COALESCE($4, value),
			reg_number = COALESCE($5, reg_number),
			colour = COALESCE($6, colour),
			cc = COALESCE($7, cc),
			manufactured = COALESCE($8, manufactured),
			premium = COALESCE($9, premium),
			accidents = COALESCE($10, accidents)
		WHERE policy_num = $1
		RETURNING id, policy_num, make, model, value, reg_number,
			colour, cc, manufactured, premium, accidents
	`

	var motor models.MotorPolicy
	err := tx.GetContext(ctx, &motor, query,
		policyNum,
		toNullString(input.Make),
		toNullString(input.Model),
		toNullFloat64(input.Value),
		toNullString(input.RegNumber),
		toNullString(input.Colour),
		toNullInt64(input.CC),
		toNullTime(input.Manufactured),
		toNullFloat64(input.Premium),
		toNullInt64(input.Accidents),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update motor details: %w", err)
	}

	return &motor, nil
}

// updateEndowmentDetailsTx updates endowment policy details.
func (r *PolicyRepository) updateEndowmentDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.EndowmentInput) (*models.EndowmentPolicy, error) {
	query := `
		UPDATE endowment_policies SET
			with_profits = COALESCE($2, with_profits),
			equities = COALESCE($3, equities),
			managed_fund = COALESCE($4, managed_fund),
			fund_name = COALESCE($5, fund_name),
			term = COALESCE($6, term),
			sum_assured = COALESCE($7, sum_assured),
			life_assured = COALESCE($8, life_assured)
		WHERE policy_num = $1
		RETURNING id, policy_num, with_profits, equities, managed_fund,
			fund_name, term, sum_assured, life_assured
	`

	var endowment models.EndowmentPolicy
	err := tx.GetContext(ctx, &endowment, query,
		policyNum,
		toNullBool(input.WithProfits),
		toNullBool(input.Equities),
		toNullBool(input.ManagedFund),
		toNullString(input.FundName),
		toNullInt64(input.Term),
		toNullFloat64(input.SumAssured),
		toNullString(input.LifeAssured),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update endowment details: %w", err)
	}

	return &endowment, nil
}

// updateHouseDetailsTx updates house policy details.
func (r *PolicyRepository) updateHouseDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.HouseInput) (*models.HousePolicy, error) {
	query := `
		UPDATE house_policies SET
			property_type = COALESCE($2, property_type),
			bedrooms = COALESCE($3, bedrooms),
			value = COALESCE($4, value),
			house_name = COALESCE($5, house_name),
			house_number = COALESCE($6, house_number),
			postcode = COALESCE($7, postcode)
		WHERE policy_num = $1
		RETURNING id, policy_num, property_type, bedrooms, value,
			house_name, house_number, postcode
	`

	var house models.HousePolicy
	err := tx.GetContext(ctx, &house, query,
		policyNum,
		toNullString(input.PropertyType),
		toNullInt64(input.Bedrooms),
		toNullFloat64(input.Value),
		toNullString(input.HouseName),
		toNullString(input.HouseNumber),
		toNullString(input.Postcode),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update house details: %w", err)
	}

	return &house, nil
}

// updateCommercialDetailsTx updates commercial policy details.
func (r *PolicyRepository) updateCommercialDetailsTx(ctx context.Context, tx *sqlx.Tx, policyNum string, input *models.CommercialInput) (*models.CommercialPolicy, error) {
	query := `
		UPDATE commercial_policies SET
			address = COALESCE($2, address),
			postcode = COALESCE($3, postcode),
			latitude = COALESCE($4, latitude),
			longitude = COALESCE($5, longitude),
			customer = COALESCE($6, customer),
			property_type = COALESCE($7, property_type),
			fire_peril = COALESCE($8, fire_peril),
			fire_premium = COALESCE($9, fire_premium),
			crime_peril = COALESCE($10, crime_peril),
			crime_premium = COALESCE($11, crime_premium),
			flood_peril = COALESCE($12, flood_peril),
			flood_premium = COALESCE($13, flood_premium),
			weather_peril = COALESCE($14, weather_peril),
			weather_premium = COALESCE($15, weather_premium),
			status = COALESCE($16, status),
			reject_reason = COALESCE($17, reject_reason)
		WHERE policy_num = $1
		RETURNING id, policy_num, address, postcode, latitude, longitude,
			customer, property_type,
			fire_peril, fire_premium, crime_peril, crime_premium,
			flood_peril, flood_premium, weather_peril, weather_premium,
			status, reject_reason
	`

	var commercial models.CommercialPolicy
	err := tx.GetContext(ctx, &commercial, query,
		policyNum,
		toNullString(input.Address),
		toNullString(input.Postcode),
		toNullString(input.Latitude),
		toNullString(input.Longitude),
		toNullString(input.Customer),
		toNullString(input.PropertyType),
		toNullInt64(input.FirePeril),
		toNullFloat64(input.FirePremium),
		toNullInt64(input.CrimePeril),
		toNullFloat64(input.CrimePremium),
		toNullInt64(input.FloodPeril),
		toNullFloat64(input.FloodPremium),
		toNullInt64(input.WeatherPeril),
		toNullFloat64(input.WeatherPremium),
		toNullInt64(input.Status),
		toNullString(input.RejectReason),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update commercial details: %w", err)
	}

	return &commercial, nil
}

// Delete removes a policy and its type-specific details (cascaded by FK).
// Equivalent to LGDPDB01 COBOL program.
func (r *PolicyRepository) Delete(ctx context.Context, policyNum string) error {
	return r.DeleteTx(ctx, r.db, policyNum)
}

// DeleteTx removes a policy within a transaction.
func (r *PolicyRepository) DeleteTx(ctx context.Context, tx Transactional, policyNum string) error {
	query := `DELETE FROM policies WHERE policy_num = $1`

	result, err := tx.ExecContext(ctx, query, policyNum)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPolicyNotFound
	}

	return nil
}

// FindByCustomerNum retrieves all policies for a customer.
func (r *PolicyRepository) FindByCustomerNum(ctx context.Context, customerNum string) ([]*models.Policy, error) {
	query := `
		SELECT id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
		FROM policies
		WHERE customer_num = $1
		ORDER BY policy_num
	`

	var policies []*models.Policy
	err := r.db.SelectContext(ctx, &policies, query, customerNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find policies by customer: %w", err)
	}

	// Load details for each policy
	for _, policy := range policies {
		if err := r.loadPolicyDetails(ctx, r.db, policy); err != nil {
			return nil, err
		}
	}

	return policies, nil
}

// List retrieves a paginated list of policies.
func (r *PolicyRepository) List(ctx context.Context, offset, limit int) ([]*models.Policy, error) {
	query := `
		SELECT id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
		FROM policies
		ORDER BY policy_num
		LIMIT $1 OFFSET $2
	`

	var policies []*models.Policy
	err := r.db.SelectContext(ctx, &policies, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list policies: %w", err)
	}

	return policies, nil
}

// ListByType retrieves policies filtered by type.
func (r *PolicyRepository) ListByType(ctx context.Context, policyType models.PolicyType, offset, limit int) ([]*models.Policy, error) {
	query := `
		SELECT id, policy_num, customer_num, policy_type,
			issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment,
			created_at, updated_at
		FROM policies
		WHERE policy_type = $1
		ORDER BY policy_num
		LIMIT $2 OFFSET $3
	`

	var policies []*models.Policy
	err := r.db.SelectContext(ctx, &policies, query, policyType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list policies by type: %w", err)
	}

	// Load details for each policy
	for _, policy := range policies {
		if err := r.loadPolicyDetails(ctx, r.db, policy); err != nil {
			return nil, err
		}
	}

	return policies, nil
}

// Count returns the total number of policies.
func (r *PolicyRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM policies`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count policies: %w", err)
	}

	return count, nil
}

// CountByType returns the number of policies of a specific type.
func (r *PolicyRepository) CountByType(ctx context.Context, policyType models.PolicyType) (int64, error) {
	query := `SELECT COUNT(*) FROM policies WHERE policy_type = $1`

	var count int64
	err := r.db.GetContext(ctx, &count, query, policyType)
	if err != nil {
		return 0, fmt.Errorf("failed to count policies by type: %w", err)
	}

	return count, nil
}

// Exists checks if a policy exists by its policy number.
func (r *PolicyRepository) Exists(ctx context.Context, policyNum string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM policies WHERE policy_num = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, policyNum)
	if err != nil {
		return false, fmt.Errorf("failed to check policy existence: %w", err)
	}

	return exists, nil
}
