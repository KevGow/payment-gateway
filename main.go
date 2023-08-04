package main

import (
	"log"
	"payment-gateway/api"
	"payment-gateway/bank"
	_ "payment-gateway/docs" // Needed for serving generated swagger docs
	"payment-gateway/payments"

	"github.com/gin-gonic/gin"                 // Gin framework
	swaggerFiles "github.com/swaggo/files"     // Swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // Gin-swagger middleware
)

// Note: Tagging above the main fucntion and handlers are for autogeneration of Swagger documents

// @title Payment Gateway API
// @description This is a simple Payment Gateway API for the ProcessOut take-home technical assessment.
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	// Create a new instance of PaymentGatewayService
	payments := payments.NewPaymentGatewayService()

	// Assign the Bank implementation to the PaymentGatewayService
	payments.Banker = new(bank.Bank)

	// Set up the router
	r := setupRouter(payments)
	// Serve Swagger UI at /swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Start the server on port 8080. This can be modified to read from an envar/config file/etc
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Could not run server with an error of: %v\n", err)
	}
}

// Function to set up the router and routes
func setupRouter(p *payments.PaymentGatewayService) *gin.Engine {
	// Create a new Gin router with default middleware
	router := gin.Default()

	// Define routes and their corresponding handler functions
	router.GET("/findpayment/:uuid", func(c *gin.Context) {
		// Handle GET requests for finding a payment
		api.HandleGetPayment(c, p)
	})
	router.POST("/pay", func(c *gin.Context) {
		// Handle POST requests for making a payment
		api.HandlePostPayment(c, p)
	})
	// Return the configured router
	return router
}
