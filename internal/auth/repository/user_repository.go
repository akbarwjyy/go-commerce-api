package repository

import (
	"github.com/akbar/go-commerce-api/internal/auth/entity"
	"gorm.io/gorm"
)

// UserRepository interface untuk akses data user
type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
}

// userRepository implementasi UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create menyimpan user baru ke database
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindByID mencari user berdasarkan ID
func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email
func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update mengupdate data user
func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete menghapus user (soft delete)
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
