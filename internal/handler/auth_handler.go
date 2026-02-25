package handler

import (
	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication operations
type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase}
}

// BusinessRegister handles POST /auth/business/register
func (h *AuthHandler) BusinessRegister(c *fiber.Ctx) error {
	var req usecase.BusinessRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.authUsecase.BusinessRegister(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(201).JSON(result)
}

// BusinessLogin handles POST /auth/business/login
func (h *AuthHandler) BusinessLogin(c *fiber.Ctx) error {
	var req usecase.BusinessLoginRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.authUsecase.BusinessLogin(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}

// CustomerLogin handles POST /auth/customer/login
func (h *AuthHandler) CustomerLogin(c *fiber.Ctx) error {
	var req usecase.CustomerLoginRequest
	if err := c.BodyParser(&req); err != nil {
		appErr := apperror.NewBadRequest("invalid request body")
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	result, err := h.authUsecase.CustomerLogin(c.Context(), req)
	if err != nil {
		appErr := apperror.From(err)
		return c.Status(appErr.StatusCode).JSON(errorResponse(appErr))
	}

	return c.Status(200).JSON(result)
}
