package dto

// CreatePaymentRequest untuk request membuat payment
type CreatePaymentRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required,oneof=BANK_TRANSFER CREDIT_CARD E_WALLET"`
}

// PaymentResponse untuk response data payment
type PaymentResponse struct {
	ID            uint    `json:"id"`
	OrderID       uint    `json:"order_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Method        string  `json:"method"`
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	PaidAt        string  `json:"paid_at,omitempty"`
	FailedReason  string  `json:"failed_reason,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// PaymentListResponse untuk response list payment
type PaymentListResponse struct {
	Payments   []PaymentResponse `json:"payments"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// PaymentQueryParams untuk filter dan pagination
type PaymentQueryParams struct {
	Page    int    `form:"page,default=1"`
	Limit   int    `form:"limit,default=10"`
	Status  string `form:"status"`
	OrderID uint   `form:"order_id"`
}

// PaymentCallbackRequest untuk simulasi callback dari payment gateway
type PaymentCallbackRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	Status        string `json:"status" binding:"required,oneof=SUCCESS FAILED"`
	FailedReason  string `json:"failed_reason,omitempty"`
}
