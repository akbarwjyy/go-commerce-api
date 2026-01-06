package handler

import (
	"net/http"
	"strconv"

	"github.com/akbarwjyy/go-commerce-api/internal/common/response"
	"github.com/akbarwjyy/go-commerce-api/internal/product/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/product/service"
	"github.com/gin-gonic/gin"
)

// ProductHandler menangani HTTP request untuk produk
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler membuat instance baru ProductHandler
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// ========================================
// Product Handlers
// ========================================

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product (Seller/Admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateProductRequest true "Create product request"
// @Success      201 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	sellerID, _ := ctx.Get("userID")

	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.productService.CreateProduct(sellerID.(uint), &req)
	if err != nil {
		switch err {
		case service.ErrCategoryNotFound:
			response.NotFound(ctx, "Category not found")
		default:
			response.InternalServerError(ctx, "Failed to create product", err.Error())
		}
		return
	}

	response.Created(ctx, "Product created successfully", result)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Get a single product by its ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid product ID", nil)
		return
	}

	result, err := h.productService.GetProduct(uint(id))
	if err != nil {
		if err == service.ErrProductNotFound {
			response.NotFound(ctx, "Product not found")
			return
		}
		response.InternalServerError(ctx, "Failed to get product", err.Error())
		return
	}

	response.OK(ctx, "Product retrieved successfully", result)
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Get all products with filters and pagination
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        search query string false "Search by name"
// @Param        category_id query int false "Filter by category ID"
// @Param        min_price query number false "Minimum price"
// @Param        max_price query number false "Maximum price"
// @Success      200 {object} response.APIResponse{data=dto.ProductListResponse}
// @Failure      400 {object} response.APIResponse
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	var params dto.ProductQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.productService.GetAllProducts(&params)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get products", err.Error())
		return
	}

	response.OK(ctx, "Products retrieved successfully", result)
}

// GetMyProducts godoc
// @Summary      Get my products
// @Description  Get products owned by the current seller
// @Tags         Seller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse{data=[]dto.ProductResponse}
// @Failure      401 {object} response.APIResponse
// @Router       /seller/products [get]
func (h *ProductHandler) GetMyProducts(ctx *gin.Context) {
	sellerID, _ := ctx.Get("userID")

	result, err := h.productService.GetMyProducts(sellerID.(uint))
	if err != nil {
		response.InternalServerError(ctx, "Failed to get products", err.Error())
		return
	}

	response.OK(ctx, "Products retrieved successfully", result)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update a product (Owner only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Param        request body dto.UpdateProductRequest true "Update product request"
// @Success      200 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	sellerID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid product ID", nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.productService.UpdateProduct(sellerID.(uint), uint(id), &req)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			response.NotFound(ctx, "Product not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to update this product")
		case service.ErrCategoryNotFound:
			response.NotFound(ctx, "Category not found")
		default:
			response.InternalServerError(ctx, "Failed to update product", err.Error())
		}
		return
	}

	response.OK(ctx, "Product updated successfully", result)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Delete a product (Owner only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Success      200 {object} response.APIResponse
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	sellerID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid product ID", nil)
		return
	}

	if err := h.productService.DeleteProduct(sellerID.(uint), uint(id)); err != nil {
		switch err {
		case service.ErrProductNotFound:
			response.NotFound(ctx, "Product not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to delete this product")
		default:
			response.InternalServerError(ctx, "Failed to delete product", err.Error())
		}
		return
	}

	response.OK(ctx, "Product deleted successfully", nil)
}

// UpdateStock godoc
// @Summary      Update product stock
// @Description  Add or reduce product stock (Owner only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Param        request body dto.UpdateStockRequest true "Update stock request"
// @Success      200 {object} response.APIResponse{data=dto.ProductResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /products/{id}/stock [patch]
func (h *ProductHandler) UpdateStock(ctx *gin.Context) {
	sellerID, _ := ctx.Get("userID")

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid product ID", nil)
		return
	}

	var req dto.UpdateStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.productService.UpdateStock(sellerID.(uint), uint(id), &req)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			response.NotFound(ctx, "Product not found")
		case service.ErrUnauthorized:
			response.Forbidden(ctx, "You are not authorized to update this product")
		case service.ErrInsufficientStock:
			response.BadRequest(ctx, "Insufficient stock", nil)
		case service.ErrInvalidStockAction:
			response.BadRequest(ctx, "Invalid stock action. Use 'add' or 'reduce'", nil)
		default:
			response.InternalServerError(ctx, "Failed to update stock", err.Error())
		}
		return
	}

	response.OK(ctx, "Stock updated successfully", result)
}

// ========================================
// Category Handlers
// ========================================

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new product category (Admin only)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCategoryRequest true "Create category request"
// @Success      201 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      409 {object} response.APIResponse
// @Router       /categories [post]
func (h *ProductHandler) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.productService.CreateCategory(&req)
	if err != nil {
		if err == service.ErrCategoryExists {
			response.Error(ctx, http.StatusConflict, "Category already exists", nil)
			return
		}
		response.InternalServerError(ctx, "Failed to create category", err.Error())
		return
	}

	response.Created(ctx, "Category created successfully", result)
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Get all product categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200 {object} response.APIResponse{data=[]dto.CategoryResponse}
// @Router       /categories [get]
func (h *ProductHandler) GetAllCategories(ctx *gin.Context) {
	result, err := h.productService.GetAllCategories()
	if err != nil {
		response.InternalServerError(ctx, "Failed to get categories", err.Error())
		return
	}

	response.OK(ctx, "Categories retrieved successfully", result)
}

// GetCategory godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /categories/{id} [get]
func (h *ProductHandler) GetCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid category ID", nil)
		return
	}

	result, err := h.productService.GetCategory(uint(id))
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(ctx, "Category not found")
			return
		}
		response.InternalServerError(ctx, "Failed to get category", err.Error())
		return
	}

	response.OK(ctx, "Category retrieved successfully", result)
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update a product category (Admin only)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Param        request body dto.UpdateCategoryRequest true "Update category request"
// @Success      200 {object} response.APIResponse{data=dto.CategoryResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /categories/{id} [put]
func (h *ProductHandler) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid category ID", nil)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := h.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(ctx, "Category not found")
			return
		}
		response.InternalServerError(ctx, "Failed to update category", err.Error())
		return
	}

	response.OK(ctx, "Category updated successfully", result)
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Delete a product category (Admin only)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} response.APIResponse
// @Failure      400 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /categories/{id} [delete]
func (h *ProductHandler) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "Invalid category ID", nil)
		return
	}

	if err := h.productService.DeleteCategory(uint(id)); err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(ctx, "Category not found")
			return
		}
		response.InternalServerError(ctx, "Failed to delete category", err.Error())
		return
	}

	response.OK(ctx, "Category deleted successfully", nil)
}
