package handler

import (
	"strconv"

	authEntity "github.com/akbarwjyy/go-commerce-api/internal/auth/entity"
	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	"github.com/akbarwjyy/go-commerce-api/internal/order/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/order/service"
	"github.com/gin-gonic/gin"
)

// OrderHandler menangani HTTP request untuk order
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler membuat instance baru OrderHandler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Checkout godoc
// @Summary      Checkout order
// @Description  Create a new order from cart items
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CheckoutRequest true "Checkout request"
// @Success      201 {object} response.APIResponse{data=dto.OrderResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /orders/checkout [post]
func (h *OrderHandler) Checkout(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req dto.CheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.orderService.Checkout(userID.(uint), &req)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			response.NotFound(ctx, "One or more products not found")
		case service.ErrInsufficientStock:
			response.BadRequest(ctx, "Insufficient stock for one or more products", nil)
		case service.ErrEmptyCart:
			response.BadRequest(ctx, "Cart is empty", nil)
		default:
			response.InternalServerError(ctx, "Failed to checkout", err.Error())
		}
		return
	}

	response.Created(ctx, "Order created successfully", result)
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Get a single order by its ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200 {object} response.APIResponse{data=dto.OrderResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrder(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid order ID", nil)
		return
	}

	result, err := h.orderService.GetOrder(userID.(uint), uint(id))
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			response.NotFound(ctx, "Order not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to view this order")
		default:
			response.InternalServerError(ctx, "Failed to get order", err.Error())
		}
		return
	}

	response.OK(ctx, "Order retrieved successfully", result)
}

// GetMyOrders godoc
// @Summary      Get my orders
// @Description  Get orders belonging to the current user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        status query string false "Filter by status" Enums(PENDING, PAID, SHIPPED, COMPLETED, CANCELLED)
// @Success      200 {object} response.APIResponse{data=dto.OrderListResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /orders [get]
func (h *OrderHandler) GetMyOrders(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var params dto.OrderQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.orderService.GetMyOrders(userID.(uint), &params)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get orders", err.Error())
		return
	}

	response.OK(ctx, "Orders retrieved successfully", result)
}

// GetAllOrders godoc
// @Summary      Get all orders (Admin)
// @Description  Get all orders with filters and pagination (Admin only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        status query string false "Filter by status" Enums(PENDING, PAID, SHIPPED, COMPLETED, CANCELLED)
// @Success      200 {object} response.APIResponse{data=dto.OrderListResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Router       /admin/orders [get]
func (h *OrderHandler) GetAllOrders(ctx *gin.Context) {
	var params dto.OrderQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.orderService.GetAllOrders(&params)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get orders", err.Error())
		return
	}

	response.OK(ctx, "Orders retrieved successfully", result)
}

// UpdateOrderStatus godoc
// @Summary      Update order status
// @Description  Update the status of an order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Param        request body dto.UpdateOrderStatusRequest true "Update status request"
// @Success      200 {object} response.APIResponse{data=dto.OrderResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	userRole, _ := ctx.Get("userRole")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid order ID", nil)
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	isAdmin := userRole.(string) == authEntity.RoleAdmin
	result, err := h.orderService.UpdateOrderStatus(userID.(uint), uint(id), req.Status, isAdmin)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			response.NotFound(ctx, "Order not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to update this order")
		case service.ErrInvalidStatus:
			response.BadRequest(ctx, "Invalid status transition", nil)
		default:
			response.InternalServerError(ctx, "Failed to update order status", err.Error())
		}
		return
	}

	response.OK(ctx, "Order status updated successfully", result)
}

// CancelOrder godoc
// @Summary      Cancel order
// @Description  Cancel an order and restore stock
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200 {object} response.APIResponse
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid order ID", nil)
		return
	}

	if err := h.orderService.CancelOrder(userID.(uint), uint(id)); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			response.NotFound(ctx, "Order not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to cancel this order")
		case service.ErrOrderNotCancellable:
			response.BadRequest(ctx, "Order cannot be cancelled", nil)
		default:
			response.InternalServerError(ctx, "Failed to cancel order", err.Error())
		}
		return
	}

	response.OK(ctx, "Order cancelled successfully", nil)
}
