package usecase

import (
	"context"

	"CLOAKBE/internal/apperror"
	"CLOAKBE/internal/domain"

	"github.com/google/uuid"
)

// ServiceUsecase handles service operations
type ServiceUsecase struct {
	serviceRepo  domain.ServiceRepository
	slotRepo     domain.SlotRepository
	businessRepo domain.BusinessRepository
}

// NewServiceUsecase creates a new service usecase
func NewServiceUsecase(
	serviceRepo domain.ServiceRepository,
	slotRepo domain.SlotRepository,
	businessRepo domain.BusinessRepository,
) *ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo:  serviceRepo,
		slotRepo:     slotRepo,
		businessRepo: businessRepo,
	}
}

// Request/Response types
type CreateServiceRequest struct {
	Name       string `json:"name"`
	TotalSlots int    `json:"total_slots"`
	BusinessID string `json:"-"`
}

type ServiceResponse struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	TotalSlots int    `json:"total_slots"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type ServiceStatsResponse struct {
	ServiceID  string `json:"service_id"`
	Name       string `json:"name"`
	TotalSlots int    `json:"total_slots"`
	Occupied   int    `json:"occupied"`
	Free       int    `json:"free"`
}

// CreateService creates a new service and generates slots
func (u *ServiceUsecase) CreateService(ctx context.Context, req CreateServiceRequest) (*ServiceResponse, error) {
	if req.Name == "" || req.TotalSlots <= 0 {
		return nil, apperror.NewValidationError("name and total_slots (>0) are required")
	}

	// Verify business exists
	_, err := u.businessRepo.FindByID(ctx, req.BusinessID)
	if err != nil {
		return nil, err
	}

	now := domain.NowTimestamp()

	// Create service
	service := &domain.Service{
		ID:         uuid.New().String(),
		BusinessID: req.BusinessID,
		Name:       req.Name,
		TotalSlots: req.TotalSlots,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := u.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	// Generate slots
	slots := make([]domain.Slot, req.TotalSlots)
	for i := 1; i <= req.TotalSlots; i++ {
		slots[i-1] = domain.Slot{
			ID:         uuid.New().String(),
			ServiceID:  service.ID,
			SlotNumber: i,
			Status:     domain.SlotStatusFree,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	if err := u.slotRepo.CreateBatch(ctx, slots); err != nil {
		return nil, err
	}

	return &ServiceResponse{
		ID:         service.ID,
		BusinessID: service.BusinessID,
		Name:       service.Name,
		TotalSlots: service.TotalSlots,
		CreatedAt:  service.CreatedAt,
		UpdatedAt:  service.UpdatedAt,
	}, nil
}

// GetService retrieves a service by ID
func (u *ServiceUsecase) GetService(ctx context.Context, serviceID, businessID string) (*ServiceResponse, error) {
	service, err := u.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if service.BusinessID != businessID {
		return nil, apperror.NewForbidden("service does not belong to this business")
	}

	return &ServiceResponse{
		ID:         service.ID,
		BusinessID: service.BusinessID,
		Name:       service.Name,
		TotalSlots: service.TotalSlots,
		CreatedAt:  service.CreatedAt,
		UpdatedAt:  service.UpdatedAt,
	}, nil
}

// ListServices lists all services for a business
func (u *ServiceUsecase) ListServices(ctx context.Context, businessID string) ([]ServiceResponse, error) {
	services, err := u.serviceRepo.ListByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	responses := make([]ServiceResponse, len(services))
	for i, service := range services {
		responses[i] = ServiceResponse{
			ID:         service.ID,
			BusinessID: service.BusinessID,
			Name:       service.Name,
			TotalSlots: service.TotalSlots,
			CreatedAt:  service.CreatedAt,
			UpdatedAt:  service.UpdatedAt,
		}
	}

	return responses, nil
}

// GetServiceStats returns occupancy statistics for a service
func (u *ServiceUsecase) GetServiceStats(ctx context.Context, serviceID, businessID string) (*ServiceStatsResponse, error) {
	service, err := u.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if service.BusinessID != businessID {
		return nil, apperror.NewForbidden("service does not belong to this business")
	}

	total, occupied, err := u.slotRepo.CountSlotsByStatus(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	return &ServiceStatsResponse{
		ServiceID:  service.ID,
		Name:       service.Name,
		TotalSlots: service.TotalSlots,
		Occupied:   occupied,
		Free:       total - occupied,
	}, nil
}
