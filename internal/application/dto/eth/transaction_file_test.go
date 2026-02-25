package eth_test

import (
	"strings"
	"testing"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
)

func TestETHTransactionFileValidate_Legacy(t *testing.T) {
	tests := []struct {
		name    string
		txFile  dtoeth.ETHTransactionFile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid unsigned legacy transaction",
			txFile: dtoeth.ETHTransactionFile{
				Version:   1,
				TxType:    "unsigned",
				EthTxType: 0,
				ChainID:   1,
				Nonce:     42,
				From:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
				To:        "0x1234567890abcdef1234567890abcdef12345678",
				Value:     "1000000000000000000",
				Gas:       21000,
				GasPrice:  "20000000000",
				RawTxHex: "0xf86c2a8504a817c80082520894" +
					"1234567890abcdef1234567890abcdef1234567888016345785d8a0000801ba0...",
			},
			wantErr: false,
		},
		{
			name: "valid signed legacy transaction",
			txFile: dtoeth.ETHTransactionFile{
				Version:     1,
				TxType:      "signed",
				EthTxType:   0,
				ChainID:     1,
				Nonce:       42,
				From:        "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
				To:          "0x1234567890abcdef1234567890abcdef12345678",
				Value:       "1000000000000000000",
				Gas:         21000,
				GasPrice:    "20000000000",
				RawTxHex:    "0xf86c2a...",
				SignedTxHex: "0xf86c2a8504a817c8...",
			},
			wantErr: false,
		},
		{
			name: "legacy missing GasPrice",
			txFile: dtoeth.ETHTransactionFile{
				Version:   1,
				TxType:    "unsigned",
				EthTxType: 0,
				ChainID:   1,
				From:      "0xSender",
				To:        "0xRecipient",
				Value:     "1000",
				Gas:       21000,
				RawTxHex:  "0xabc",
			},
			wantErr: true,
			errMsg:  "gas_price",
		},
		{
			name: "signed missing SignedTxHex",
			txFile: dtoeth.ETHTransactionFile{
				Version:   1,
				TxType:    "signed",
				EthTxType: 0,
				ChainID:   1,
				From:      "0xSender",
				To:        "0xRecipient",
				Value:     "1000",
				Gas:       21000,
				GasPrice:  "20000000000",
				RawTxHex:  "0xabc",
			},
			wantErr: true,
			errMsg:  "signed_tx_hex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.txFile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestETHTransactionFileValidate_EIP1559(t *testing.T) {
	tests := []struct {
		name    string
		txFile  dtoeth.ETHTransactionFile
		wantErr bool
	}{
		{
			name: "valid unsigned EIP-1559 transaction",
			txFile: dtoeth.ETHTransactionFile{
				Version:              1,
				TxType:               "unsigned",
				EthTxType:            2,
				ChainID:              1,
				Nonce:                42,
				From:                 "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
				To:                   "0x1234567890abcdef1234567890abcdef12345678",
				Value:                "1000000000000000000",
				Gas:                  21000,
				MaxFeePerGas:         "30000000000",
				MaxPriorityFeePerGas: "2000000000",
				RawTxHex:             "0x02f86c...",
			},
			wantErr: false,
		},
		{
			name: "EIP-1559 missing MaxFeePerGas",
			txFile: dtoeth.ETHTransactionFile{
				Version:              1,
				TxType:               "unsigned",
				EthTxType:            2,
				ChainID:              1,
				From:                 "0xSender",
				To:                   "0xRecipient",
				Value:                "1000",
				Gas:                  21000,
				MaxPriorityFeePerGas: "2000000000",
				RawTxHex:             "0x02abc",
			},
			wantErr: true,
		},
		{
			name: "EIP-1559 missing MaxPriorityFeePerGas",
			txFile: dtoeth.ETHTransactionFile{
				Version:      1,
				TxType:       "unsigned",
				EthTxType:    2,
				ChainID:      1,
				From:         "0xSender",
				To:           "0xRecipient",
				Value:        "1000",
				Gas:          21000,
				MaxFeePerGas: "30000000000",
				RawTxHex:     "0x02abc",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.txFile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestETHTransactionFileValidate_CommonFields(t *testing.T) {
	validBase := dtoeth.ETHTransactionFile{
		Version:   1,
		TxType:    "unsigned",
		EthTxType: 0,
		ChainID:   1,
		From:      "0xSender",
		To:        "0xRecipient",
		Value:     "1000",
		Gas:       21000,
		GasPrice:  "20000000000",
		RawTxHex:  "0xabc",
	}

	tests := []struct {
		name    string
		modify  func(f *dtoeth.ETHTransactionFile)
		wantErr bool
	}{
		{
			name:    "valid base",
			modify:  func(_ *dtoeth.ETHTransactionFile) {},
			wantErr: false,
		},
		{
			name:    "invalid version (zero)",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.Version = 0 },
			wantErr: true,
		},
		{
			name:    "invalid TxType",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.TxType = "invalid" },
			wantErr: true,
		},
		{
			name:    "invalid EthTxType (1 not supported)",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.EthTxType = 1 },
			wantErr: true,
		},
		{
			name:    "invalid EthTxType (3 not supported)",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.EthTxType = 3 },
			wantErr: true,
		},
		{
			name:    "zero ChainID",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.ChainID = 0 },
			wantErr: true,
		},
		{
			name:    "empty From",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.From = "" },
			wantErr: true,
		},
		{
			name:    "empty To",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.To = "" },
			wantErr: true,
		},
		{
			name:    "empty Value",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.Value = "" },
			wantErr: true,
		},
		{
			name:    "zero Gas",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.Gas = 0 },
			wantErr: true,
		},
		{
			name:    "empty RawTxHex",
			modify:  func(f *dtoeth.ETHTransactionFile) { f.RawTxHex = "" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := validBase
			tt.modify(&f)
			err := f.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
