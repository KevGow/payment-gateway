package payments

import (
	"payment-gateway/bank"
	"payment-gateway/data"
	"payment-gateway/validation"

	"github.com/google/uuid"
)

// PaymentGatewayService represents the payment gateway service that handles payment operations.
type PaymentGatewayService struct {
	data.GatewayData // Embedding GatewayData to inherit its fields and methods
	bank.Banker      // Embedding Banker interface to use bank-related functionality
}

// NewPaymentGatewayService creates a new instance of PaymentGatewayService and initializes the PaymentData map.
func NewPaymentGatewayService() *PaymentGatewayService {
	p := new(PaymentGatewayService)
	// Payment data is our in-memory data store and can easily be ripped out
	p.GatewayData.PaymentData = make(map[data.PaymentID]data.Payment)
	return p
}

// GetPayment retrieves payment information based on the provided payment ID.
func (p *PaymentGatewayService) GetPayment(paymentId data.PaymentID) (bool, data.BankPaymentStatus, data.CardData) {
	// Check if the paymentId exists in the PaymentData map
	exists, bankPaymentStatus, maskedCardData := p.GatewayData.RetrievePayment(paymentId)

	// return the details of the payment
	return exists, bankPaymentStatus, maskedCardData
}

// MakePayment initiates a new payment transaction with the provided card data.
func (p *PaymentGatewayService) MakePayment(cd data.CardData) data.PaymentID {
	// Generate a payment id to record the payment
	paymentId := data.PaymentID(uuid.New())

	// Use the embedded Banker interface to make a payment to the bank
	// Note this also returns an UUID, which is our reference to the
	// Payment for the bank
	bstatus, bpid := p.Banker.MakePaymentToBank(cd)

	// Add the payment to the PaymentData
	p.GatewayData.AddPayment(bstatus, bpid, paymentId, cd)
	// returns the payment id to the client
	return paymentId
}

// ValidatePayment validates the card data before processing the payment.
func ValidatePayment(cd data.CardData) (bool, string) {
	// Validate card number using Luhn's algorithm
	if !validation.LuhnCheck(cd.CardNumber) {
		return false, "Invalid card number"
	}
	// Validate expiration date of the card
	if !validation.ValidateExpirationDate(cd.ExpiryDate) {
		return false, "Card has expired"
	}
	// Validate CVV number of the card
	if !validation.ValidateCVV(cd.Cvv) {
		return false, "Invalid CVV"
	}
	// Validate payment amount
	if !validation.ValidatePaymentAmount(cd.Amount) {
		return false, "Invalid payment amount"
	}
	// If all validations pass, return true and an empty error message
	return true, ""
}
