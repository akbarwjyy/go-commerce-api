package repository

import (
	"github.com/akbarwjyy/go-commerce-api/internal/payment/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/entity"
	"gorm.io/gorm"
)

// PaymentRepository interface untuk akses data payment
type PaymentRepository interface {
	Create(payment *entity.Payment) error
	FindByID(id uint) (*entity.Payment, error)
	FindByOrderID(orderID uint) (*entity.Payment, error)
	FindByTransactionID(transactionID string) (*entity.Payment, error)
	FindByUserID(userID uint, params *dto.PaymentQueryParams) ([]entity.Payment, int64, error)
	FindAll(params *dto.PaymentQueryParams) ([]entity.Payment, int64, error)
	Update(payment *entity.Payment) error
	WithTx(tx *gorm.DB) PaymentRepository
}

// paymentRepository implementasi PaymentRepository
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository membuat instance baru PaymentRepository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// WithTx mengembalikan repository dengan transaction
func (r *paymentRepository) WithTx(tx *gorm.DB) PaymentRepository {
	return &paymentRepository{db: tx}
}

// Create menyimpan payment baru ke database
func (r *paymentRepository) Create(payment *entity.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID mencari payment berdasarkan ID
func (r *paymentRepository) FindByID(id uint) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByOrderID mencari payment berdasarkan Order ID
func (r *paymentRepository) FindByOrderID(orderID uint) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByTransactionID mencari payment berdasarkan Transaction ID
func (r *paymentRepository) FindByTransactionID(transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByUserID mengambil payment berdasarkan user ID dengan pagination
func (r *paymentRepository) FindByUserID(userID uint, params *dto.PaymentQueryParams) ([]entity.Payment, int64, error) {
	var payments []entity.Payment
	var total int64

	query := r.db.Model(&entity.Payment{}).Where("user_id = ?", userID)

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.OrderID > 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// FindAll mengambil semua payment dengan pagination (untuk admin)
func (r *paymentRepository) FindAll(params *dto.PaymentQueryParams) ([]entity.Payment, int64, error) {
	var payments []entity.Payment
	var total int64

	query := r.db.Model(&entity.Payment{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.OrderID > 0 {
		query = query.Where("order_id = ?", params.OrderID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// Update mengupdate data payment
func (r *paymentRepository) Update(payment *entity.Payment) error {
	return r.db.Save(payment).Error
}
