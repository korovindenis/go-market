package entity

import (
	"strconv"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
)

type Order struct {
	Number     uint64    `json:"number"`
	Status     Status    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Luhn algorithm
func (o *Order) IsValidNumber() error {
	if err := goluhn.Validate(strconv.FormatUint(o.Number, 10)); err != nil {
		return err
	}
	return nil
}

type Status uint8

const (
	StatusNew Status = iota
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func (s Status) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}[s]
}
