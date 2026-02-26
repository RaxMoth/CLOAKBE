package middleware

import (
	"strings"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			appErr := apperror.NewUnauthorized("Authorization header is required")
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := apperror.NewUnauthorized("Invalid authorization header format")
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &usecase.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			appErr := apperror.NewUnauthorized("Invalid or expired token")
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
		}

		// Extract claims and set in context
		if claims, ok := token.Claims.(*usecase.CustomClaims); ok {
			c.Locals("user_id", claims.UserID)
			c.Locals("email", claims.Email)
			c.Locals("role", claims.Role)
		} else {
			appErr := apperror.NewUnauthorized("Invalid token claims")
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
		}

		return c.Next()
	}
}

// RoleMiddleware enforces role-based access control
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		appErr := apperror.NewForbidden("insufficient permissions")
		return c.Status(appErr.StatusCode).JSON(fiber.Map{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
	}
}

// GetUserIDFromContext is a helper function to get user ID from Fiber context
func GetUserIDFromContext(c *fiber.Ctx) string {
	userID := c.Locals("user_id")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// GetEmailFromContext is a helper function to get email from Fiber context
func GetEmailFromContext(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetRoleFromContext is a helper function to get role from Fiber context
func GetRoleFromContext(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}
