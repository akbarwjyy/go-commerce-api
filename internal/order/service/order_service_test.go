package service

import (
	"testing"

	"github.com/akbarwjyy/go-commerce-api/internal/order/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/order/entity"
	"github.com/stretchr/testify/assert"
)

// Test Order Entity Methods
func TestOrderEntity_IsOwner(t *testing.T) {
	order := &entity.Order{
		ID:     1,
		UserID: 10,
	}

	assert.True(t, order.IsOwner(10))
	assert.False(t, order.IsOwner(20))
}

func TestOrderEntity_IsPending(t *testing.T) {
	order := &entity.Order{
		ID:     1,
		Status: entity.OrderStatusPending,
	}

	assert.True(t, order.IsPending())
	assert.False(t, order.IsPaid())
}

func TestOrderEntity_CanBeCancelled(t *testing.T) {
	pendingOrder := &entity.Order{Status: entity.OrderStatusPending}
	paidOrder := &entity.Order{Status: entity.OrderStatusPaid}
	shippedOrder := &entity.Order{Status: entity.OrderStatusShipped}

	assert.True(t, pendingOrder.CanBeCancelled())
	assert.False(t, paidOrder.CanBeCancelled())
	assert.False(t, shippedOrder.CanBeCancelled())
}

func TestOrderEntity_CanBeShipped(t *testing.T) {
	pendingOrder := &entity.Order{Status: entity.OrderStatusPending}
	paidOrder := &entity.Order{Status: entity.OrderStatusPaid}

	assert.False(t, pendingOrder.CanBeShipped())
	assert.True(t, paidOrder.CanBeShipped())
}

func TestOrderEntity_UpdateStatus(t *testing.T) {
	// Test PENDING -> PAID
	order := &entity.Order{Status: entity.OrderStatusPending}
	result := order.UpdateStatus(entity.OrderStatusPaid)
	assert.True(t, result)
	assert.Equal(t, entity.OrderStatusPaid, order.Status)

	// Test PAID -> SHIPPED
	result = order.UpdateStatus(entity.OrderStatusShipped)
	assert.True(t, result)
	assert.Equal(t, entity.OrderStatusShipped, order.Status)

	// Test SHIPPED -> COMPLETED
	result = order.UpdateStatus(entity.OrderStatusCompleted)
	assert.True(t, result)
	assert.Equal(t, entity.OrderStatusCompleted, order.Status)

	// Test invalid transition (COMPLETED -> PENDING)
	result = order.UpdateStatus(entity.OrderStatusPending)
	assert.False(t, result)
}

func TestOrderEntity_CalculateTotal(t *testing.T) {
	order := &entity.Order{
		Items: []entity.OrderItem{
			{Price: 100, Quantity: 2, Subtotal: 200},
			{Price: 50, Quantity: 3, Subtotal: 150},
		},
	}

	total := order.CalculateTotal()
	assert.Equal(t, 350.0, total)
	assert.Equal(t, 350.0, order.TotalAmount)
}

// Test OrderItem Entity
func TestOrderItemEntity_CalculateSubtotal(t *testing.T) {
	item := &entity.OrderItem{
		Price:    100.50,
		Quantity: 3,
	}

	subtotal := item.CalculateSubtotal()
	assert.Equal(t, 301.50, subtotal)
	assert.Equal(t, 301.50, item.Subtotal)
}

// Test Checkout Request
func TestCheckoutRequest(t *testing.T) {
	req := &dto.CheckoutRequest{
		Items: []dto.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
		ShippingAddress: "123 Main St",
		Notes:           "Leave at door",
	}

	assert.Len(t, req.Items, 2)
	assert.Equal(t, "123 Main St", req.ShippingAddress)
}

// Test Order Status Constants
func TestOrderStatusConstants(t *testing.T) {
	assert.Equal(t, "PENDING", entity.OrderStatusPending)
	assert.Equal(t, "PAID", entity.OrderStatusPaid)
	assert.Equal(t, "SHIPPED", entity.OrderStatusShipped)
	assert.Equal(t, "COMPLETED", entity.OrderStatusCompleted)
	assert.Equal(t, "CANCELLED", entity.OrderStatusCancelled)
}
