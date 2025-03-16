package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wp-wachi/stripe-payment-golang/services"
	"github.com/wp-wachi/stripe-payment-golang/utils"
)

// CreatePaymentHandler handles payment requests
func CreatePaymentHandler(c *gin.Context) {
	var req struct {
		Amount   int64  `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, "Invalid request data")
		return
	}
	
	var currency string = "thb"
	clientSecret, err := services.CreatePaymentIntent(req.Amount, currency)
	if err != nil {
		utils.ErrorResponse(c, "Payment failed")
		return
	}

	utils.SuccessResponse(c, gin.H{"client_secret": clientSecret})
}