package entity

import (
	"time"

	"gorm.io/gorm"
)

// Role constants
const (
	RoleAdmin  = "admin"
	RoleSeller = "seller"
	RoleUser   = "user"
)

// User entity untuk tabel users
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName menentukan nama tabel di database
func (User) TableName() string {
	return "users"
}

// IsAdmin mengecek apakah user adalah admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsSeller mengecek apakah user adalah seller
func (u *User) IsSeller() bool {
	return u.Role == RoleSeller
}

// IsValidRole memvalidasi role yang valid
func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleSeller || role == RoleUser
}
