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

// CustomerCheckIn handles POST /tickets/checkin - Customer checks in to an event
func (h *TicketHandler) CustomerCheckIn(c *fiber.Ctx) error {
	_ = c.Locals("user_id").(string) // customerID - for future use (tracking which customer checked in)

	var req struct {
		ServiceID string `json:"service_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	if req.ServiceID == "" {
		appErr := apperror.NewBadRequest("service_id is required")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.ticketUsecase.CustomerCheckIn(c.Context(), req.ServiceID)
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

// GetCustomerTickets handles GET /customers/:id/tickets - Get all tickets for a customer
func (h *TicketHandler) GetCustomerTickets(c *fiber.Ctx) error {
	customerID := c.Params("id")
	if customerID == "" {
		appErr := apperror.NewBadRequest("invalid customer ID")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	tickets, err := h.ticketUsecase.GetCustomerTickets(c.Context(), customerID)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(fiber.Map{
		"tickets": tickets,
	})
}

func errorResponse(err *apperror.AppError) fiber.Map {
	return fiber.Map{
		"code":    err.Code,
		"message": err.Message,
		"details": err.Details,
	}
}
