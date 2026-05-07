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
