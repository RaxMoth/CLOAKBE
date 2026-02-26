package repository

import (
	"context"
	"errors"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PostgresSlotRepository struct {
	db *database.Pool
}

func NewPostgresSlotRepository(db *database.Pool) *PostgresSlotRepository {
	return &PostgresSlotRepository{db: db}
}

func (r *PostgresSlotRepository) Create(ctx context.Context, slot *domain.Slot) error {
	query := `
		INSERT INTO slots (id, service_id, slot_number, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		slot.ID,
		slot.ServiceID,
		slot.SlotNumber,
		slot.Status,
		slot.CreatedAt,
		slot.UpdatedAt,
	)
	if err != nil {
		return apperror.NewDatabaseError("failed to create slot", err)
	}

	return nil
}

// CreateBatch inserts multiple slots in one operation (used for service creation)
func (r *PostgresSlotRepository) CreateBatch(ctx context.Context, slots []domain.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO slots (id, service_id, slot_number, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, slot := range slots {
		batch.Queue(query, slot.ID, slot.ServiceID, slot.SlotNumber, slot.Status, slot.CreatedAt, slot.UpdatedAt)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(slots); i++ {
		_, err := results.Exec()
		if err != nil {
			return apperror.NewDatabaseError("failed to create slot in batch", err)
		}
	}

	return nil
}

func (r *PostgresSlotRepository) FindByID(ctx context.Context, id string) (*domain.Slot, error) {
	query := `
		SELECT id, service_id, slot_number, status, created_at, updated_at
		FROM slots
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	slot := &domain.Slot{}

	err := row.Scan(&slot.ID, &slot.ServiceID, &slot.SlotNumber, &slot.Status, &slot.CreatedAt, &slot.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("slot not found")
		}
		return nil, apperror.NewDatabaseError("failed to find slot", err)
	}

	return slot, nil
}

// ListByServiceID retrieves all slots for a service
func (r *PostgresSlotRepository) ListByServiceID(ctx context.Context, serviceID string) ([]domain.Slot, error) {
	query := `
		SELECT id, service_id, slot_number, status, created_at, updated_at
		FROM slots
		WHERE service_id = $1
		ORDER BY slot_number ASC
	`

	rows, err := r.db.Query(ctx, query, serviceID)
	if err != nil {
		return nil, apperror.NewDatabaseError("failed to list slots", err)
	}
	defer rows.Close()

	slots := []domain.Slot{}
	for rows.Next() {
		slot := domain.Slot{}
		if err := rows.Scan(&slot.ID, &slot.ServiceID, &slot.SlotNumber, &slot.Status, &slot.CreatedAt, &slot.UpdatedAt); err != nil {
			return nil, apperror.NewDatabaseError("failed to scan slot", err)
		}
		slots = append(slots, slot)
	}

	if err = rows.Err(); err != nil {
		return nil, apperror.NewDatabaseError("failed to iterate slots", err)
	}

	return slots, nil
}

// ClaimNextFreeSlot claims the next available free slot using row-level locking
// This prevents race conditions when multiple check-ins occur simultaneously
// Uses: SELECT ... FOR UPDATE SKIP LOCKED to prevent deadlocks and allow concurrent operations
func (r *PostgresSlotRepository) ClaimNextFreeSlot(ctx context.Context, serviceID string) (*domain.Slot, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperror.NewDatabaseError("failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Find and lock the next free slot (SKIP LOCKED prevents waiting on locked rows)
	selectQuery := `
		SELECT id, service_id, slot_number, status, created_at, updated_at
		FROM slots
		WHERE service_id = $1 AND status = $2
		ORDER BY slot_number ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	row := tx.QueryRow(ctx, selectQuery, serviceID, domain.SlotStatusFree)
	slot := &domain.Slot{}

	err = row.Scan(&slot.ID, &slot.ServiceID, &slot.SlotNumber, &slot.Status, &slot.CreatedAt, &slot.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewConflict("no free slots available")
		}
		return nil, apperror.NewDatabaseError("failed to find free slot", err)
	}

	// Mark slot as occupied within same transaction
	updateQuery := `
		UPDATE slots
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, service_id, slot_number, status, created_at, updated_at
	`

	now := domain.NowTimestamp()
	updateRow := tx.QueryRow(ctx, updateQuery, domain.SlotStatusOccupied, now, slot.ID)

	err = updateRow.Scan(&slot.ID, &slot.ServiceID, &slot.SlotNumber, &slot.Status, &slot.CreatedAt, &slot.UpdatedAt)
	if err != nil {
		return nil, apperror.NewDatabaseError("failed to update slot status", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperror.NewDatabaseError("failed to commit transaction", err)
	}

	return slot, nil
}

// UpdateStatus updates a slot's status (used for releasing tickets)
func (r *PostgresSlotRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE slots
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, status, domain.NowTimestamp(), id)
	if err != nil {
		return apperror.NewDatabaseError("failed to update slot", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("slot not found")
	}

	return nil
}

// CountSlotsByStatus returns counts of slots in each status for a service
func (r *PostgresSlotRepository) CountSlotsByStatus(ctx context.Context, serviceID string) (total, occupied int, err error) {
	query := `
		SELECT COUNT(*) total,
		       COUNT(CASE WHEN status = $1 THEN 1 END) occupied
		FROM slots
		WHERE service_id = $2
	`

	row := r.db.QueryRow(ctx, query, domain.SlotStatusOccupied, serviceID)
	if err := row.Scan(&total, &occupied); err != nil {
		return 0, 0, apperror.NewDatabaseError("failed to count slots", err)
	}

	return total, occupied, nil
}
