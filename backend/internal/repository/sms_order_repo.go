package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/smsorder"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type smsOrderRepository struct {
	client *dbent.Client
}

// NewSmsOrderRepository creates a new SMS order repository.
func NewSmsOrderRepository(client *dbent.Client) service.SmsOrderRepository {
	return &smsOrderRepository{client: client}
}

func (r *smsOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*service.SmsOrder, error) {
	m, err := r.client.SmsOrder.Query().
		Where(smsorder.OrderNoEQ(orderNo)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSmsOrderNotFound, nil)
	}
	return smsOrderEntityToService(m), nil
}

func smsOrderEntityToService(m *dbent.SmsOrder) *service.SmsOrder {
	if m == nil {
		return nil
	}
	return &service.SmsOrder{
		ID:          m.ID,
		OrderNo:     m.OrderNo,
		PhoneNumber: m.PhoneNumber,
		SmsContent:  m.SmsContent,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
