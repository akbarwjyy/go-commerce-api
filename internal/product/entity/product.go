package entity

import (
	"time"

	"gorm.io/gorm"
)

// Product entity untuk tabel products
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	CategoryID  uint           `gorm:"index" json:"category_id"`
	SellerID    uint           `gorm:"index;not null" json:"seller_id"`
	ImageURL    string         `gorm:"size:255" json:"image_url,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName menentukan nama tabel di database
func (Product) TableName() string {
	return "products"
}

// IsOwner mengecek apakah user adalah pemilik produk
func (p *Product) IsOwner(userID uint) bool {
	return p.SellerID == userID
}

// HasStock mengecek apakah stok tersedia
func (p *Product) HasStock(quantity int) bool {
	return p.Stock >= quantity
}

// ReduceStock mengurangi stok produk
func (p *Product) ReduceStock(quantity int) bool {
	if !p.HasStock(quantity) {
		return false
	}
	p.Stock -= quantity
	return true
}

// AddStock menambah stok produk
func (p *Product) AddStock(quantity int) {
	p.Stock += quantity
}
