package entity

import (
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
)

// struct for user Balance
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// struct for update user Balance
type BalanceUpdate struct {
	Order      string    `json:"order"`
	Sum        float64   `json:"sum"`
	UploadedAt time.Time `json:"processed_at,omitempty"`
}

// Luhn algorithm
func (b *BalanceUpdate) IsValidNumber() error {
	if err := goluhn.Validate(b.Order); err != nil {
		return err
	}
	return nil
}
