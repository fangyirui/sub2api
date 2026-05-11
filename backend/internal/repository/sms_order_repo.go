package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/smsorder"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type smsOrderRepository struct {
	client *dbent.Client
}

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

func (r *smsOrderRepository) Create(ctx context.Context, order *service.SmsOrder) (*service.SmsOrder, error) {
	builder := r.client.SmsOrder.Create().
		SetOrderNo(order.OrderNo).
		SetServiceType(order.ServiceType).
		SetPhoneNumber(order.PhoneNumber).
		SetHeroSmsID(order.HeroSmsID).
		SetStatus(order.Status)
	if order.PendingAt != nil {
		builder = builder.SetPendingAt(*order.PendingAt)
	}
	m, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return smsOrderEntityToService(m), nil
}

func (r *smsOrderRepository) List(ctx context.Context, filter service.SmsOrderListFilter) (*service.SmsOrderListResult, error) {
	q := r.client.SmsOrder.Query()
	if filter.Status != "" {
		q = q.Where(smsorder.StatusEQ(filter.Status))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	items, err := q.
		Order(dbent.Desc(smsorder.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := &service.SmsOrderListResult{
		Items: make([]*service.SmsOrder, 0, len(items)),
		Total: total,
	}
	for _, m := range items {
		result.Items = append(result.Items, smsOrderEntityToService(m))
	}
	return result, nil
}

func (r *smsOrderRepository) ListPending(ctx context.Context) ([]*service.SmsOrder, error) {
	items, err := r.client.SmsOrder.Query().
		Where(smsorder.StatusEQ(service.SmsOrderStatusPending)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*service.SmsOrder, 0, len(items))
	for _, m := range items {
		result = append(result, smsOrderEntityToService(m))
	}
	return result, nil
}

func (r *smsOrderRepository) AssignNumber(ctx context.Context, id int64, phone, heroSmsID string, pendingAt time.Time) error {
	return r.client.SmsOrder.UpdateOneID(id).
		SetPhoneNumber(phone).
		SetHeroSmsID(heroSmsID).
		SetPendingAt(pendingAt).
		SetStatus(service.SmsOrderStatusPending).
		Exec(ctx)
}

func (r *smsOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.client.SmsOrder.UpdateOneID(id).
		SetStatus(status).
		Exec(ctx)
}

func (r *smsOrderRepository) UpdateSmsContent(ctx context.Context, id int64, content string, status string) error {
	return r.client.SmsOrder.UpdateOneID(id).
		SetSmsContent(content).
		SetStatus(status).
		Exec(ctx)
}

func smsOrderEntityToService(m *dbent.SmsOrder) *service.SmsOrder {
	if m == nil {
		return nil
	}
	return &service.SmsOrder{
		ID:          m.ID,
		OrderNo:     m.OrderNo,
		ServiceType: m.ServiceType,
		PhoneNumber: m.PhoneNumber,
		HeroSmsID:   m.HeroSmsID,
		SmsContent:  m.SmsContent,
		Status:      m.Status,
		PendingAt:   m.PendingAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
