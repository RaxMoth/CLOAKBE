package repository

import (
	"context"
	"errors"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresCustomerRepository struct {
	db *database.Pool
}

func NewPostgresCustomerRepository(db *database.Pool) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (id, email, phone, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, customer.ID, customer.Email, customer.Phone, customer.CreatedAt)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"customers_email_key\" (SQLSTATE 23505)" {
			return apperror.NewConflict("email already exists")
		}
		return apperror.NewDatabaseError("failed to create customer", err)
	}

	return nil
}

func (r *PostgresCustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
		SELECT id, email, phone, created_at
		FROM customers
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	customer := &domain.Customer{}

	err := row.Scan(&customer.ID, &customer.Email, &customer.Phone, &customer.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("customer")
		}
		return nil, apperror.NewDatabaseError("failed to find customer", err)
	}

	return customer, nil
}

func (r *PostgresCustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	query := `
		SELECT id, email, phone, created_at
		FROM customers
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)
	customer := &domain.Customer{}

	err := row.Scan(&customer.ID, &customer.Email, &customer.Phone, &customer.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("customer")
		}
		return nil, apperror.NewDatabaseError("failed to find customer", err)
	}

	return customer, nil
}

// FindOrCreate returns existing customer or creates new one (upsert pattern)
func (r *PostgresCustomerRepository) FindOrCreate(ctx context.Context, email, phone string) (*domain.Customer, error) {
	// Try to find existing
	customer, err := r.FindByEmail(ctx, email)
	if err == nil {
		return customer, nil
	}

	// If not found, create
	if apperror.IsNotFound(err) {
		newCustomer := &domain.Customer{
			ID:        uuid.New().String(),
			Email:     email,
			Phone:     phone,
			CreatedAt: domain.NowTimestamp(),
		}

		if err := r.Create(ctx, newCustomer); err != nil {
			return nil, err
		}

		return newCustomer, nil
	}

	return nil, err
}
