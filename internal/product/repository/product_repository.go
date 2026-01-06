package repository

import (
	"github.com/akbarwjyy/go-commerce-api/internal/product/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/product/entity"
	"gorm.io/gorm"
)

// ProductRepository interface untuk akses data produk
type ProductRepository interface {
	Create(product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	FindByIDWithCategory(id uint) (*entity.Product, error)
	FindAll(params *dto.ProductQueryParams) ([]entity.Product, int64, error)
	FindBySellerID(sellerID uint) ([]entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
	WithTx(tx *gorm.DB) ProductRepository
}

// productRepository implementasi ProductRepository
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository membuat instance baru ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// WithTx mengembalikan repository dengan transaction
func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}

// Create menyimpan produk baru ke database
func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

// FindByID mencari produk berdasarkan ID
func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByIDWithCategory mencari produk dengan relasi kategori
func (r *productRepository) FindByIDWithCategory(id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// FindAll mengambil semua produk dengan filter dan pagination
func (r *productRepository) FindAll(params *dto.ProductQueryParams) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.Model(&entity.Product{})

	// Apply filters
	if params.Search != "" {
		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
	}
	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.SellerID > 0 {
		query = query.Where("seller_id = ?", params.SellerID)
	}
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Preload("Category").Offset(offset).Limit(params.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindBySellerID mengambil produk berdasarkan seller ID
func (r *productRepository) FindBySellerID(sellerID uint) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.Where("seller_id = ?", sellerID).Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Update mengupdate data produk
func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

// Delete menghapus produk (soft delete)
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}

// UpdateStock mengupdate stok produk dengan row-level locking
func (r *productRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&entity.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
