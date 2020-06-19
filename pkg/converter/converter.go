package converter

import (
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/types"
)

// Converter is type converter
type Converter interface {
	FloatToDecimal(f float64) types.Decimal
}

type convert struct{}

// NewConverter returns Converter interface
func NewConverter() Converter {
	return &convert{}
}

// FloatToDecimal converts float to decimal
func (c convert) FloatToDecimal(f float64) types.Decimal {
	strAmt := fmt.Sprintf("%f", f)
	dAmt := types.Decimal{Big: new(decimal.Big)}
	dAmt.Big, _ = dAmt.SetString(strAmt)
	return dAmt
}
