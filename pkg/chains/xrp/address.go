package xrp

import "github.com/LanfordCai/ava/pkg/ripple"

// ValidateAddress validates an XRP address.
func ValidateAddress(addr string) bool {
	isValid, _ := ripple.New().ValidateAddress(addr, false)
	return isValid
}
