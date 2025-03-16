package routes

import (
	"github.com/wp-wachi/stripe-payment-golang/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/create-payment-intent", controllers.CreatePaymentHandler)
	r.POST("/payment-intent-webhook", controllers.StripeWebhookHandler)
}