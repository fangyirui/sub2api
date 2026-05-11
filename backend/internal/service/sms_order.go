package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrSmsOrderNotFound   = infraerrors.NotFound("SMS_ORDER_NOT_FOUND", "SMS order not found")
	ErrSmsOrderFetchFailed = infraerrors.ServiceUnavailable("SMS_ORDER_FETCH_FAILED", "Failed to fetch phone number, please contact support")
	ErrSmsOrderExpired    = infraerrors.ServiceUnavailable("SMS_ORDER_EXPIRED", "SMS order has expired")
)

const (
	SmsOrderStatusCreated  = "created"
	SmsOrderStatusPending  = "pending"
	SmsOrderStatusReceived = "received"
	SmsOrderStatusExpired  = "expired"
	SmsOrderStatusFailed   = "failed"
)

type SmsOrder struct {
	ID          int64
	OrderNo     string
	ServiceType string
	PhoneNumber string
	HeroSmsID   string
	SmsContent  string
	Status      string
	PendingAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type HeroSmsNumber struct {
	ID    string
	Phone string
}

type HeroSmsStatus struct {
	Received bool
	Code     string
	Text     string
}

type HeroSmsClient interface {
	GetNumber(ctx context.Context, serviceType string) (*HeroSmsNumber, error)
	GetNumberWithRetry(ctx context.Context, serviceType string) (*HeroSmsNumber, error)
	GetStatus(ctx context.Context, id string) (*HeroSmsStatus, error)
}

type SmsOrderListFilter struct {
	Page      int
	PageSize  int
	Status    string
	StartTime *time.Time
	EndTime   *time.Time
	SortBy    string
	SortDesc  bool
}

type SmsOrderListResult struct {
	Items []*SmsOrder
	Total int
}

type SmsOrderRepository interface {
	GetByOrderNo(ctx context.Context, orderNo string) (*SmsOrder, error)
	Create(ctx context.Context, order *SmsOrder) (*SmsOrder, error)
	List(ctx context.Context, filter SmsOrderListFilter) (*SmsOrderListResult, error)
	ListPending(ctx context.Context) ([]*SmsOrder, error)
	AssignNumber(ctx context.Context, id int64, phone, heroSmsID string, pendingAt time.Time) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateSmsContent(ctx context.Context, id int64, content string, status string) error
}

type SmsOrderService struct {
	repo       SmsOrderRepository
	heroClient HeroSmsClient
	cfg        *config.Config
	stopOnce   sync.Once
	stopCh     chan struct{}
}

func NewSmsOrderService(repo SmsOrderRepository, heroClient HeroSmsClient, cfg *config.Config) *SmsOrderService {
	return &SmsOrderService{
		repo:       repo,
		heroClient: heroClient,
		cfg:        cfg,
		stopCh:     make(chan struct{}),
	}
}

func (s *SmsOrderService) GetByOrderNo(ctx context.Context, orderNo string) (*SmsOrder, error) {
	if orderNo == "" {
		return nil, fmt.Errorf("order number is required")
	}
	return s.repo.GetByOrderNo(ctx, orderNo)
}

func (s *SmsOrderService) List(ctx context.Context, filter SmsOrderListFilter) (*SmsOrderListResult, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
		filter.SortDesc = true
	}
	return s.repo.List(ctx, filter)
}

func (s *SmsOrderService) BatchCreate(ctx context.Context, count int, serviceType string) ([]*SmsOrder, error) {
	if count <= 0 || count > 50 {
		return nil, fmt.Errorf("count must be between 1 and 50")
	}
	if serviceType == "" {
		serviceType = "claude"
	}

	var orders []*SmsOrder
	for i := 0; i < count; i++ {
		orderNo := generateOrderNo(serviceType)
		order := &SmsOrder{
			OrderNo:     orderNo,
			ServiceType: serviceType,
			Status:      SmsOrderStatusCreated,
		}

		created, err := s.repo.Create(ctx, order)
		if err != nil {
			return orders, fmt.Errorf("create order #%d: %w", i+1, err)
		}
		orders = append(orders, created)
	}
	return orders, nil
}

// RefreshSmsContent advances the order state machine:
//
//	created  -> [getNumberWithRetry] -> pending  | failed
//	pending  -> [getStatus]          -> received | expired (after 20min)
//	terminal states (received/expired/failed) are returned as-is.
func (s *SmsOrderService) RefreshSmsContent(ctx context.Context, orderNo string) (*SmsOrder, error) {
	order, err := s.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	switch order.Status {
	case SmsOrderStatusCreated:
		return s.assignNumber(ctx, order)
	case SmsOrderStatusPending:
		return s.pollPending(ctx, order, s.pendingTimeout())
	default:
		return order, nil
	}
}

func (s *SmsOrderService) assignNumber(ctx context.Context, order *SmsOrder) (*SmsOrder, error) {
	num, err := s.heroClient.GetNumberWithRetry(ctx, order.ServiceType)
	if err != nil {
		if updErr := s.repo.UpdateStatus(ctx, order.ID, SmsOrderStatusFailed); updErr != nil {
			slog.Error("mark sms order failed", "order_no", order.OrderNo, "error", updErr)
		}
		order.Status = SmsOrderStatusFailed
		return order, ErrSmsOrderFetchFailed
	}

	now := time.Now()
	if err := s.repo.AssignNumber(ctx, order.ID, num.Phone, num.ID, now); err != nil {
		return nil, err
	}
	order.PhoneNumber = num.Phone
	order.HeroSmsID = num.ID
	order.PendingAt = &now
	order.Status = SmsOrderStatusPending
	return order, nil
}

func (s *SmsOrderService) pollPending(ctx context.Context, order *SmsOrder, timeout time.Duration) (*SmsOrder, error) {
	if order.PendingAt != nil && time.Since(*order.PendingAt) > timeout {
		if err := s.repo.UpdateStatus(ctx, order.ID, SmsOrderStatusExpired); err != nil {
			slog.Error("expire sms order", "order_no", order.OrderNo, "error", err)
		}
		order.Status = SmsOrderStatusExpired
		return order, nil
	}

	if order.HeroSmsID == "" {
		return order, nil
	}

	status, err := s.heroClient.GetStatus(ctx, order.HeroSmsID)
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	if status.Received {
		if err := s.repo.UpdateSmsContent(ctx, order.ID, status.Text, SmsOrderStatusReceived); err != nil {
			return nil, err
		}
		order.SmsContent = status.Text
		order.Status = SmsOrderStatusReceived
	}
	return order, nil
}

func (s *SmsOrderService) pendingTimeout() time.Duration {
	timeout := time.Duration(s.cfg.HeroSms.PollTimeoutMinutes) * time.Minute
	if timeout <= 0 {
		timeout = 20 * time.Minute
	}
	return timeout
}

func (s *SmsOrderService) StartPolling(ctx context.Context) {
	interval := time.Duration(s.cfg.HeroSms.PollIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	timeout := s.pendingTimeout()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.stopCh:
				return
			case <-ticker.C:
				s.pollPendingOrders(ctx, timeout)
			}
		}
	}()
}

func (s *SmsOrderService) StopPolling() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
}

func (s *SmsOrderService) pollPendingOrders(ctx context.Context, timeout time.Duration) {
	orders, err := s.repo.ListPending(ctx)
	if err != nil {
		slog.Error("poll pending sms orders", "error", err)
		return
	}

	now := time.Now()
	for _, order := range orders {
		if order.PendingAt != nil && now.Sub(*order.PendingAt) > timeout {
			if err := s.repo.UpdateStatus(ctx, order.ID, SmsOrderStatusExpired); err != nil {
				slog.Error("expire sms order", "order_no", order.OrderNo, "error", err)
			}
			continue
		}

		if order.HeroSmsID == "" {
			continue
		}

		status, err := s.heroClient.GetStatus(ctx, order.HeroSmsID)
		if err != nil {
			slog.Error("poll sms status", "order_no", order.OrderNo, "error", err)
			continue
		}

		if status.Received {
			if err := s.repo.UpdateSmsContent(ctx, order.ID, status.Text, SmsOrderStatusReceived); err != nil {
				slog.Error("update sms content", "order_no", order.OrderNo, "error", err)
			}
		}
	}
}

func generateOrderNo(serviceType string) string {
	var prefix string
	switch serviceType {
	case "claude", "acz":
		prefix = "CC"
	case "openai", "dr":
		prefix = "OP"
	default:
		prefix = "DE"
	}
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}
