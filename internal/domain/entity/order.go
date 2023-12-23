package entity

import (
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
)

// order status
const (
	StatusNew        = "NEW"
	StatusRegistered = "REGISTERED"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
)

// struct for user Order
type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"Status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Sum        float64   `json:"sum,omitempty"`
}

// Luhn algorithm
func (o *Order) IsValidNumber() error {
	if err := goluhn.Validate(o.Number); err != nil {
		return err
	}
	return nil
}
