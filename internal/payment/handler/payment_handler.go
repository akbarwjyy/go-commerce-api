package handler

import (
	"strconv"

	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/service"
	"github.com/gin-gonic/gin"
)

// PaymentHandler menangani HTTP request untuk payment
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler membuat instance baru PaymentHandler
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment godoc
// @Summary      Create payment
// @Description  Create a new payment for an order (triggers async processing)
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreatePaymentRequest true "Create payment request"
// @Success      201 {object} response.APIResponse{data=dto.PaymentResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /payments [post]
func (h *PaymentHandler) CreatePayment(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req dto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.paymentService.CreatePayment(userID.(uint), &req)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			response.NotFound(ctx, "Order not found")
		case service.ErrOrderNotPending:
			response.BadRequest(ctx, "Order is not in pending status", nil)
		case service.ErrPaymentAlreadyExists:
			response.BadRequest(ctx, "Payment already exists for this order", nil)
		case service.ErrInvalidPaymentMethod:
			response.BadRequest(ctx, "Invalid payment method", nil)
		default:
			response.InternalServerError(ctx, "Failed to create payment", err.Error())
		}
		return
	}

	response.Created(ctx, "Payment created successfully. Processing async...", result)
}

// GetPayment godoc
// @Summary      Get payment by ID
// @Description  Get a single payment by its ID
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Payment ID"
// @Success      200 {object} response.APIResponse{data=dto.PaymentResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetPayment(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid payment ID", nil)
		return
	}

	result, err := h.paymentService.GetPayment(userID.(uint), uint(id))
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			response.NotFound(ctx, "Payment not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to view this payment")
		default:
			response.InternalServerError(ctx, "Failed to get payment", err.Error())
		}
		return
	}

	response.OK(ctx, "Payment retrieved successfully", result)
}

// GetMyPayments godoc
// @Summary      Get my payments
// @Description  Get payments belonging to the current user
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        status query string false "Filter by status" Enums(PENDING, PROCESSING, SUCCESS, FAILED)
// @Success      200 {object} response.APIResponse{data=dto.PaymentListResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /payments [get]
func (h *PaymentHandler) GetMyPayments(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var params dto.PaymentQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.paymentService.GetMyPayments(userID.(uint), &params)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get payments", err.Error())
		return
	}

	response.OK(ctx, "Payments retrieved successfully", result)
}

// GetAllPayments godoc
// @Summary      Get all payments (Admin)
// @Description  Get all payments with filters and pagination (Admin only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        status query string false "Filter by status" Enums(PENDING, PROCESSING, SUCCESS, FAILED)
// @Success      200 {object} response.APIResponse{data=dto.PaymentListResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Router       /admin/payments [get]
func (h *PaymentHandler) GetAllPayments(ctx *gin.Context) {
	var params dto.PaymentQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.paymentService.GetAllPayments(&params)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get payments", err.Error())
		return
	}

	response.OK(ctx, "Payments retrieved successfully", result)
}

// GetPaymentByOrder godoc
// @Summary      Get payment by order ID
// @Description  Get the payment associated with an order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200 {object} response.APIResponse{data=dto.PaymentResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /orders/{id}/payment [get]
func (h *PaymentHandler) GetPaymentByOrder(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid order ID", nil)
		return
	}

	result, err := h.paymentService.GetPaymentByOrderID(uint(id))
	if err != nil {
		if err == service.ErrPaymentNotFound {
			response.NotFound(ctx, "Payment not found for this order")
			return
		}
		response.InternalServerError(ctx, "Failed to get payment", err.Error())
		return
	}

	response.OK(ctx, "Payment retrieved successfully", result)
}

// PaymentCallback godoc
// @Summary      Payment callback (Testing)
// @Description  Manual payment callback for testing purposes
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.PaymentCallbackRequest true "Payment callback request"
// @Success      200 {object} response.APIResponse
// @Failure      400 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /payments/callback [post]
func (h *PaymentHandler) PaymentCallback(ctx *gin.Context) {
	var req dto.PaymentCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	if err := h.paymentService.ProcessPaymentCallback(req.TransactionID, req.Status, req.FailedReason); err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			response.NotFound(ctx, "Payment not found")
		case service.ErrPaymentAlreadyProcessed:
			response.BadRequest(ctx, "Payment has already been processed", nil)
		default:
			response.InternalServerError(ctx, "Failed to process callback", err.Error())
		}
		return
	}

	response.OK(ctx, "Payment callback processed successfully", nil)
}
