package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SmsOrderHandler handles SMS order query requests.
type SmsOrderHandler struct {
	smsOrderService *service.SmsOrderService
}

// NewSmsOrderHandler creates a new SmsOrderHandler.
func NewSmsOrderHandler(smsOrderService *service.SmsOrderService) *SmsOrderHandler {
	return &SmsOrderHandler{
		smsOrderService: smsOrderService,
	}
}

// Query handles querying an SMS order by order number.
// GET /api/v1/sms-orders/:order_no
func (h *SmsOrderHandler) Query(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "Order number is required")
		return
	}

	order, err := h.smsOrderService.GetByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.SmsOrderFromService(order))
}

// Refresh triggers a manual refresh of SMS content for an order.
// POST /api/v1/sms-orders/:order_no/refresh
//
// Optional JSON body: {"service_type": "claude" | "openai"}.
// The service_type is only honored when the order is in the `created` state
// (i.e., no phone number assigned yet). Once a number has been fetched, the
// service type is locked and the field is ignored.
func (h *SmsOrderHandler) Refresh(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "Order number is required")
		return
	}

	var req struct {
		ServiceType string `json:"service_type"`
	}
	// Body is optional; ShouldBindJSON returns an error on empty bodies which we ignore.
	_ = c.ShouldBindJSON(&req)

	if req.ServiceType != "" && req.ServiceType != "claude" && req.ServiceType != "openai" {
		response.BadRequest(c, "service_type must be 'claude' or 'openai'")
		return
	}

	order, err := h.smsOrderService.RefreshSmsContent(c.Request.Context(), orderNo, req.ServiceType)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.SmsOrderFromService(order))
}

// Reassign replaces the current phone number on a pending order with a freshly
// fetched one. Service type is locked to the order's existing value.
// POST /api/v1/sms-orders/:order_no/reassign
func (h *SmsOrderHandler) Reassign(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		response.BadRequest(c, "Order number is required")
		return
	}

	order, err := h.smsOrderService.ReassignNumber(c.Request.Context(), orderNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.SmsOrderFromService(order))
}
