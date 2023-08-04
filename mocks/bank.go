package mocks

import (
	"payment-gateway/bank"
	"payment-gateway/data"

	"github.com/google/uuid"
)

// BankMock is a mock implementation of the bank.Banker interface.
type BankMock struct {
	bank.Banker // Embedding the bank.Banker interface to satisfy the interface contract.
}

// MakePaymentToBank is the mocked version of the bank.Banker's MakePaymentToBank function.
// This function simulates making a payment to the bank and returns a predefined failure status and payment ID.
func (b *BankMock) MakePaymentToBank(cd data.CardData) (data.BankPaymentStatus, data.BankPaymentID) {
	bankPaymentStatus := data.BankPaymentStatus("Failure")
	bankPaymentId := data.BankPaymentID(uuid.New())
	return bankPaymentStatus, bankPaymentId
}
