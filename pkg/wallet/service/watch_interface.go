package service

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch"
)

//-----------------------------------------------------------------------------
// Watch Wallet Services - Type aliases for backward compatibility
// These interfaces have been moved to pkg/wallet/service/watch/interfaces.go
//-----------------------------------------------------------------------------

// AddressImporter is AddressImporter interface (for now btc/bch only)
type AddressImporter = watch.AddressImporter

// PaymentRequestCreator is PaymentRequestCreate interface
type PaymentRequestCreator = watch.PaymentRequestCreator

// TxCreator is TxCreator interface (for now btc/bch only)
type TxCreator = watch.TxCreator

// TxMonitorer is TxMonitor interface
type TxMonitorer = watch.TxMonitorer

// TxSender is TxSender interface
type TxSender = watch.TxSender
