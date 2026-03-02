package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cicsdev/genapp/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrClaimNotFound is returned when a claim is not found.
var ErrClaimNotFound = errors.New("claim not found")

// ClaimRepository provides data access operations for claims.
type ClaimRepository struct {
	db *DB
}

// NewClaimRepository creates a new ClaimRepository.
func NewClaimRepository(db *DB) *ClaimRepository {
	return &ClaimRepository{db: db}
}

// Create inserts a new claim and returns the created claim with the generated number.
func (r *ClaimRepository) Create(ctx context.Context, input *models.ClaimInput) (*models.Claim, error) {
	return r.CreateTx(ctx, r.db, input)
}

// CreateTx inserts a new claim within a transaction.
func (r *ClaimRepository) CreateTx(ctx context.Context, tx Transactional, input *models.ClaimInput) (*models.Claim, error) {
	query := `
		INSERT INTO claims (
			claim_num, policy_num, claim_date, paid, value, cause, observations
		) VALUES (
			next_claim_num(), $1, $2, $3, $4, $5, $6
		)
		RETURNING id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
	`

	var claim models.Claim
	err := sqlx.GetContext(ctx, tx, &claim, query,
		input.PolicyNum,
		toNullTime(input.ClaimDate),
		toNullFloat64(input.Paid),
		toNullFloat64(input.Value),
		toNullString(input.Cause),
		toNullString(input.Observations),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}

	return &claim, nil
}

// FindByNum retrieves a claim by its claim number.
func (r *ClaimRepository) FindByNum(ctx context.Context, claimNum string) (*models.Claim, error) {
	return r.FindByNumTx(ctx, r.db, claimNum)
}

// FindByNumTx retrieves a claim within a transaction.
func (r *ClaimRepository) FindByNumTx(ctx context.Context, tx Transactional, claimNum string) (*models.Claim, error) {
	query := `
		SELECT id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
		FROM claims
		WHERE claim_num = $1
	`

	var claim models.Claim
	err := sqlx.GetContext(ctx, tx, &claim, query, claimNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClaimNotFound
		}
		return nil, fmt.Errorf("failed to find claim: %w", err)
	}

	return &claim, nil
}

// FindByID retrieves a claim by its internal ID.
func (r *ClaimRepository) FindByID(ctx context.Context, id int64) (*models.Claim, error) {
	query := `
		SELECT id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
		FROM claims
		WHERE id = $1
	`

	var claim models.Claim
	err := r.db.GetContext(ctx, &claim, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClaimNotFound
		}
		return nil, fmt.Errorf("failed to find claim: %w", err)
	}

	return &claim, nil
}

// Update updates an existing claim.
func (r *ClaimRepository) Update(ctx context.Context, claimNum string, input *models.ClaimInput) (*models.Claim, error) {
	return r.UpdateTx(ctx, r.db, claimNum, input)
}

// UpdateTx updates a claim within a transaction.
func (r *ClaimRepository) UpdateTx(ctx context.Context, tx Transactional, claimNum string, input *models.ClaimInput) (*models.Claim, error) {
	query := `
		UPDATE claims SET
			claim_date = COALESCE($2, claim_date),
			paid = COALESCE($3, paid),
			value = COALESCE($4, value),
			cause = COALESCE($5, cause),
			observations = COALESCE($6, observations)
		WHERE claim_num = $1
		RETURNING id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
	`

	var claim models.Claim
	err := sqlx.GetContext(ctx, tx, &claim, query,
		claimNum,
		toNullTime(input.ClaimDate),
		toNullFloat64(input.Paid),
		toNullFloat64(input.Value),
		toNullString(input.Cause),
		toNullString(input.Observations),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClaimNotFound
		}
		return nil, fmt.Errorf("failed to update claim: %w", err)
	}

	return &claim, nil
}

// Delete removes a claim by its claim number.
func (r *ClaimRepository) Delete(ctx context.Context, claimNum string) error {
	return r.DeleteTx(ctx, r.db, claimNum)
}

// DeleteTx removes a claim within a transaction.
func (r *ClaimRepository) DeleteTx(ctx context.Context, tx Transactional, claimNum string) error {
	query := `DELETE FROM claims WHERE claim_num = $1`

	result, err := tx.ExecContext(ctx, query, claimNum)
	if err != nil {
		return fmt.Errorf("failed to delete claim: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrClaimNotFound
	}

	return nil
}

// FindByPolicyNum retrieves all claims for a policy.
func (r *ClaimRepository) FindByPolicyNum(ctx context.Context, policyNum string) ([]*models.Claim, error) {
	query := `
		SELECT id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
		FROM claims
		WHERE policy_num = $1
		ORDER BY claim_date DESC, claim_num
	`

	var claims []*models.Claim
	err := r.db.SelectContext(ctx, &claims, query, policyNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find claims by policy: %w", err)
	}

	return claims, nil
}

// List retrieves a paginated list of claims.
func (r *ClaimRepository) List(ctx context.Context, offset, limit int) ([]*models.Claim, error) {
	query := `
		SELECT id, claim_num, policy_num, claim_date, paid, value, cause, observations,
			created_at, updated_at
		FROM claims
		ORDER BY claim_num
		LIMIT $1 OFFSET $2
	`

	var claims []*models.Claim
	err := r.db.SelectContext(ctx, &claims, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list claims: %w", err)
	}

	return claims, nil
}

// Count returns the total number of claims.
func (r *ClaimRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM claims`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count claims: %w", err)
	}

	return count, nil
}

// CountByPolicy returns the number of claims for a specific policy.
func (r *ClaimRepository) CountByPolicy(ctx context.Context, policyNum string) (int64, error) {
	query := `SELECT COUNT(*) FROM claims WHERE policy_num = $1`

	var count int64
	err := r.db.GetContext(ctx, &count, query, policyNum)
	if err != nil {
		return 0, fmt.Errorf("failed to count claims by policy: %w", err)
	}

	return count, nil
}

// TotalPaidByPolicy returns the total paid amount for claims on a policy.
func (r *ClaimRepository) TotalPaidByPolicy(ctx context.Context, policyNum string) (float64, error) {
	query := `SELECT COALESCE(SUM(paid), 0) FROM claims WHERE policy_num = $1`

	var total float64
	err := r.db.GetContext(ctx, &total, query, policyNum)
	if err != nil {
		return 0, fmt.Errorf("failed to sum paid by policy: %w", err)
	}

	return total, nil
}

// TotalValueByPolicy returns the total claim value for a policy.
func (r *ClaimRepository) TotalValueByPolicy(ctx context.Context, policyNum string) (float64, error) {
	query := `SELECT COALESCE(SUM(value), 0) FROM claims WHERE policy_num = $1`

	var total float64
	err := r.db.GetContext(ctx, &total, query, policyNum)
	if err != nil {
		return 0, fmt.Errorf("failed to sum value by policy: %w", err)
	}

	return total, nil
}

// Exists checks if a claim exists by its claim number.
func (r *ClaimRepository) Exists(ctx context.Context, claimNum string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM claims WHERE claim_num = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, claimNum)
	if err != nil {
		return false, fmt.Errorf("failed to check claim existence: %w", err)
	}

	return exists, nil
}
