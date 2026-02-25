package qr

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Payload represents the QR code data that gets signed and encoded
type Payload struct {
	Version    int    `json:"v"`   // Version 1
	TicketID   string `json:"tid"` // Ticket UUID
	ServiceID  string `json:"sid"` // Service UUID
	BusinessID string `json:"bid"` // Business UUID
	SlotNumber int    `json:"slot"`
	IssuedAt   int64  `json:"iat"` // Unix timestamp
	HMAC       string `json:"hmac"`
}

// New creates a new QR payload with current timestamp
func New(ticketID, serviceID, businessID string, slotNumber int) *Payload {
	return &Payload{
		Version:    1,
		TicketID:   ticketID,
		ServiceID:  serviceID,
		BusinessID: businessID,
		SlotNumber: slotNumber,
		IssuedAt:   time.Now().Unix(),
	}
}

// canonicalString creates the string to sign over
func (p *Payload) canonicalString() string {
	return fmt.Sprintf("v=%d&tid=%s&sid=%s&bid=%s&slot=%d&iat=%d",
		p.Version,
		p.TicketID,
		p.ServiceID,
		p.BusinessID,
		p.SlotNumber,
		p.IssuedAt,
	)
}

// Sign computes HMAC over the payload
func (p *Payload) Sign(secret string) error {
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	canonical := p.canonicalString()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(canonical))
	p.HMAC = fmt.Sprintf("%x", h.Sum(nil))
	return nil
}

// Verify checks if the HMAC is valid
func (p *Payload) Verify(secret string) error {
	if p.HMAC == "" {
		return fmt.Errorf("payload not signed")
	}

	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	canonical := p.canonicalString()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(canonical))
	expectedHMAC := fmt.Sprintf("%x", h.Sum(nil))

	if expectedHMAC != p.HMAC {
		return fmt.Errorf("invalid hmac signature")
	}

	return nil
}

// Encode converts payload to base64 JSON string (for QR code)
func (p *Payload) Encode() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

// Decode decodes base64 JSON string back to Payload
func Decode(encoded string) (*Payload, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encoding: %w", err)
	}

	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("invalid payload JSON: %w", err)
	}

	// Regenerate ticket ID if it's not a valid UUID (for validation)
	if _, err := uuid.Parse(payload.TicketID); err != nil {
		return nil, fmt.Errorf("invalid ticket ID format: %w", err)
	}

	return &payload, nil
}
