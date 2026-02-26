package handler

import (
	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// TicketHandler handles ticket operations
type TicketHandler struct {
	ticketUsecase *usecase.TicketUsecase
}

// NewTicketHandler creates a new ticket handler
func NewTicketHandler(ticketUsecase *usecase.TicketUsecase) *TicketHandler {
	return &TicketHandler{ticketUsecase}
}

// CheckIn handles POST /tickets/checkin - Business creates QR for customer
func (h *TicketHandler) CheckIn(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)

	var req usecase.CheckInRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	req.BusinessID = businessID

	result, err := h.ticketUsecase.CheckIn(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}

// Scan handles POST /tickets/scan - Business scans QR code
func (h *TicketHandler) Scan(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)

	var req usecase.ScanRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	req.BusinessID = businessID

	result, err := h.ticketUsecase.Scan(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}

// Release handles POST /tickets/:id/release - Business releases ticket
func (h *TicketHandler) Release(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)
	ticketID := c.Params("id")

	if ticketID == "" {
		appErr := apperror.NewBadRequest("invalid ticket ID")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	if err := h.ticketUsecase.Release(c.Context(), ticketID, businessID); err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

// GetTicket handles GET /tickets/:id - Customer views their ticket
func (h *TicketHandler) GetTicket(c *fiber.Ctx) error {
	_ = c.Locals("user_id").(string) // customerID - will be used when implementing actual logic
	ticketID := c.Params("id")

	if ticketID == "" {
		appErr := apperror.NewBadRequest("invalid ticket ID")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	// TODO: Implement GetTicket in usecase to fetch and verify customer owns ticket
	// For now returning not found
	err := apperror.NewNotFound("ticket endpoint not yet implemented")
	return c.Status(err.StatusCode).JSON(errorResponse(err))
}

func errorResponse(err *apperror.AppError) fiber.Map {
	return fiber.Map{
		"code":    err.Code,
		"message": err.Message,
		"details": err.Details,
	}
}
