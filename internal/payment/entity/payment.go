package entity

import (
	"time"

	"gorm.io/gorm"
)

// Payment status constants
const (
	PaymentStatusPending    = "PENDING"
	PaymentStatusProcessing = "PROCESSING"
	PaymentStatusSuccess    = "SUCCESS"
	PaymentStatusFailed     = "FAILED"
)

// Payment method constants
const (
	PaymentMethodBankTransfer = "BANK_TRANSFER"
	PaymentMethodCreditCard   = "CREDIT_CARD"
	PaymentMethodEWallet      = "E_WALLET"
)

// Payment entity untuk tabel payments
type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       uint           `gorm:"index;not null" json:"order_id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Amount        float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method        string         `gorm:"size:50;not null" json:"method"`
	Status        string         `gorm:"size:20;default:PENDING" json:"status"`
	TransactionID string         `gorm:"size:100;uniqueIndex" json:"transaction_id"`
	PaidAt        *time.Time     `json:"paid_at,omitempty"`
	FailedReason  string         `gorm:"size:255" json:"failed_reason,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (Payment) TableName() string {
	return "payments"
}

// IsPending mengecek apakah payment masih pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsProcessing mengecek apakah payment sedang diproses
func (p *Payment) IsProcessing() bool {
	return p.Status == PaymentStatusProcessing
}

// IsSuccess mengecek apakah payment berhasil
func (p *Payment) IsSuccess() bool {
	return p.Status == PaymentStatusSuccess
}

// IsFailed mengecek apakah payment gagal
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// MarkAsProcessing mengubah status menjadi processing
func (p *Payment) MarkAsProcessing() {
	p.Status = PaymentStatusProcessing
}

// MarkAsSuccess mengubah status menjadi success
func (p *Payment) MarkAsSuccess() {
	p.Status = PaymentStatusSuccess
	now := time.Now()
	p.PaidAt = &now
}

// MarkAsFailed mengubah status menjadi failed
func (p *Payment) MarkAsFailed(reason string) {
	p.Status = PaymentStatusFailed
	p.FailedReason = reason
}

// IsValidMethod memvalidasi method payment
func IsValidMethod(method string) bool {
	return method == PaymentMethodBankTransfer ||
		method == PaymentMethodCreditCard ||
		method == PaymentMethodEWallet
}
