package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category entity untuk tabel categories
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Products    []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName menentukan nama tabel di database
func (Category) TableName() string {
	return "categories"
}
