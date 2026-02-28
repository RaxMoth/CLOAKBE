package repository

import (
	"context"
	"errors"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"

	"github.com/jackc/pgx/v5"
)

// PostgresBusinessRepository implements BusinessRepository for PostgreSQL
type PostgresBusinessRepository struct {
	db *database.Pool
}

// NewPostgresBusinessRepository creates a new business repository
func NewPostgresBusinessRepository(db *database.Pool) *PostgresBusinessRepository {
	return &PostgresBusinessRepository{db}
}

// Create creates a new business
func (r *PostgresBusinessRepository) Create(ctx context.Context, b *domain.Business) error {
	query := `
		INSERT INTO businesses (id, name, email, password, role, hmac_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		b.ID, b.Name, b.Email, b.Password, b.Role, b.HMACKey, b.CreatedAt, b.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"businesses_email_key\" (SQLSTATE 23505)" {
			return apperror.NewConflict("email already registered")
		}
		return apperror.NewDatabaseError("failed to create business", err)
	}

	return nil
}

// FindByID finds a business by ID
func (r *PostgresBusinessRepository) FindByID(ctx context.Context, id string) (*domain.Business, error) {
	query := `
		SELECT id, name, email, password, role, hmac_key, created_at, updated_at
		FROM businesses WHERE id = $1
	`

	b := &domain.Business{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Email, &b.Password, &b.Role, &b.HMACKey, &b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("business")
		}
		return nil, apperror.NewDatabaseError("failed to find business by id", err)
	}

	return b, nil
}

// FindByEmail finds a business by email
func (r *PostgresBusinessRepository) FindByEmail(ctx context.Context, email string) (*domain.Business, error) {
	query := `
		SELECT id, name, email, password, role, hmac_key, created_at, updated_at
		FROM businesses WHERE email = $1
	`

	b := &domain.Business{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&b.ID, &b.Name, &b.Email, &b.Password, &b.Role, &b.HMACKey, &b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("business")
		}
		return nil, apperror.NewDatabaseError("failed to find business by email", err)
	}

	return b, nil
}

// Update updates a business
func (r *PostgresBusinessRepository) Update(ctx context.Context, b *domain.Business) error {
	query := `
		UPDATE businesses
		SET name = $2, email = $3, password = $4, hmac_key = $5, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		b.ID, b.Name, b.Email, b.Password, b.HMACKey,
	)

	if err != nil {
		return apperror.NewDatabaseError("failed to update business", err)
	}

	return nil
}
