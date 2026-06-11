package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// GetPayment returns a payment (with client) or a 404 error.
func (s *Service) GetPayment(id uint) (*models.Payment, error) {
	var p models.Payment
	if err := s.DB.Preload("Client").First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(fmt.Sprintf("Payment #%d не найден", id))
		}
		return nil, err
	}
	return &p, nil
}

// ListPayments returns payments filtered by client and rental.
func (s *Service) ListPayments(clientID, rentalID uint, page, perPage, offset int) (*ListResult, error) {
	apply := func(db *gorm.DB) *gorm.DB {
		db = db.Model(&models.Payment{})
		if clientID > 0 {
			db = db.Where("client_id = ?", clientID)
		}
		if rentalID > 0 {
			db = db.Where("rental_id = ?", rentalID)
		}
		return db
	}
	var total int64
	apply(s.DB).Count(&total)

	var items []models.Payment
	if err := apply(s.DB).Preload("Client").Order("id desc").Limit(perPage).Offset(offset).Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for i := range items {
		out = append(out, items[i].ToMap())
	}
	return &ListResult{Items: out, Total: total, Page: page, PerPage: perPage}, nil
}

func (s *Service) generateReceiptNumber() string {
	return fmt.Sprintf("KV-%d-%05d", time.Now().Year(), s.nextSequence(&models.Payment{}))
}

// PaymentInput holds fields to create a payment.
type PaymentInput struct {
	ClientID    uint    `json:"client_id"`
	RentalID    *uint   `json:"rental_id"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	PaymentType string  `json:"payment_type"`
	Notes       *string `json:"notes"`
}

// CreatePayment records a money transaction and generates a receipt PDF.
func (s *Service) CreatePayment(in PaymentInput, actor *Actor) (*models.Payment, error) {
	client, err := s.GetClient(in.ClientID)
	if err != nil {
		return nil, err
	}
	if in.Amount <= 0 {
		return nil, apperr.Validation("Сумма платежа должна быть больше нуля")
	}
	if in.RentalID != nil {
		if _, err := s.GetRental(*in.RentalID); err != nil {
			return nil, err
		}
	}

	method := models.PaymentMethod(in.Method)
	if in.Method == "" {
		method = models.PaymentCash
	}
	ptype := models.PaymentType(in.PaymentType)
	if in.PaymentType == "" {
		ptype = models.PaymentTypeRental
	}
	if !method.Valid() || !ptype.Valid() {
		return nil, apperr.Validation("Некорректный способ или тип оплаты")
	}

	payment := &models.Payment{
		ReceiptNumber: s.generateReceiptNumber(),
		ClientID:      client.ID,
		RentalID:      in.RentalID,
		CashierID:     actor.idOrNil(),
		Amount:        in.Amount,
		Method:        method,
		PaymentType:   ptype,
		PaidAt:        time.Now(),
		Notes:         in.Notes,
	}
	if err := s.DB.Create(payment).Error; err != nil {
		return nil, err
	}

	// Generate the receipt PDF (best effort).
	payment.Client = client
	settings, _ := s.GetSettings()
	if path, err := s.Media.GenerateReceiptPDF(payment, settings); err == nil {
		s.DB.Model(payment).Update("receipt_path", path)
		payment.ReceiptPath = &path
	} else {
		log.Printf("payment: receipt PDF generation failed: %v", err)
	}

	s.record(actor, "create_payment", "payment", uintPtr(payment.ID),
		fmt.Sprintf("Платёж %s на сумму %.2f", payment.ReceiptNumber, in.Amount))
	return payment, nil
}
