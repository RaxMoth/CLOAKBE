package usecase

import (
	"context"
	"time"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase handles authentication logic
type AuthUsecase struct {
	businessRepo domain.BusinessRepository
	customerRepo domain.CustomerRepository
	jwtSecret    string
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(
	businessRepo domain.BusinessRepository,
	customerRepo domain.CustomerRepository,
	jwtSecret string,
) *AuthUsecase {
	return &AuthUsecase{
		businessRepo: businessRepo,
		customerRepo: customerRepo,
		jwtSecret:    jwtSecret,
	}
}

// Request/Response types
type BusinessRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BusinessLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustomerLoginRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// JWT Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// BusinessRegister handles business registration
func (u *AuthUsecase) BusinessRegister(ctx context.Context, req BusinessRegisterRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, apperror.NewValidationError("email, password, and name are required", map[string]string{})
	}

	// Check if business already exists
	_, err := u.businessRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, apperror.NewConflict("email already registered")
	}
	if !apperror.IsNotFound(err) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewInternalServer("password hashing failed", err)
	}

	business := &domain.Business{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "business",
		HMACKey:   uuid.New().String(), // Secret key for QR signing
		CreatedAt: domain.NowTimestamp(),
		UpdatedAt: domain.NowTimestamp(),
	}

	if err := u.businessRepo.Create(ctx, business); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := u.generateToken(business.ID, business.Email, "business")
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:  token,
		UserID: business.ID,
		Role:   "business",
	}, nil
}

// BusinessLogin handles business login
func (u *AuthUsecase) BusinessLogin(ctx context.Context, req BusinessLoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, apperror.NewValidationError("email and password are required", map[string]string{})
	}

	// Find business by email
	business, err := u.businessRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if apperror.IsNotFound(err) {
			return nil, apperror.NewUnauthorized("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(business.Password), []byte(req.Password)); err != nil {
		return nil, apperror.NewUnauthorized("invalid credentials")
	}

	// Generate JWT
	token, err := u.generateToken(business.ID, business.Email, "business")
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:  token,
		UserID: business.ID,
		Role:   "business",
	}, nil
}

// CustomerLogin handles customer login (upsert pattern)
func (u *AuthUsecase) CustomerLogin(ctx context.Context, req CustomerLoginRequest) (*AuthResponse, error) {
	if req.Email == "" {
		return nil, apperror.NewValidationError("email is required", map[string]string{})
	}

	// Find or create customer
	customer, err := u.customerRepo.FindOrCreate(ctx, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := u.generateToken(customer.ID, customer.Email, "customer")
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:  token,
		UserID: customer.ID,
		Role:   "customer",
	}, nil
}

// generateToken creates a JWT token
func (u *AuthUsecase) generateToken(userID, email, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", apperror.NewInternalServer("token generation failed", err)
	}

	return tokenString, nil
}
