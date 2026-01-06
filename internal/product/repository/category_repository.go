package repository

import (
	"github.com/akbarwjyy/go-commerce-api/internal/product/entity"
	"gorm.io/gorm"
)

// CategoryRepository interface untuk akses data kategori
type CategoryRepository interface {
	Create(category *entity.Category) error
	FindByID(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	FindAll() ([]entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
}

// categoryRepository implementasi CategoryRepository
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository membuat instance baru CategoryRepository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create menyimpan kategori baru ke database
func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

// FindByID mencari kategori berdasarkan ID
func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName mencari kategori berdasarkan nama
func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll mengambil semua kategori
func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update mengupdate data kategori
func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete menghapus kategori (soft delete)
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}
