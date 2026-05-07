package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// SmsOrderResponse is the JSON response for an SMS order query.
type SmsOrderResponse struct {
	OrderNo     string    `json:"order_no"`
	PhoneNumber string    `json:"phone_number"`
	SmsContent  string    `json:"sms_content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SmsOrderFromService maps a service-layer SmsOrder to the DTO response.
func SmsOrderFromService(o *service.SmsOrder) *SmsOrderResponse {
	if o == nil {
		return nil
	}
	return &SmsOrderResponse{
		OrderNo:     o.OrderNo,
		PhoneNumber: o.PhoneNumber,
		SmsContent:  o.SmsContent,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}
