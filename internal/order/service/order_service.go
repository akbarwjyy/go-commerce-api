package service

import (
	"errors"
	"math"
	"time"

	"github.com/akbarwjyy/go-commerce-api/internal/order/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/order/entity"
	"github.com/akbarwjyy/go-commerce-api/internal/order/repository"
	productService "github.com/akbarwjyy/go-commerce-api/internal/product/service"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrUnauthorized       = errors.New("you are not authorized to perform this action")
	ErrInvalidStatus      = errors.New("invalid status transition")
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("insufficient stock for one or more products")
	ErrEmptyCart          = errors.New("cart is empty")
	ErrOrderNotCancellable = errors.New("order cannot be cancelled")
)

// OrderService interface untuk business logic order
type OrderService interface {
	Checkout(userID uint, req *dto.CheckoutRequest) (*dto.OrderResponse, error)
	GetOrder(userID uint, orderID uint) (*dto.OrderResponse, error)
	GetMyOrders(userID uint, params *dto.OrderQueryParams) (*dto.OrderListResponse, error)
	GetAllOrders(params *dto.OrderQueryParams) (*dto.OrderListResponse, error)
	UpdateOrderStatus(userID uint, orderID uint, status string, isAdmin bool) (*dto.OrderResponse, error)
	CancelOrder(userID uint, orderID uint) error

	// Untuk Payment Module callback
	MarkAsPaid(orderID uint) error
}

// orderService implementasi OrderService
type orderService struct {
	orderRepo      repository.OrderRepository
	productService productService.ProductService
	db             *gorm.DB
}

// NewOrderService membuat instance baru OrderService
func NewOrderService(
	orderRepo repository.OrderRepository,
	productSvc productService.ProductService,
	db *gorm.DB,
) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		productService: productSvc,
		db:             db,
	}
}

// Checkout membuat order baru dari checkout
func (s *orderService) Checkout(userID uint, req *dto.CheckoutRequest) (*dto.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, ErrEmptyCart
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var orderItems []entity.OrderItem
	var totalAmount float64

	// Validate and process each item
	for _, item := range req.Items {
		// Get product details
		product, err := s.productService.GetProductByID(item.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, ErrProductNotFound
		}

		// Check stock
		if !product.HasStock(item.Quantity) {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}

		// Reduce stock
		if err := s.productService.ReduceStock(item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create order item
		subtotal := product.Price * float64(item.Quantity)
		orderItem := entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += subtotal
	}

	// Create order
	order := &entity.Order{
		UserID:       userID,
		TotalAmount:  totalAmount,
		Status:       entity.OrderStatusPending,
		ShippingAddr: req.ShippingAddress,
		Notes:        req.Notes,
		Items:        orderItems,
	}

	orderRepoWithTx := s.orderRepo.WithTx(tx)
	if err := orderRepoWithTx.Create(order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload order with items
	order, _ = s.orderRepo.FindByIDWithItems(order.ID)

	return s.toOrderResponse(order), nil
}

// GetOrder mengambil order berdasarkan ID
func (s *orderService) GetOrder(userID uint, orderID uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByIDWithItems(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Check ownership
	if !order.IsOwner(userID) {
		return nil, ErrUnauthorized
	}

	return s.toOrderResponse(order), nil
}

// GetMyOrders mengambil order milik user
func (s *orderService) GetMyOrders(userID uint, params *dto.OrderQueryParams) (*dto.OrderListResponse, error) {
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

	orders, total, err := s.orderRepo.FindByUserID(userID, params)
	if err != nil {
		return nil, err
	}

	var orderResponses []dto.OrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, *s.toOrderResponse(&o))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetAllOrders mengambil semua order (untuk admin)
func (s *orderService) GetAllOrders(params *dto.OrderQueryParams) (*dto.OrderListResponse, error) {
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

	orders, total, err := s.orderRepo.FindAll(params)
	if err != nil {
		return nil, err
	}

	var orderResponses []dto.OrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, *s.toOrderResponse(&o))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus mengupdate status order
func (s *orderService) UpdateOrderStatus(userID uint, orderID uint, status string, isAdmin bool) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByIDWithItems(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// User hanya bisa update status tertentu (cancel)
	// Admin bisa update semua status
	if !isAdmin && !order.IsOwner(userID) {
		return nil, ErrUnauthorized
	}

	// Validate status transition
	if !order.UpdateStatus(status) {
		return nil, ErrInvalidStatus
	}

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// CancelOrder membatalkan order dan mengembalikan stok
func (s *orderService) CancelOrder(userID uint, orderID uint) error {
	order, err := s.orderRepo.FindByIDWithItems(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	// Check ownership
	if !order.IsOwner(userID) {
		return ErrUnauthorized
	}

	// Check if order can be cancelled
	if !order.CanBeCancelled() {
		return ErrOrderNotCancellable
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Restore stock for each item
	for _, item := range order.Items {
		if err := s.productService.RestoreStock(item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update status to cancelled
	order.Status = entity.OrderStatusCancelled
	orderRepoWithTx := s.orderRepo.WithTx(tx)
	if err := orderRepoWithTx.Update(order); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// MarkAsPaid dipanggil oleh Payment Module untuk update status
func (s *orderService) MarkAsPaid(orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if !order.IsPending() {
		return ErrInvalidStatus
	}

	order.Status = entity.OrderStatusPaid
	return s.orderRepo.Update(order)
}

// Helper Functions

func (s *orderService) toOrderResponse(o *entity.Order) *dto.OrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range o.Items {
		items = append(items, dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal,
		})
	}

	return &dto.OrderResponse{
		ID:              o.ID,
		UserID:          o.UserID,
		TotalAmount:     o.TotalAmount,
		Status:          o.Status,
		ShippingAddress: o.ShippingAddr,
		Notes:           o.Notes,
		Items:           items,
		CreatedAt:       o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       o.UpdatedAt.Format(time.RFC3339),
	}
}
