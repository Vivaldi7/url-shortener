package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/models"
	"github.com/vivaldi7/golang_code/subscription_service/innternal/repository"
	"github.com/vivaldi7/golang_code/subscription_service/logger"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {

	//Parse UserID который в структуре Subscription идет с типом uuid, а в CreateSubscriptionRequest он тип string
	userID, err := uuid.Parse("UserID")
	if err != nil {
		logger.Log.Error("Invalid user_id format", "user_id", req.UserID)
		return nil, fmt.Errorf("Invalid user_id format %w", err)
	}

	//Parse StartDate который в структуре Subscription идет с типом time.Time, а в CreateSubscriptionRequest он тип string
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		logger.Log.Error("Invalid stert_date format", "start_date", req.StartDate)
		return nil, fmt.Errorf("Invalid stert_date format %w", err)
	}

	//Parse EndDate который в структуре Subscription идет с типом time.Time (если она указана),
	// а в CreateSubscriptionRequest он тип string
	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			logger.Log.Error("Invalid erd_date format", "end_date", req.EndDate)
			return nil, fmt.Errorf("Invalid end_date format %w", err)
		}
		endDate = &parsed
	}

	subscriptoin := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	_, err = s.repo.Create(ctx, subscriptoin)
	if err != nil {
		return nil, err
	}

	return subscriptoin, nil
}

func (s *SubscriptionService) Get(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *SubscriptionService) List(ctx context.Context, page, pageSize int64) ([]models.Subscription, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	return s.repo.List(ctx, pageSize, offset)
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) error {

	err := s.repo.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("subscription not found")
	}

	return err
}

func (s *SubscriptionService) Update(ctx context.Context, id int, req *models.UpdateSubscriptionRequest) error {

	updates := make(map[string]interface{})

	if req.ServiceName != nil {
		updates["service_name"] = req.ServiceName
	}

	if req.Price != nil {
		updates["price"] = req.Price
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("01-2006", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format: %w", err)
		}

		updates["start_date"] = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format: %w", err)
		}

		updates["end_date"] = endDate
	}

	return s.repo.Update(ctx, id, updates)

}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, req *models.TotalCostRequest) (int, error) {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date format: %w", err)
	}

	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		return 0, fmt.Errorf("invalid and_date format: %w", err)
	}

	if err := validateServiceName(req.ServiceName); err != nil {
		return 0, err
	}

	return s.repo.GetTotalCost(ctx, userID, req.ServiceName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}

func validateServiceName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("service_name cannot be empty")
	}
	return nil
}
