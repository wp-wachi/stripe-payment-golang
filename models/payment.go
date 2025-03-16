package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	StripeID    string `json:"stripe_id"`
}