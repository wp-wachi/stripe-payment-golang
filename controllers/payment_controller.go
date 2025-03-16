package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/wp-wachi/stripe-payment-golang/config"
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

// StripeWebhookHandler listens for Stripe events
func StripeWebhookHandler(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorResponse(c, "Request body too large")
		return
	}

	// Get Stripe signature from headers
	signatureHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := config.GetEnv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEvent(body, signatureHeader, endpointSecret)
	if err != nil {
		utils.ErrorResponse(c, "Invalid webhook signature")
		return
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			log.Println("Error parsing event:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing event"})
			return
		}

		// Update payment status in DB
		// repositories.UpdatePaymentStatus(intent.ID, "succeeded")
		log.Println("✅ Payment successful:", intent.ID)

		// Send message to LINE
		message := fmt.Sprintf("✅ Payment Successful!\nAmount: %d %s\nTransaction ID: %s", intent.Amount/100, intent.Currency, intent.ID)
		utils.SendMessageToLINE(message)

	case "payment_intent.payment_failed":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			log.Println("Error parsing event:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing event"})
			return
		}

		// Update payment status in DB
		// repositories.UpdatePaymentStatus(intent.ID, "failed")
		log.Println("❌ Payment failed:", intent.ID)

		// Send failure message to LINE
		message := fmt.Sprintf("❌ Payment Failed!\nAmount: %d %s\nTransaction ID: %s", intent.Amount/100, intent.Currency, intent.ID)
		utils.SendMessageToLINE(message)

	default:
		log.Println("Unhandled event type:", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}