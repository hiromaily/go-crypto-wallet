package xrp

import (
	"strconv"
)

// ToFloat64 converts string to float64
func ToFloat64(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0
	}
	return f
}
