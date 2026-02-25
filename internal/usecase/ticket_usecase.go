package usecase

import (
	"context"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/domain"
	"CLOAKBE/internal/qr"

	"github.com/google/uuid"
)

// TicketUsecase handles ticket operations
type TicketUsecase struct {
	ticketRepo   domain.TicketRepository
	slotRepo     domain.SlotRepository
	serviceRepo  domain.ServiceRepository
	businessRepo domain.BusinessRepository
}

// NewTicketUsecase creates a new ticket usecase
func NewTicketUsecase(
	ticketRepo domain.TicketRepository,
	slotRepo domain.SlotRepository,
	serviceRepo domain.ServiceRepository,
	businessRepo domain.BusinessRepository,
) *TicketUsecase {
	return &TicketUsecase{
		ticketRepo:   ticketRepo,
		slotRepo:     slotRepo,
		serviceRepo:  serviceRepo,
		businessRepo: businessRepo,
	}
}

// Request/Response types
type CheckInRequest struct {
	ServiceID  string `json:"service_id"`
	BusinessID string `json:"-"`
}

type CheckInResponse struct {
	TicketID   string `json:"ticket_id"`
	SlotNumber int    `json:"slot_number"`
	QRPayload  string `json:"qr_payload"` // base64 encoded
	IssuedAt   int64  `json:"issued_at"`
}

type ScanRequest struct {
	QRPayload  string `json:"qr_payload"` // base64 encoded payload
	BusinessID string `json:"-"`
}

type ScanResponse struct {
	SlotNumber int    `json:"slot_number"`
	ServiceID  string `json:"service_id"`
	Status     string `json:"status"`
	ReleasedAt *int64 `json:"released_at"`
}

// CheckIn claims a slot and creates a QR code ticket
func (u *TicketUsecase) CheckIn(ctx context.Context, req CheckInRequest) (*CheckInResponse, error) {
	// Verify service ownership
	service, err := u.serviceRepo.FindByID(ctx, req.ServiceID)
	if err != nil {
		return nil, err
	}

	if service.BusinessID != req.BusinessID {
		return nil, apperror.NewForbidden("service does not belong to this business")
	}

	// Get business to retrieve HMAC key
	business, err := u.businessRepo.FindByID(ctx, req.BusinessID)
	if err != nil {
		return nil, err
	}

	// Claim next free slot (with row locking to prevent race conditions)
	slot, err := u.slotRepo.ClaimNextFreeSlot(ctx, req.ServiceID)
	if err != nil {
		return nil, err
	}

	// Create QR payload
	payload := qr.New(uuid.New().String(), req.ServiceID, req.BusinessID, slot.SlotNumber)

	// Sign payload with business HMAC key
	if err := payload.Sign(business.HMACKey); err != nil {
		return nil, apperror.NewInternalServer("QR signing failed")
	}

	// Encode payload to base64
	encoded, err := payload.Encode()
	if err != nil {
		return nil, apperror.NewInternalServer("QR encoding failed")
	}

	// Create ticket record
	ticket := &domain.Ticket{
		ID:         payload.TicketID,
		ServiceID:  req.ServiceID,
		SlotID:     slot.ID,
		SlotNumber: slot.SlotNumber,
		Status:     domain.TicketStatusActive,
		HMACDigest: payload.HMAC,
		IssuedAt:   payload.IssuedAt,
		CreatedAt:  domain.NowTimestamp(),
		UpdatedAt:  domain.NowTimestamp(),
	}

	if err := u.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	return &CheckInResponse{
		TicketID:   ticket.ID,
		SlotNumber: slot.SlotNumber,
		QRPayload:  encoded,
		IssuedAt:   payload.IssuedAt,
	}, nil
}

// Scan verifies a QR code and returns ticket status
func (u *TicketUsecase) Scan(ctx context.Context, req ScanRequest) (*ScanResponse, error) {
	// Decode QR payload
	payload, err := qr.Decode(req.QRPayload)
	if err != nil {
		return nil, apperror.NewBadRequest("invalid QR payload")
	}

	// Verify business ownership
	if payload.BusinessID != req.BusinessID {
		return nil, apperror.NewForbidden("QR code does not belong to this business")
	}

	// Get business HMAC key
	business, err := u.businessRepo.FindByID(ctx, req.BusinessID)
	if err != nil {
		return nil, err
	}

	// Verify HMAC signature
	if err := payload.Verify(business.HMACKey); err != nil {
		return nil, apperror.NewBadRequest("invalid QR signature")
	}

	// Find ticket by HMAC for audit trail
	ticket, err := u.ticketRepo.FindByHMAC(ctx, payload.HMAC)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, apperror.NewBadRequest("ticket not found")
		}
		return nil, err
	}

	return &ScanResponse{
		SlotNumber: ticket.SlotNumber,
		ServiceID:  ticket.ServiceID,
		Status:     ticket.Status,
		ReleasedAt: (*int64)(nil), // Will populate if released
	}, nil
}

// Release frees a slot and marks ticket as released
func (u *TicketUsecase) Release(ctx context.Context, ticketID, businessID string) error {
	// Find ticket
	ticket, err := u.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return err
	}

	// Verify service ownership
	service, err := u.serviceRepo.FindByID(ctx, ticket.ServiceID)
	if err != nil {
		return err
	}

	if service.BusinessID != businessID {
		return apperror.NewForbidden("ticket does not belong to this business")
	}

	// Mark ticket as released
	now := domain.NowTimestamp()
	if err := u.ticketRepo.UpdateStatus(ctx, ticketID, domain.TicketStatusReleased); err != nil {
		return err
	}

	// Free the slot
	if ticket.SlotID != "" {
		if err := u.slotRepo.UpdateStatus(ctx, ticket.SlotID, domain.SlotStatusFree); err != nil {
			return err
		}
	}

	_ = now
	return nil
}
