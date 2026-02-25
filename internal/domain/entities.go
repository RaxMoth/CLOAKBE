package domain

import (
	"context"
	"time"
)

// Slot Status Constants
const (
	SlotStatusFree     = "free"
	SlotStatusOccupied = "occupied"
)

// Ticket Status Constants
const (
	TicketStatusActive   = "active"
	TicketStatusReleased = "released"
)

// NowTimestamp returns current time as Unix timestamp
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// Business represents a business entity (e.g., a nightclub, event venue)
type Business struct {
	ID        string
	Name      string
	Email     string
	Password  string // bcrypt hash
	Role      string // "business"
	HMACKey   string // Secret key for QR signing
	CreatedAt int64
	UpdatedAt int64
}

// Customer represents a customer entity
type Customer struct {
	ID        string
	Email     string
	Phone     string
	CreatedAt int64
}

// Service represents a ticketing service (e.g., VIP table reservation, door entry)
type Service struct {
	ID         string
	BusinessID string
	Name       string
	TotalSlots int
	CreatedAt  int64
	UpdatedAt  int64
}

// Slot represents a single slot/ticket in a service
type Slot struct {
	ID         string
	ServiceID  string
	SlotNumber int
	Status     string // "free" or "occupied"
	CreatedAt  int64
	UpdatedAt  int64
}

// Ticket represents an issued ticket
type Ticket struct {
	ID         string
	ServiceID  string
	SlotID     string
	SlotNumber int
	CustomerID string // nullable for anonymous tickets
	Status     string // "active" or "released"
	HMACDigest string // Store the HMAC for audit trail
	IssuedAt   int64  // Unix timestamp when ticket was created
	ReleasedAt int64  // Unix timestamp when ticket was released (nullable)
	CreatedAt  int64
	UpdatedAt  int64
}

// Repository Interfaces

// BusinessRepository defines business persistence operations
type BusinessRepository interface {
	Create(ctx context.Context, business *Business) error
	FindByID(ctx context.Context, id string) (*Business, error)
	FindByEmail(ctx context.Context, email string) (*Business, error)
	Update(ctx context.Context, business *Business) error
}

// CustomerRepository defines customer persistence operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindByEmail(ctx context.Context, email string) (*Customer, error)
	FindOrCreate(ctx context.Context, email, phone string) (*Customer, error)
}

// ServiceRepository defines service persistence operations
type ServiceRepository interface {
	Create(ctx context.Context, service *Service) error
	FindByID(ctx context.Context, id string) (*Service, error)
	ListByBusinessID(ctx context.Context, businessID string) ([]Service, error)
	Update(ctx context.Context, service *Service) error
	Delete(ctx context.Context, id string) error
}

// SlotRepository defines slot persistence operations
type SlotRepository interface {
	Create(ctx context.Context, slot *Slot) error
	CreateBatch(ctx context.Context, slots []Slot) error
	FindByID(ctx context.Context, id string) (*Slot, error)
	ListByServiceID(ctx context.Context, serviceID string) ([]Slot, error)
	ClaimNextFreeSlot(ctx context.Context, serviceID string) (*Slot, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	CountSlotsByStatus(ctx context.Context, serviceID string) (total, occupied int, err error)
}

// TicketRepository defines ticket persistence operations
type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	FindByID(ctx context.Context, id string) (*Ticket, error)
	FindByHMAC(ctx context.Context, hmacDigest string) (*Ticket, error)
	ListByCustomerID(ctx context.Context, customerID string) ([]Ticket, error)
	ListActiveByServiceID(ctx context.Context, serviceID string) ([]Ticket, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
