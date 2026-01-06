package service

import (
	"errors"
	"math"

	"github.com/akbarwjyy/go-commerce-api/internal/product/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/product/entity"
	"github.com/akbarwjyy/go-commerce-api/internal/product/repository"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrProductNotFound     = errors.New("product not found")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrUnauthorized        = errors.New("you are not authorized to perform this action")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrCategoryExists      = errors.New("category already exists")
	ErrInvalidStockAction  = errors.New("invalid stock action")
)

// ProductService interface untuk business logic produk
type ProductService interface {
	// Product operations
	CreateProduct(sellerID uint, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(id uint) (*dto.ProductResponse, error)
	GetAllProducts(params *dto.ProductQueryParams) (*dto.ProductListResponse, error)
	GetMyProducts(sellerID uint) ([]dto.ProductResponse, error)
	UpdateProduct(sellerID uint, productID uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(sellerID uint, productID uint) error
	UpdateStock(sellerID uint, productID uint, req *dto.UpdateStockRequest) (*dto.ProductResponse, error)

	// Category operations
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAllCategories() ([]dto.CategoryResponse, error)
	GetCategory(id uint) (*dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error

	// For inter-module communication
	GetProductByID(id uint) (*entity.Product, error)
	ReduceStock(productID uint, quantity int) error
	RestoreStock(productID uint, quantity int) error
}

// productService implementasi ProductService
type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	db           *gorm.DB
}

// NewProductService membuat instance baru ProductService
func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	db *gorm.DB,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		db:           db,
	}
}

// ========================================
// Product Operations
// ========================================

// CreateProduct membuat produk baru
func (s *productService) CreateProduct(sellerID uint, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Validate category if provided
	if req.CategoryID > 0 {
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
	}

	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		SellerID:    sellerID,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Reload product with category
	product, _ = s.productRepo.FindByIDWithCategory(product.ID)

	return s.toProductResponse(product), nil
}

// GetProduct mengambil produk berdasarkan ID
func (s *productService) GetProduct(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByIDWithCategory(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return s.toProductResponse(product), nil
}

// GetAllProducts mengambil semua produk dengan filter dan pagination
func (s *productService) GetAllProducts(params *dto.ProductQueryParams) (*dto.ProductListResponse, error) {
	// Set default pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	products, total, err := s.productRepo.FindAll(params)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *s.toProductResponse(&p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.ProductListResponse{
		Products:   productResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetMyProducts mengambil produk milik seller
func (s *productService) GetMyProducts(sellerID uint) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.FindBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, p := range products {
		responses = append(responses, *s.toProductResponse(&p))
	}
	return responses, nil
}

// UpdateProduct mengupdate produk
func (s *productService) UpdateProduct(sellerID uint, productID uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Check ownership
	if !product.IsOwner(sellerID) {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.CategoryID > 0 {
		// Validate category
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		product.CategoryID = req.CategoryID
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Reload with category
	product, _ = s.productRepo.FindByIDWithCategory(product.ID)

	return s.toProductResponse(product), nil
}

// DeleteProduct menghapus produk
func (s *productService) DeleteProduct(sellerID uint, productID uint) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	// Check ownership
	if !product.IsOwner(sellerID) {
		return ErrUnauthorized
	}

	return s.productRepo.Delete(productID)
}

// UpdateStock mengupdate stok produk
func (s *productService) UpdateStock(sellerID uint, productID uint, req *dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Check ownership
	if !product.IsOwner(sellerID) {
		return nil, ErrUnauthorized
	}

	switch req.Action {
	case "add":
		product.AddStock(req.Quantity)
	case "reduce":
		if !product.ReduceStock(req.Quantity) {
			return nil, ErrInsufficientStock
		}
	default:
		return nil, ErrInvalidStockAction
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// ========================================
// Category Operations
// ========================================

// CreateCategory membuat kategori baru
func (s *productService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Check if category exists
	existing, _ := s.categoryRepo.FindByName(req.Name)
	if existing != nil {
		return nil, ErrCategoryExists
	}

	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

// GetAllCategories mengambil semua kategori
func (s *productService) GetAllCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.CategoryResponse
	for _, c := range categories {
		responses = append(responses, *s.toCategoryResponse(&c))
	}
	return responses, nil
}

// GetCategory mengambil kategori berdasarkan ID
func (s *productService) GetCategory(id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return s.toCategoryResponse(category), nil
}

// UpdateCategory mengupdate kategori
func (s *productService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

// DeleteCategory menghapus kategori
func (s *productService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	return s.categoryRepo.Delete(id)
}

// ========================================
// Inter-Module Communication
// ========================================

// GetProductByID mengambil entity produk (untuk modul lain)
func (s *productService) GetProductByID(id uint) (*entity.Product, error) {
	return s.productRepo.FindByID(id)
}

// ReduceStock mengurangi stok (dipanggil dari Order Module)
func (s *productService) ReduceStock(productID uint, quantity int) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return ErrProductNotFound
	}

	if !product.HasStock(quantity) {
		return ErrInsufficientStock
	}

	return s.productRepo.UpdateStock(productID, -quantity)
}

// RestoreStock mengembalikan stok (jika order dibatalkan)
func (s *productService) RestoreStock(productID uint, quantity int) error {
	return s.productRepo.UpdateStock(productID, quantity)
}

// ========================================
// Helper Functions
// ========================================

func (s *productService) toProductResponse(p *entity.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		SellerID:    p.SellerID,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
	}

	if p.Category != nil {
		resp.Category = s.toCategoryResponse(p.Category)
	}

	return resp
}

func (s *productService) toCategoryResponse(c *entity.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}
