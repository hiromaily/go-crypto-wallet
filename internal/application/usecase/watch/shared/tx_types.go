package shared

import (
	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
)

// ParsedTx holds parsed UTXO data for transaction creation.
// This struct is shared between BTC and BCH use cases.
type ParsedTx struct {
	TxInputs       []dtobtc.TransactionInput
	TxRepoTxInputs []*domainBTC.BTCTxInput
	PrevTxs        []dtobtc.PreviousTx
	Addresses      []string // input, sender's address
}

// UserPayment represents user's payment address and amount.
// This struct is shared between BTC and BCH use cases for payment transactions.
type UserPayment struct {
	SenderAddr   string          // sender address for just checking
	ReceiverAddr string          // receiver address
	ValidRecAddr btcutil.Address // decoded receiver address
	Amount       float64         // amount
	ValidAmount  btcutil.Amount  // decoded amount
}
