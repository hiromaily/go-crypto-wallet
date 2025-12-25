package watch

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/application/ports/persistence"
)

// Type aliases for backward compatibility.
// These interfaces have been moved to pkg/application/ports/persistence.

// AddressRepositorier is AddressRepository interface
type AddressRepositorier = persistence.AddressRepositorier

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier = persistence.BTCTxRepositorier

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier = persistence.TxInputRepositorier

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier = persistence.TxOutputRepositorier

// TxRepositorier is TxRepository interface
type TxRepositorier = persistence.TxRepositorier

// PaymentRequestRepositorier is PaymentRequestRepository interface
type PaymentRequestRepositorier = persistence.PaymentRequestRepositorier

// EthDetailTxRepositorier is EthDetailTxRepository interface
type EthDetailTxRepositorier = persistence.EthDetailTxRepositorier

// XrpDetailTxRepositorier is XrpDetailTxRepository interface
type XrpDetailTxRepositorier = persistence.XrpDetailTxRepositorier
