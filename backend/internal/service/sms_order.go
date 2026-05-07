package service

import (
	"context"
	"fmt"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrSmsOrderNotFound = infraerrors.NotFound("SMS_ORDER_NOT_FOUND", "SMS order not found")
)

// SmsOrder is the domain model for SMS order query.
type SmsOrder struct {
	ID          int64
	OrderNo     string
	PhoneNumber string
	SmsContent  string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SmsOrderRepository defines the persistence interface for SMS orders.
type SmsOrderRepository interface {
	GetByOrderNo(ctx context.Context, orderNo string) (*SmsOrder, error)
}

// SmsOrderService provides SMS order business logic.
type SmsOrderService struct {
	repo SmsOrderRepository
}

// NewSmsOrderService creates a new SmsOrderService.
func NewSmsOrderService(repo SmsOrderRepository) *SmsOrderService {
	return &SmsOrderService{repo: repo}
}

// GetByOrderNo retrieves an SMS order by its order number.
func (s *SmsOrderService) GetByOrderNo(ctx context.Context, orderNo string) (*SmsOrder, error) {
	if orderNo == "" {
		return nil, fmt.Errorf("order number is required")
	}
	return s.repo.GetByOrderNo(ctx, orderNo)
}
