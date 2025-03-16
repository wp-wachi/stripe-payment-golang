package services

import (
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/wp-wachi/stripe-payment-golang/config"
	"github.com/wp-wachi/stripe-payment-golang/models"
)

// CreatePaymentIntent calls Stripe API to create a payment
func CreatePaymentIntent(amount int64, currency string) (string, error) {
	stripe.Key = config.GetEnv("STRIPE_SECRET_KEY")

	// Create PaymentIntent
	params := &stripe.PaymentIntentParams{
		PaymentMethodTypes: []*string{stripe.String("promptpay")},
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Println("Stripe error:", err)
		return "", err
	}

	// Save payment details in the database
	payment := models.Payment{
		Amount:   amount,
		Currency: currency,
		Status:   "pending",
		StripeID: pi.ID,
	}

	// Print payment details to the log
    log.Printf("Payment details: %+v\n", payment)
    

	// repositories.SavePayment(payment)

	return pi.ClientSecret, nil
}