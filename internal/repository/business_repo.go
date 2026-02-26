package repository

import (
	"context"

	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"
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
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	err := r.db.QueryRow(ctx, query,
		b.ID, b.Name, b.Email, b.Password, b.Role, b.HMACKey,
	).Scan()

	return err
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
		return nil, err
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
		return nil, err
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

	err := r.db.QueryRow(ctx, query,
		b.ID, b.Name, b.Email, b.Password, b.HMACKey,
	).Scan()

	return err
}
