package entity

import (
	"time"

	"gorm.io/gorm"
)

// Order status constants
const (
	OrderStatusPending   = "PENDING"
	OrderStatusPaid      = "PAID"
	OrderStatusShipped   = "SHIPPED"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCancelled = "CANCELLED"
)

// Order entity untuk tabel orders
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	TotalAmount   float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status        string         `gorm:"size:20;default:PENDING" json:"status"`
	ShippingAddr  string         `gorm:"type:text" json:"shipping_address"`
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName menentukan nama tabel di database
func (Order) TableName() string {
	return "orders"
}

// IsOwner mengecek apakah user adalah pemilik order
func (o *Order) IsOwner(userID uint) bool {
	return o.UserID == userID
}

// IsPending mengecek apakah order masih pending
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid mengecek apakah order sudah dibayar
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// CanBeCancelled mengecek apakah order bisa dibatalkan
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending
}

// CanBeShipped mengecek apakah order bisa dikirim
func (o *Order) CanBeShipped() bool {
	return o.Status == OrderStatusPaid
}

// CanBeCompleted mengecek apakah order bisa diselesaikan
func (o *Order) CanBeCompleted() bool {
	return o.Status == OrderStatusShipped
}

// UpdateStatus mengupdate status order
func (o *Order) UpdateStatus(newStatus string) bool {
	switch newStatus {
	case OrderStatusPaid:
		if o.IsPending() {
			o.Status = OrderStatusPaid
			return true
		}
	case OrderStatusShipped:
		if o.CanBeShipped() {
			o.Status = OrderStatusShipped
			return true
		}
	case OrderStatusCompleted:
		if o.CanBeCompleted() {
			o.Status = OrderStatusCompleted
			return true
		}
	case OrderStatusCancelled:
		if o.CanBeCancelled() {
			o.Status = OrderStatusCancelled
			return true
		}
	}
	return false
}

// CalculateTotal menghitung total dari semua items
func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.TotalAmount = total
	return total
}
