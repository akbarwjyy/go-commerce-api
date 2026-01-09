package entity

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem entity untuk tabel order_items
type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"index;not null" json:"order_id"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	Subtotal  float64        `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (OrderItem) TableName() string {
	return "order_items"
}

// CalculateSubtotal menghitung subtotal item
func (oi *OrderItem) CalculateSubtotal() float64 {
	oi.Subtotal = oi.Price * float64(oi.Quantity)
	return oi.Subtotal
}
