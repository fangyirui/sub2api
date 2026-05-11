package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminSmsOrderHandler struct {
	smsOrderService *service.SmsOrderService
}

func NewAdminSmsOrderHandler(smsOrderService *service.SmsOrderService) *AdminSmsOrderHandler {
	return &AdminSmsOrderHandler{smsOrderService: smsOrderService}
}

type batchGenerateRequest struct {
	Count       int    `json:"count" binding:"required,min=1,max=50"`
	ServiceType string `json:"service_type"`
}

type adminSmsOrderResponse struct {
	OrderNo     string  `json:"order_no"`
	ServiceType string  `json:"service_type"`
	PhoneNumber string  `json:"phone_number"`
	SmsContent  string  `json:"sms_content"`
	Status      string  `json:"status"`
	PendingAt   *string `json:"pending_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func toAdminSmsOrderResponse(o *service.SmsOrder) adminSmsOrderResponse {
	resp := adminSmsOrderResponse{
		OrderNo:     o.OrderNo,
		ServiceType: o.ServiceType,
		PhoneNumber: o.PhoneNumber,
		SmsContent:  o.SmsContent,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if o.PendingAt != nil {
		s := o.PendingAt.Format("2006-01-02 15:04:05")
		resp.PendingAt = &s
	}
	return resp
}

// BatchGenerate creates multiple SMS orders without calling hero-sms.
// Customers will trigger getNumber lazily on first refresh.
// POST /api/v1/admin/sms-orders/generate
func (h *AdminSmsOrderHandler) BatchGenerate(c *gin.Context) {
	var req batchGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: count must be between 1 and 50")
		return
	}

	orders, err := h.smsOrderService.BatchCreate(c.Request.Context(), req.Count, req.ServiceType)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	result := make([]adminSmsOrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toAdminSmsOrderResponse(o))
	}

	response.Success(c, gin.H{"items": result})
}

type adminSmsOrderListResponse struct {
	Items []adminSmsOrderResponse `json:"items"`
	Total int                     `json:"total"`
}

// List returns SMS orders with pagination, date filtering, and sorting.
// GET /api/v1/admin/sms-orders?page=1&page_size=20&status=created&start_date=2026-05-11&end_date=2026-05-11&timezone=Asia/Shanghai&sort_by=pending_at&sort_order=desc
func (h *AdminSmsOrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	userTZ := c.Query("timezone")
	sortBy := c.DefaultQuery("sort_by", "pending_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	filter := service.SmsOrderListFilter{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		SortBy:   sortBy,
		SortDesc: sortOrder == "desc",
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", startDateStr, userTZ)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
			return
		}
		filter.StartTime = &t
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		t, err := timezone.ParseInUserLocation("2006-01-02", endDateStr, userTZ)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
			return
		}
		t = t.AddDate(0, 0, 1)
		filter.EndTime = &t
	}

	res, err := h.smsOrderService.List(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	items := make([]adminSmsOrderResponse, 0, len(res.Items))
	for _, o := range res.Items {
		items = append(items, toAdminSmsOrderResponse(o))
	}

	response.Success(c, adminSmsOrderListResponse{
		Items: items,
		Total: res.Total,
	})
}
