package converter

import (
	"github.com/quagmt/udecimal"
)

// Converter is type converter
type Converter interface {
	FloatToDecimal(f float64) (udecimal.Decimal, error)
}

type convert struct{}

// NewConverter returns Converter interface
func NewConverter() Converter {
	return &convert{}
}

// FloatToDecimal converts float to decimal
func (convert) FloatToDecimal(f float64) (udecimal.Decimal, error) {
	return udecimal.NewFromFloat64(f)
}
