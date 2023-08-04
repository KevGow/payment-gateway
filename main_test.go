package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"payment-gateway/api"
	"payment-gateway/bank"
	"payment-gateway/data"
	"payment-gateway/mocks"
	"payment-gateway/payments"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandlePostPayment tests the payment creation endpoint with valid payment data.
func TestHandlePostPayment(t *testing.T) {
	// Create a new PaymentGatewayService and set up the router.
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)
	router := setupRouter(p)

	// Create valid payment data in the request body.
	var cd api.PostJsonRequest
	cd.CardNumber = "4658585018481009"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/24"
	cd.Cvv = "555"

	// Marshal the payment data to JSON.
	jsonData, err := json.Marshal(cd)
	require.NoError(t, err)

	// Create a new HTTP request for the POST endpoint and record the response.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	// Assert that the response status code is 200 (OK).
	assert.Equal(t, 200, w.Code)

	// Unmarshal the response body into a PostResponse struct.
	var resp api.PostResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
}

// TestHandlePostPayment tests the payment creation endpoint with valid payment data
// and simulates concurrent requests.
func TestHandlePostPaymentConcurrent(t *testing.T) {
	// Create a new PaymentGatewayService and set up the router.
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)
	router := setupRouter(p)

	// Number of concurrent requests to simulate.
	numConcurrentRequests := 10

	// Create a WaitGroup to synchronize the completion of all goroutines.
	var wg sync.WaitGroup
	wg.Add(numConcurrentRequests)

	for i := 0; i < numConcurrentRequests; i++ {
		// Create valid payment data in the request body.
		var cd api.PostJsonRequest
		cd.CardNumber = "4658585018481009"
		cd.Amount = 100.00
		cd.Currency = "GBP"
		cd.ExpiryDate = "11/24"
		cd.Cvv = "555"

		// Marshal the payment data to JSON.
		jsonData, err := json.Marshal(cd)
		require.NoError(t, err)

		// Create a new HTTP request for the POST endpoint and record the response.
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))

		// Concurrently send the request using a goroutine.
		go func() {
			defer wg.Done()
			router.ServeHTTP(w, req)
		}()
	}

	// Wait for all goroutines to complete.
	wg.Wait()
}

func TestHandlePostPaymentWithIncorrectBody(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)
	router := setupRouter(p)

	// Using the get response struct out of convenience, the point is that pit's the wrong json body
	var getResp api.GetResponse
	getResp.MaskCardNumber = ""

	jsonData, err := json.Marshal(getResp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.PostResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Invalid json body"}`, w.Body.String())
}

func TestHandlePostPaymentWithInvalidCardNo(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	router := setupRouter(p)

	var cd api.PostJsonRequest
	cd.CardNumber = "4658585018481009123"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/24"
	cd.Cvv = "555"

	jsonData, err := json.Marshal(cd)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Invalid card number"}`, w.Body.String())
}

func TestHandlePostPaymentWithInvalidAmount(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	router := setupRouter(p)

	var cd api.PostJsonRequest
	cd.CardNumber = "4658585018481009"
	cd.Amount = -100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/24"
	cd.Cvv = "555"

	jsonData, err := json.Marshal(cd)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Invalid payment amount"}`, w.Body.String())
}

func TestHandlePostPaymentWithInvalidExpiry(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	router := setupRouter(p)

	var cd api.PostJsonRequest
	cd.CardNumber = "4658585018481009"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/09"
	cd.Cvv = "555"

	jsonData, err := json.Marshal(cd)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Card has expired"}`, w.Body.String())
}

func TestHandlePostPaymentWithInvalidCvv(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	router := setupRouter(p)

	var cd api.PostJsonRequest
	cd.CardNumber = "4658585018481009"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/24"
	cd.Cvv = "55555"

	jsonData, err := json.Marshal(cd)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pay", bytes.NewBuffer(jsonData))

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.PostResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Invalid CVV"}`, w.Body.String())
}

func TestHandleGetPayment(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	var cd data.CardData
	cd.CardNumber = "4658585018481009"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/22"
	cd.Cvv = "555"

	// Adding the payment to the in memory data store.
	pId := p.MakePayment(cd)
	router := setupRouter(p)

	w := httptest.NewRecorder()
	uuidValue := uuid.UUID(pId)
	strPaymentID := uuidValue.String()
	req, _ := http.NewRequest("GET", "/findpayment/"+strPaymentID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp api.GetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"amount":100, "bank-payment-status":"Success", "card-number-masked":"****1009", "currency":"GBP", "expiry-date":"11/22"}`, w.Body.String())
}

func TestHandleGetPaymentWithInvalidPaymentId(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)

	router := setupRouter(p)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/findpayment/InvalidID", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"Invalid uuid"}`, w.Body.String())
}

func TestHandleGetPaymentForNonExistantPayment(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	p.Banker = new(bank.Bank)
	router := setupRouter(p)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/findpayment/f2a5dd12-dad1-487c-ad60-d9f79a8aa6c6", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"error":"payment not found"}`, w.Body.String())
}

func TestHandleGetPaymentWithFailedBankPayment(t *testing.T) {
	p := payments.NewPaymentGatewayService()
	// Using the bank mock allows us to mock out different responses from the bank
	// In this case, we make a payment to the bank and receive an unsuccessful payment
	p.Banker = new(mocks.BankMock)

	var cd data.CardData
	cd.CardNumber = "4658585018481009"
	cd.Amount = 100.00
	cd.Currency = "GBP"
	cd.ExpiryDate = "11/22"
	cd.Cvv = "555"

	// Adding the payment to the in memory data store.
	pId := p.MakePayment(cd)

	router := setupRouter(p)

	w := httptest.NewRecorder()
	uuidValue := uuid.UUID(pId)
	strPaymentID := uuidValue.String()
	req, _ := http.NewRequest("GET", "/findpayment/"+strPaymentID, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp api.GetResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("ERROR : %v\n", err)
	}
	require.JSONEq(t, `{"amount":100, "bank-payment-status":"Failure", "card-number-masked":"****1009", "currency":"GBP", "expiry-date":"11/22"}`, w.Body.String())
}
