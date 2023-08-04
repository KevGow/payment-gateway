package bank

import (
	"payment-gateway/data"

	"github.com/google/uuid"
)

// Bank represents a concrete implementation of the Banker interface.
type Bank struct {
	Banker // Embedding the Banker interface to satisfy the interface contract.
}

// Banker is the interface that defines the contract for a bank service.
type Banker interface {
	MakePaymentToBank(cd data.CardData) (data.BankPaymentStatus, data.BankPaymentID)
}

// MakePaymentToBank simulates making a payment to the bank and receiving a response.
// We get back a resonse message, as well as uuid for refernce, This Uuid is NOT the
// Same as the payment uuid, and is simply a reference for the bank transaction
func (b *Bank) MakePaymentToBank(cd data.CardData) (data.BankPaymentStatus, data.BankPaymentID) {
	bankPaymentStatus := data.BankPaymentStatus("Success")
	bankPaymentId := data.BankPaymentID(uuid.New())
	return bankPaymentStatus, bankPaymentId
}
