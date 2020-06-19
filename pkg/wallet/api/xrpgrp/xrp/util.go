package xrp

import (
	"strconv"

	"github.com/LanfordCai/ava/pkg/ripple"
)

// ToFloat64 converts string to float64
func ToFloat64(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0
	}
	return f
}

// XRPToDrops converts an XRP amount to drops. 1 XRP = 1,000,000 drops, so 1 drop = 0.000001 XRP
// - https://xrpl.org/rippleapi-reference.html#xrptodrops
// nolint:golint
func XRPToDrops(val float64) float64 {
	return val * 0.000001
}

// ValidateAddress validates address
func ValidateAddress(addr string) bool {
	isValid, _ := ripple.New().ValidateAddress(addr, false)
	return isValid
}
