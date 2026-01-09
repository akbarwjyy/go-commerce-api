package dto

// OrderItemRequest untuk request item dalam checkout
type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// CheckoutRequest untuk request checkout
type CheckoutRequest struct {
	Items           []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
	ShippingAddress string             `json:"shipping_address" binding:"required"`
	Notes           string             `json:"notes,omitempty"`
}

// UpdateOrderStatusRequest untuk request update status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PAID SHIPPED COMPLETED CANCELLED"`
}

// OrderItemResponse untuk response item dalam order
type OrderItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name,omitempty"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

// OrderResponse untuk response data order
type OrderResponse struct {
	ID              uint                `json:"id"`
	UserID          uint                `json:"user_id"`
	TotalAmount     float64             `json:"total_amount"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
	Notes           string              `json:"notes,omitempty"`
	Items           []OrderItemResponse `json:"items,omitempty"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

// OrderListResponse untuk response list order dengan pagination
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// OrderQueryParams untuk filter dan pagination
type OrderQueryParams struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Status string `form:"status"`
}
