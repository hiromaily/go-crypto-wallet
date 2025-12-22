package converter

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

// Converter is type converter
type Converter interface {
	FloatToDecimal(f float64) *decimal.Big
}

type convert struct{}

// NewConverter returns Converter interface
func NewConverter() Converter {
	return &convert{}
}

// FloatToDecimal converts float to decimal
func (convert) FloatToDecimal(f float64) *decimal.Big {
	strAmt := fmt.Sprintf("%f", f)
	dAmt := new(decimal.Big)
	dAmt, _ = dAmt.SetString(strAmt)
	return dAmt
}
