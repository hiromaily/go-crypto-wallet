package decimal

import (
	"github.com/quagmt/udecimal"
)

// FloatToDecimal converts float64 to udecimal.Decimal.
// This is a simple wrapper around udecimal.NewFromFloat64 for consistency
// across the codebase.
func FloatToDecimal(f float64) (udecimal.Decimal, error) {
	return udecimal.NewFromFloat64(f)
}

