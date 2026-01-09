package repository

import (
	"github.com/akbarwjyy/go-commerce-api/internal/order/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/order/entity"
	"gorm.io/gorm"
)

// OrderRepository interface untuk akses data order
type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(id uint) (*entity.Order, error)
	FindByIDWithItems(id uint) (*entity.Order, error)
	FindByUserID(userID uint, params *dto.OrderQueryParams) ([]entity.Order, int64, error)
	FindAll(params *dto.OrderQueryParams) ([]entity.Order, int64, error)
	Update(order *entity.Order) error
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) OrderRepository
}

// orderRepository implementasi OrderRepository
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository membuat instance baru OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// WithTx mengembalikan repository dengan transaction
func (r *orderRepository) WithTx(tx *gorm.DB) OrderRepository {
	return &orderRepository{db: tx}
}

// Create menyimpan order baru ke database
func (r *orderRepository) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

// FindByID mencari order berdasarkan ID
func (r *orderRepository) FindByID(id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByIDWithItems mencari order dengan items
func (r *orderRepository) FindByIDWithItems(id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("Items").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID mengambil order berdasarkan user ID dengan pagination
func (r *orderRepository) FindByUserID(userID uint, params *dto.OrderQueryParams) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	query := r.db.Model(&entity.Order{}).Where("user_id = ?", userID)

	// Apply status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Preload("Items").Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// FindAll mengambil semua order dengan pagination (untuk admin)
func (r *orderRepository) FindAll(params *dto.OrderQueryParams) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	query := r.db.Model(&entity.Order{})

	// Apply status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Preload("Items").Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Update mengupdate data order
func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus mengupdate status order
func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&entity.Order{}).Where("id = ?", id).Update("status", status).Error
}

// Delete menghapus order (soft delete)
func (r *orderRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Order{}, id).Error
}
