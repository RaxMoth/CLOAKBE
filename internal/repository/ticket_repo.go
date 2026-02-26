package repository

import (
	"context"

	"CLOAKBE/internal/database"
	"CLOAKBE/internal/domain"
)

// PostgresTicketRepository implements TicketRepository for PostgreSQL
type PostgresTicketRepository struct {
	db *database.Pool
}

// NewPostgresTicketRepository creates a new ticket repository
func NewPostgresTicketRepository(db *database.Pool) *PostgresTicketRepository {
	return &PostgresTicketRepository{db}
}

// Create inserts a new ticket
func (r *PostgresTicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err := r.db.Exec(ctx, query,
		ticket.ID,
		ticket.ServiceID,
		ticket.SlotID,
		ticket.SlotNumber,
		ticket.CustomerID,
		ticket.Status,
		ticket.HMACDigest,
		ticket.IssuedAt,
	)
	return err
}

// FindByID retrieves a ticket by ID
func (r *PostgresTicketRepository) FindByID(ctx context.Context, id string) (*domain.Ticket, error) {
	query := `
		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
		FROM tickets WHERE id = $1
	`
	t := &domain.Ticket{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// FindByHMAC finds a ticket by HMAC digest
func (r *PostgresTicketRepository) FindByHMAC(ctx context.Context, hmacDigest string) (*domain.Ticket, error) {
	query := `
		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
		FROM tickets WHERE hmac_digest = $1
	`
	t := &domain.Ticket{}
	err := r.db.QueryRow(ctx, query, hmacDigest).Scan(
		&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ListByCustomerID lists tickets for a customer
func (r *PostgresTicketRepository) ListByCustomerID(ctx context.Context, customerID string) ([]domain.Ticket, error) {
	query := `
		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
		FROM tickets WHERE customer_id = $1 ORDER BY issued_at DESC
	`
	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		if err := rows.Scan(
			&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

// ListActiveByServiceID lists active tickets for a service
func (r *PostgresTicketRepository) ListActiveByServiceID(ctx context.Context, serviceID string) ([]domain.Ticket, error) {
	query := `
		SELECT id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
		FROM tickets WHERE service_id = $1 AND status = 'active'
	`
	rows, err := r.db.Query(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		if err := rows.Scan(
			&t.ID, &t.ServiceID, &t.SlotID, &t.SlotNumber, &t.CustomerID, &t.Status, &t.HMACDigest, &t.IssuedAt, &t.ReleasedAt, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

// UpdateStatus updates ticket status
func (r *PostgresTicketRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE tickets
		SET status = $2, released_at = CASE WHEN $2 = 'released' THEN NOW() ELSE released_at END, updated_at = NOW()
		WHERE id = $1
		RETURNING id, service_id, slot_id, slot_number, customer_id, status, hmac_digest, issued_at, released_at, created_at, updated_at
	`
	_, err := r.db.Exec(ctx, query, id, status)
	return err
}
