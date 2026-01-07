package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/akbarwjyy/go-commerce-api/internal/order/service"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/dto"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/entity"
	"github.com/akbarwjyy/go-commerce-api/internal/payment/repository"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderNotPending       = errors.New("order is not in pending status")
	ErrPaymentAlreadyExists  = errors.New("payment already exists for this order")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrUnauthorized          = errors.New("you are not authorized to perform this action")
	ErrPaymentAlreadyProcessed = errors.New("payment has already been processed")
)

// PaymentService interface untuk business logic payment
type PaymentService interface {
	CreatePayment(userID uint, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetPayment(userID uint, paymentID uint) (*dto.PaymentResponse, error)
	GetPaymentByOrderID(orderID uint) (*dto.PaymentResponse, error)
	GetMyPayments(userID uint, params *dto.PaymentQueryParams) (*dto.PaymentListResponse, error)
	GetAllPayments(params *dto.PaymentQueryParams) (*dto.PaymentListResponse, error)

	// Untuk callback simulasi
	ProcessPaymentCallback(transactionID string, status string, failedReason string) error
}

// paymentService implementasi PaymentService
type paymentService struct {
	paymentRepo  repository.PaymentRepository
	orderService service.OrderService
	db           *gorm.DB
}

// NewPaymentService membuat instance baru PaymentService
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderSvc service.OrderService,
	db *gorm.DB,
) PaymentService {
	return &paymentService{
		paymentRepo:  paymentRepo,
		orderService: orderSvc,
		db:           db,
	}
}

// CreatePayment membuat payment baru dan memulai proses async
func (s *paymentService) CreatePayment(userID uint, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Validate payment method
	if !entity.IsValidMethod(req.Method) {
		return nil, ErrInvalidPaymentMethod
	}

	// Check if payment already exists for this order
	existingPayment, _ := s.paymentRepo.FindByOrderID(req.OrderID)
	if existingPayment != nil && !existingPayment.IsFailed() {
		return nil, ErrPaymentAlreadyExists
	}

	// Get order details
	order, err := s.orderService.GetOrder(userID, req.OrderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Validate order status (must be PENDING)
	if order.Status != "PENDING" {
		return nil, ErrOrderNotPending
	}

	// Generate transaction ID
	transactionID := generateTransactionID()

	// Create payment record
	payment := &entity.Payment{
		OrderID:       req.OrderID,
		UserID:        userID,
		Amount:        order.TotalAmount,
		Method:        req.Method,
		Status:        entity.PaymentStatusPending,
		TransactionID: transactionID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	// Start async payment processing (Goroutine)
	go s.processPaymentAsync(payment.ID, transactionID)

	return s.toPaymentResponse(payment), nil
}

// processPaymentAsync memproses payment secara async dengan Goroutine
func (s *paymentService) processPaymentAsync(paymentID uint, transactionID string) {
	log.Printf("[Payment] Starting async processing for transaction: %s", transactionID)

	// Update status to PROCESSING
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		log.Printf("[Payment] Error finding payment: %v", err)
		return
	}
	payment.MarkAsProcessing()
	s.paymentRepo.Update(payment)

	// Simulate payment gateway delay (2-5 seconds)
	delay := time.Duration(2+rand.Intn(4)) * time.Second
	log.Printf("[Payment] Processing payment %s, waiting %v...", transactionID, delay)
	time.Sleep(delay)

	// Simulate success/failure (90% success rate)
	isSuccess := rand.Float32() < 0.9

	if isSuccess {
		// Mark payment as success
		payment.MarkAsSuccess()
		if err := s.paymentRepo.Update(payment); err != nil {
			log.Printf("[Payment] Error updating payment status: %v", err)
			return
		}

		// Callback to Order Module - Mark order as PAID
		if err := s.orderService.MarkAsPaid(payment.OrderID); err != nil {
			log.Printf("[Payment] Error marking order as paid: %v", err)
			return
		}

		log.Printf("[Payment] Payment %s SUCCESS! Order %d marked as PAID", transactionID, payment.OrderID)
	} else {
		// Mark payment as failed
		payment.MarkAsFailed("Payment declined by gateway (simulated)")
		if err := s.paymentRepo.Update(payment); err != nil {
			log.Printf("[Payment] Error updating payment status: %v", err)
			return
		}

		log.Printf("[Payment] Payment %s FAILED!", transactionID)
	}
}

// GetPayment mengambil payment berdasarkan ID
func (s *paymentService) GetPayment(userID uint, paymentID uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	// Check ownership
	if payment.UserID != userID {
		return nil, ErrUnauthorized
	}

	return s.toPaymentResponse(payment), nil
}

// GetPaymentByOrderID mengambil payment berdasarkan Order ID
func (s *paymentService) GetPaymentByOrderID(orderID uint) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return s.toPaymentResponse(payment), nil
}

// GetMyPayments mengambil payment milik user
func (s *paymentService) GetMyPayments(userID uint, params *dto.PaymentQueryParams) (*dto.PaymentListResponse, error) {
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

	payments, total, err := s.paymentRepo.FindByUserID(userID, params)
	if err != nil {
		return nil, err
	}

	var paymentResponses []dto.PaymentResponse
	for _, p := range payments {
		paymentResponses = append(paymentResponses, *s.toPaymentResponse(&p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.PaymentListResponse{
		Payments:   paymentResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetAllPayments mengambil semua payment (untuk admin)
func (s *paymentService) GetAllPayments(params *dto.PaymentQueryParams) (*dto.PaymentListResponse, error) {
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

	payments, total, err := s.paymentRepo.FindAll(params)
	if err != nil {
		return nil, err
	}

	var paymentResponses []dto.PaymentResponse
	for _, p := range payments {
		paymentResponses = append(paymentResponses, *s.toPaymentResponse(&p))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.PaymentListResponse{
		Payments:   paymentResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// ProcessPaymentCallback memproses callback dari payment gateway (untuk manual testing)
func (s *paymentService) ProcessPaymentCallback(transactionID string, status string, failedReason string) error {
	payment, err := s.paymentRepo.FindByTransactionID(transactionID)
	if err != nil {
		return ErrPaymentNotFound
	}

	// Only process if payment is still pending or processing
	if payment.IsSuccess() || payment.IsFailed() {
		return ErrPaymentAlreadyProcessed
	}

	if status == "SUCCESS" {
		payment.MarkAsSuccess()
		if err := s.paymentRepo.Update(payment); err != nil {
			return err
		}
		return s.orderService.MarkAsPaid(payment.OrderID)
	} else {
		payment.MarkAsFailed(failedReason)
		return s.paymentRepo.Update(payment)
	}
}

// Helper Functions

func (s *paymentService) toPaymentResponse(p *entity.Payment) *dto.PaymentResponse {
	resp := &dto.PaymentResponse{
		ID:            p.ID,
		OrderID:       p.OrderID,
		UserID:        p.UserID,
		Amount:        p.Amount,
		Method:        p.Method,
		Status:        p.Status,
		TransactionID: p.TransactionID,
		FailedReason:  p.FailedReason,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
	}

	if p.PaidAt != nil {
		resp.PaidAt = p.PaidAt.Format(time.RFC3339)
	}

	return resp
}

// generateTransactionID membuat transaction ID unik
func generateTransactionID() string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	return fmt.Sprintf("TXN-%d-%04d", timestamp, random)
}
