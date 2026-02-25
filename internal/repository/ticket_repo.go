package repository
// Database repositories for CLOAK
// Note: These are simplified examples. In production, use sqlc for generated code.
// Source of truth: internal/domain/entities.go defines the interfaces

package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/RaxMoth/qrcheck-backend/internal/database"
	"github.com/RaxMoth/qrcheck-backend/internal/domain"
)

// PostgresTicketRepository implements TicketRepository for PostgreSQL
type PostgresTicketRepository struct {
	db *database.Pool
}

// NewPostgresTicketRepository creates a new ticket repository
func NewPostgresTicketRepository(db *database.Pool) *PostgresTicketRepository {
	return &PostgresTicketRepository{db}
}

// Create creates a new ticket
func (r *PostgresTicketRepository) Create(ctx context.Context, t *domain.Ticket) (*domain.Ticket, error) {
	query := `
		INSERT INTO tickets (id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, to_timestamp($8), NOW(), NOW())
		RETURNING id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		t.ID, t.ServiceID, t.SlotID, t.SlotNumber, t.CustomerID, t.Status, t.HMACDigest, t.IssuedAt,
	).Scan(
		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}




















































































































}	return t, nil	}		return nil, err	if err != nil {	)		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,	err := r.db.QueryRow(ctx, query, id, status).Scan(	t := &domain.Ticket{}	`		RETURNING id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at		WHERE id = $1		SET status = $2, released_at = CASE WHEN $2 = 'released' THEN NOW() ELSE released_at END, updated_at = NOW()		UPDATE tickets	query := `func (r *PostgresTicketRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) (*domain.Ticket, error) {// UpdateStatus updates ticket status}	return tickets, rows.Err()	}		tickets = append(tickets, t)		}			return nil, err		); err != nil {			&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,		if err := rows.Scan(		var t domain.Ticket	for rows.Next() {	var tickets []domain.Ticket	defer rows.Close()	}		return nil, err	if err != nil {	rows, err := r.db.Query(ctx, query, serviceID)	`		FROM tickets WHERE service_id = $1 AND status = 'active'		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at	query := `func (r *PostgresTicketRepository) ListActiveByServiceID(ctx context.Context, serviceID uuid.UUID) ([]domain.Ticket, error) {// ListActiveByServiceID lists active tickets for a service}	return tickets, rows.Err()	}		tickets = append(tickets, t)		}			return nil, err		); err != nil {			&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,		if err := rows.Scan(		var t domain.Ticket	for rows.Next() {	var tickets []domain.Ticket	defer rows.Close()	}		return nil, err	if err != nil {	rows, err := r.db.Query(ctx, query, customerID)	`		FROM tickets WHERE customer_id = $1 ORDER BY issued_at DESC		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at	query := `func (r *PostgresTicketRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Ticket, error) {// ListByCustomerID lists tickets for a customer}	return t, nil	}		return nil, err	if err != nil {	)		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,	err := r.db.QueryRow(ctx, query, hmac).Scan(	t := &domain.Ticket{}	`		FROM tickets WHERE hmac_digest = $1		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at	query := `func (r *PostgresTicketRepository) FindByHMAC(ctx context.Context, hmac string) (*domain.Ticket, error) {// FindByHMAC finds a ticket by HMAC digest}	return t, nil	}		return nil, err	if err != nil {	)		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,	err := r.db.QueryRow(ctx, query, id).Scan(	t := &domain.Ticket{}	`		FROM tickets WHERE id = $1		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at	query := `func (r *PostgresTicketRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {// FindByID finds a ticket by ID}	return t, nil