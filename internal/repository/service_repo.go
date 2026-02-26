package repository

import (
	"context"
	"errors"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PostgresServiceRepository struct {
	db *database.Pool
}

func NewPostgresServiceRepository(db *database.Pool) *PostgresServiceRepository {
	return &PostgresServiceRepository{db: db}
}

func (r *PostgresServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	query := `
		INSERT INTO services (id, business_id, name, total_slots, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		service.ID,
		service.BusinessID,
		service.Name,
		service.TotalSlots,
		service.CreatedAt,
		service.UpdatedAt,
	)
	if err != nil {
		return apperror.NewDatabaseError("failed to create service", err)
	}

	return nil
}

func (r *PostgresServiceRepository) FindByID(ctx context.Context, id string) (*domain.Service, error) {
	query := `
		SELECT id, business_id, name, total_slots, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	service := &domain.Service{}

	err := row.Scan(&service.ID, &service.BusinessID, &service.Name, &service.TotalSlots, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("service")
		}
		return nil, apperror.NewDatabaseError("failed to find service", err)
	}

	return service, nil
}

func (r *PostgresServiceRepository) ListByBusinessID(ctx context.Context, businessID string) ([]domain.Service, error) {
	query := `
		SELECT id, business_id, name, total_slots, created_at, updated_at
		FROM services
		WHERE business_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, apperror.NewDatabaseError("failed to list services", err)
	}
	defer rows.Close()

	services := []domain.Service{}
	for rows.Next() {
		service := domain.Service{}
		if err := rows.Scan(&service.ID, &service.BusinessID, &service.Name, &service.TotalSlots, &service.CreatedAt, &service.UpdatedAt); err != nil {
			return nil, apperror.NewDatabaseError("failed to scan service", err)
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, apperror.NewDatabaseError("failed to iterate services", err)
	}

	return services, nil
}

func (r *PostgresServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	query := `
		UPDATE services
		SET name = $1, total_slots = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(ctx, query, service.Name, service.TotalSlots, service.UpdatedAt, service.ID)
	if err != nil {
		return apperror.NewDatabaseError("failed to update service", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("service not found")
	}

	return nil
}

func (r *PostgresServiceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewDatabaseError("failed to delete service", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("service not found")
	}

	return nil
}
