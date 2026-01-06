package dto

// CreateProductRequest untuk request membuat produk baru
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest untuk request update produk
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateStockRequest untuk request update stok
type UpdateStockRequest struct {
	Quantity int    `json:"quantity" binding:"required"`
	Action   string `json:"action" binding:"required,oneof=add reduce"`
}

// ProductResponse untuk response data produk
type ProductResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	CategoryID  uint              `json:"category_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	SellerID    uint              `json:"seller_id"`
	ImageURL    string            `json:"image_url,omitempty"`
	IsActive    bool              `json:"is_active"`
}

// ProductListResponse untuk response list produk dengan pagination
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// CreateCategoryRequest untuk request membuat kategori baru
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

// UpdateCategoryRequest untuk request update kategori
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

// CategoryResponse untuk response data kategori
type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProductQueryParams untuk filter dan pagination
type ProductQueryParams struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	CategoryID uint   `form:"category_id"`
	SellerID   uint   `form:"seller_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	IsActive   *bool  `form:"is_active"`
}
