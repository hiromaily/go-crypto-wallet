package shared

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/quagmt/udecimal"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// AmountConverter provides methods for amount conversion.
// This interface allows decoupling from the full Bitcoiner interface.
type AmountConverter interface {
	FloatToDecimal(f float64) (udecimal.Decimal, error)
}

// ParseListUnspentTx parses the result of listUnspent and creates transaction inputs.
// This function is shared between BTC and BCH use cases.
//
// Parameters:
//   - unspentList: List of unspent outputs from listunspent RPC
//   - amount: Required amount (0 means collect all UTXOs)
//   - converter: Amount converter for decimal conversion
//
// Returns:
//   - parsedTx: Parsed transaction data
//   - inputTotal: Total input amount
//   - isDone: Whether enough UTXOs were collected to meet the required amount
func ParseListUnspentTx(
	unspentList []dtobtc.UnspentOutput,
	amount btcutil.Amount,
	converter AmountConverter,
) (*ParsedTx, btcutil.Amount, bool) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]*domainBTC.BTCTxInput, 0, len(unspentList))
	prevTxs := make([]dtobtc.PreviousTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	var isDone bool // if isDone is false, sender can't meet amount
	if amount == 0 {
		isDone = true
	}

	for _, txItem := range unspentList {
		// Amount (already btcutil.Amount, no conversion needed)
		inputTotal += txItem.Amount

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: txItem.TxID,
			Vout: txItem.Vout,
		})

		// Convert btcutil.Amount to float64 for decimal conversion
		inputAmount, err := converter.FloatToDecimal(txItem.Amount.ToBTC())
		if err != nil {
			logger.Error("fail to convert input amount to decimal", "error", err)
			continue
		}
		input, err := domainBTC.NewBTCTxInput(
			0, // TxID will be set later
			txItem.TxID,
			txItem.Vout,
			txItem.Address,
			txItem.Label,
			inputAmount.String(),
			uint64(txItem.Confirmations),
		)
		if err != nil {
			logger.Error("fail to create BtcTxInput", "error", err)
			continue
		}
		txRepoTxInputs = append(txRepoTxInputs, input)

		// TODO: if sender is client account (non-multisig address), RedeemScript is blank
		prevTxs = append(prevTxs, dtobtc.PreviousTx{
			TxID:          txItem.TxID,
			Vout:          txItem.Vout,
			ScriptPubKey:  txItem.ScriptPubKey,
			RedeemScript:  txItem.RedeemScript,  // required if target account is multisig address
			WitnessScript: txItem.WitnessScript, // for SegWit addresses
			Amount:        txItem.Amount,        // Keep as btcutil.Amount
		})

		addresses = append(addresses, txItem.Address)

		// check total if amount is set as parameter
		if amount == 0 {
			continue
		}
		if inputTotal > amount {
			isDone = true
			break
		}
	}

	return &ParsedTx{
		TxInputs:       txInputs,
		TxRepoTxInputs: txRepoTxInputs,
		PrevTxs:        prevTxs,
		Addresses:      addresses,
	}, inputTotal, isDone
}
