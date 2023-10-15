package entity

import (
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
)

type Order struct {
	Number uint64
}

// Luhn algorithm
func (o *Order) IsValidNumber() error {
	if err := goluhn.Validate(strconv.FormatUint(o.Number, 10)); err != nil {
		return err
	}
	return nil
}
