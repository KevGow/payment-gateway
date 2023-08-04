package data

import (
	"sync"

	"github.com/google/uuid"
)

// GatewayData holds the payment data, which is held in an in-memory map,
// As well as a mutex lock to protect the PyamentData map, as it is shared.
// This can be expanded to contain other meta data about the system/payments.
type GatewayData struct {
	PaymentData map[PaymentID]Payment // A map that associates PaymentID with Payment.
	mu          sync.Mutex            // Mutex to protect concurrent access to PaymentData
}

// Payment represents a payment transaction.
type Payment struct {
	BankTransactionData // Embedding BankTransactionData to inherit its fields.
	CardData            // Embedding CardData to inherit its fields.
}

// BankTransactionData represents data related to a bank transaction.
type BankTransactionData struct {
	BankPaymentID     // Embedding BankPaymentID to inherit its fields.
	BankPaymentStatus // Embedding BankPaymentStatus to inherit its fields.
}

// CardData represents data related to the payment card.
type CardData struct {
	CardNumber string
	ExpiryDate string
	Amount     float64
	Currency   string
	Cvv        string
}

// PaymentID is a custom type representing a unique identifier for a payment.
type PaymentID uuid.UUID

// BankPaymentID is a custom type representing a unique identifier for a bank payment transaction.
type BankPaymentID uuid.UUID

// BankPaymentStatus is a custom type representing the status of a bank payment transaction.
type BankPaymentStatus string

func (g *GatewayData) AddPayment(bstatus BankPaymentStatus, bpid BankPaymentID, paymentId PaymentID, cd CardData) {
	// Lock the mutex to protect concurrent access to PaymentData
	g.mu.Lock()
	defer g.mu.Unlock()
	// Create a new BankTransactionData with the received bank payment status and ID
	var btd BankTransactionData
	btd.BankPaymentStatus = bstatus
	btd.BankPaymentID = bpid

	// Create a new Payment with the card data and bank transaction data
	var payment Payment
	payment.CardData = cd
	payment.BankTransactionData = btd

	// Add the payment to the PaymentData map with the generated payment ID
	g.PaymentData[paymentId] = payment
}

func (g *GatewayData) RetrievePayment(paymentId PaymentID) (bool, BankPaymentStatus, CardData) {
	// Lock the mutex to protect concurrent access to PaymentData
	g.mu.Lock()
	defer g.mu.Unlock()
	var maskedCardData CardData
	var bankPaymentStatus BankPaymentStatus
	if payment, ok := g.PaymentData[paymentId]; ok {
		// If found, return true, masked card number and data, and bank payment status
		maskedCardData.Amount = payment.Amount
		maskedCardData.CardNumber = MaskCardNumber(payment.CardData)
		maskedCardData.Currency = payment.Currency
		maskedCardData.ExpiryDate = payment.ExpiryDate
		bankPaymentStatus = payment.BankPaymentStatus
		return true, bankPaymentStatus, maskedCardData
	}
	// If payment not found, return false and empty strings for masked card number and bank payment status
	return false, bankPaymentStatus, maskedCardData
}

// MaskCardNumber masks the card number, keeping only the last four digits visible.
func MaskCardNumber(cd CardData) string {
	return "****" + cd.CardNumber[len(cd.CardNumber)-4:]
}
