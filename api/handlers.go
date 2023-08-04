package api

import (
	"net/http"
	"payment-gateway/data"
	"payment-gateway/payments"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Make a payment
// @Description Make a payment
// @ID make-payment
// @Accept json
// @Produce json
// @Param paymentData body PostJsonRequest true "Payment Data"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Router /pay [post]
func HandlePostPayment(c *gin.Context, p *payments.PaymentGatewayService) {
	// Bind the JSON data from the request body to the PostJsonRequest struct
	var body PostJsonRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid json body"})
		return
	}

	// Convert the PostJsonRequest data to a CardData struct
	cd := data.CardData{
		CardNumber: strings.ReplaceAll(body.CardNumber, " ", ""),
		ExpiryDate: body.ExpiryDate,
		Amount:     body.Amount,
		Currency:   body.Currency,
		Cvv:        body.Cvv,
	}

	// Validate the payment data using the ValidatePayment function
	isValid, message := payments.ValidatePayment(cd)
	if isValid {
		// If the payment data is valid, call the MakePayment method of the PaymentGatewayService
		paymentId := p.MakePayment(cd)
		// Respond with the generated UUID for the payment
		c.IndentedJSON(http.StatusOK, gin.H{"uuid": uuid.UUID(paymentId).String()})
		return
	}

	// If the payment data is invalid, respond with 400 status and error message
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": message})
}

// @Summary Get payment information by UUID
// @Description Get payment information by UUID
// @ID get-payment-by-uuid
// @Produce json
// @Param uuid path string true "Payment UUID"
// @Success 200 {object} GetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /findpayment/{uuid} [get]
func HandleGetPayment(c *gin.Context, p *payments.PaymentGatewayService) {
	// Parse the UUID parameter from the request URL
	u, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}

	// Convert the parsed UUID to a custom PaymentID type
	paymentId := data.PaymentID(u)

	// Call the GetPayment method of the PaymentGatewayService to retrieve payment information
	if ok, bankPaymentStatus, maskedCardData := p.GetPayment(paymentId); ok {
		c.IndentedJSON(http.StatusOK, gin.H{
			"bank-payment-status": bankPaymentStatus,
			"amount":              maskedCardData.Amount,
			"currency":            maskedCardData.Currency,
			"card-number-masked":  maskedCardData.CardNumber,
			"expiry-date":         maskedCardData.ExpiryDate,
		})
		return
	}

	// If payment not found, respond with 404 status and error message
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "payment not found"})
}

// swagger:model
type PostResponse struct {
	Uuid uuid.UUID `json:"uuid"`
}

// swagger:model
type GetResponse struct {
	BankPaymentStatus string  `json:"bank-payment-status" example:"Success"`
	Amount            float64 `json:"amount" example:"100.00"`
	Currency          string  `json:"currency" example:"GBP"`
	MaskCardNumber    string  `json:"card-number-masked" example:"****5070"`
	ExpiryDate        string  `json:"expiry-date" example:"11/26"`
}

// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}

// PostJsonRequest represents the JSON data expected in POST requests for making a payment.
type PostJsonRequest struct {
	CardNumber string  `json:"card-number" example:"4032 0341 3083 5070" binding:"required"`
	ExpiryDate string  `json:"expiry-date" example:"11/26" binding:"required"`
	Amount     float64 `json:"amount" example:"100.00" binding:"required"`
	Currency   string  `json:"currency" example:"GBP" binding:"required"`
	Cvv        string  `json:"cvv" example:"975" binding:"required"`
}
