package handler

import (
	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// ServiceHandler handles service operations
type ServiceHandler struct {
	serviceUsecase *usecase.ServiceUsecase
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(serviceUsecase *usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{serviceUsecase}
}

// CreateService handles POST /services
func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)

	var req usecase.CreateServiceRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	req.BusinessID = businessID

	result, err := h.serviceUsecase.CreateService(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(201).JSON(result)
}

// GetService handles GET /services/:id
func (h *ServiceHandler) GetService(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)
	serviceID := c.Params("id")

	if serviceID == "" {
		appErr := apperror.NewBadRequest("invalid service ID")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.serviceUsecase.GetService(c.Context(), serviceID, businessID)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}

// ListServices handles GET /services
func (h *ServiceHandler) ListServices(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)

	result, err := h.serviceUsecase.ListServices(c.Context(), businessID)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}

// GetServiceStats handles GET /services/:id/stats
func (h *ServiceHandler) GetServiceStats(c *fiber.Ctx) error {
	businessID := c.Locals("user_id").(string)
	serviceID := c.Params("id")

	if serviceID == "" {
		appErr := apperror.NewBadRequest("invalid service ID")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.serviceUsecase.GetServiceStats(c.Context(), serviceID, businessID)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}
