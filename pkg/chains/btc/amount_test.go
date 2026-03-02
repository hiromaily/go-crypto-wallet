package btc

import (
	"math"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAmountString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		amt  btcutil.Amount
		want string
	}{
		{
			name: "1 BTC",
			amt:  btcutil.Amount(100000000),
			want: "1",
		},
		{
			name: "0.54 BTC",
			amt:  btcutil.Amount(54000000),
			want: "0.54",
		},
		{
			name: "1 satoshi",
			amt:  btcutil.Amount(1),
			want: "0.00000001",
		},
		{
			name: "0 BTC",
			amt:  btcutil.Amount(0),
			want: "0",
		},
		{
			name: "10 BTC",
			amt:  btcutil.Amount(1000000000),
			want: "10",
		},
		{
			name: "0.001 BTC",
			amt:  btcutil.Amount(100000),
			want: "0.001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := AmountString(tt.amt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAmountToDecimal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		amt     btcutil.Amount
		want    string
		wantErr bool
	}{
		{
			name: "1 BTC",
			amt:  btcutil.Amount(100000000),
			want: "1",
		},
		{
			name: "0.54 BTC",
			amt:  btcutil.Amount(54000000),
			want: "0.54",
		},
		{
			name: "0 BTC",
			amt:  btcutil.Amount(0),
			want: "0",
		},
		{
			name: "0.001 BTC",
			amt:  btcutil.Amount(100000),
			want: "0.001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := AmountToDecimal(tt.amt)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestFloatToDecimal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   float64
		want    string
		wantErr bool
	}{
		{
			name:  "1.0",
			input: 1.0,
			want:  "1",
		},
		{
			name:  "0.54",
			input: 0.54,
			want:  "0.54",
		},
		{
			name:  "0.0",
			input: 0.0,
			want:  "0",
		},
		{
			name:  "10.5",
			input: 10.5,
			want:  "10.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FloatToDecimal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestFloatToAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   float64
		want    btcutil.Amount
		wantErr bool
	}{
		{
			name:  "1 BTC",
			input: 1.0,
			want:  btcutil.Amount(100000000),
		},
		{
			name:  "0.54 BTC",
			input: 0.54,
			want:  btcutil.Amount(54000000),
		},
		{
			name:  "0 BTC",
			input: 0.0,
			want:  btcutil.Amount(0),
		},
		{
			name:  "0.001 BTC",
			input: 0.001,
			want:  btcutil.Amount(100000),
		},
		{
			name:    "NaN",
			input:   math.NaN(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FloatToAmount(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStrToAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    btcutil.Amount
		wantErr bool
	}{
		{
			name:  "1 BTC as string",
			input: "1.0",
			want:  btcutil.Amount(100000000),
		},
		{
			name:  "0.54 BTC as string",
			input: "0.54",
			want:  btcutil.Amount(54000000),
		},
		{
			name:  "0 BTC as string",
			input: "0",
			want:  btcutil.Amount(0),
		},
		{
			name:  "0.001 BTC as string",
			input: "0.001",
			want:  btcutil.Amount(100000),
		},
		{
			name:    "invalid string",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := StrToAmount(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStrSatoshiToAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    btcutil.Amount
		wantErr bool
	}{
		{
			name:  "100000000 satoshi = 1 BTC",
			input: "100000000",
			want:  btcutil.Amount(100000000),
		},
		{
			name:  "54000000 satoshi = 0.54 BTC",
			input: "54000000",
			want:  btcutil.Amount(54000000),
		},
		{
			name:  "0 satoshi",
			input: "0",
			want:  btcutil.Amount(0),
		},
		{
			name:  "100000 satoshi = 0.001 BTC",
			input: "100000",
			want:  btcutil.Amount(100000),
		},
		{
			name:    "invalid string",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := StrSatoshiToAmount(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateNewFee(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fee           btcutil.Amount
		adjustmentFee float64
		want          btcutil.Amount
		wantErr       bool
	}{
		{
			name:          "increase fee by 1.5x",
			fee:           btcutil.Amount(100000), // 0.001 BTC
			adjustmentFee: 1.5,
			want:          btcutil.Amount(150000), // 0.0015 BTC
		},
		{
			name:          "no adjustment (1.0x)",
			fee:           btcutil.Amount(100000),
			adjustmentFee: 1.0,
			want:          btcutil.Amount(100000),
		},
		{
			name:          "double fee (2.0x)",
			fee:           btcutil.Amount(50000),
			adjustmentFee: 2.0,
			want:          btcutil.Amount(100000),
		},
		{
			name:          "zero fee",
			fee:           btcutil.Amount(0),
			adjustmentFee: 1.5,
			want:          btcutil.Amount(0),
		},
		{
			name:          "NaN adjustment",
			fee:           btcutil.Amount(100000),
			adjustmentFee: math.NaN(),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := CalculateNewFee(tt.fee, tt.adjustmentFee)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
