package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// Bitcoin includes client to call Json-RPC
type Bitcoin struct {
	client            *rpcclient.Client
	chainConf         *chaincfg.Params
	coinTypeCode      coin.CoinTypeCode //btc
	version           coin.BTCVersion   //179900
	confirmationBlock uint64
	feeRange          FeeAdjustmentRate
	logger            *zap.Logger
}

// FeeAdjustmentRate range of fee adjustment rate
type FeeAdjustmentRate struct {
	min float64
	max float64
}

// NewBitcoin creates bitcoin object
func NewBitcoin(
	client *rpcclient.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Bitcoin,
	logger *zap.Logger) (*Bitcoin, error) {

	bit := Bitcoin{
		client: client,
		logger: logger,
	}

	bit.coinTypeCode = coinTypeCode

	switch conf.NetworkType {
	case coin.NetworkTypeMainNet:
		bit.chainConf = &chaincfg.MainNetParams
	case coin.NetworkTypeTestNet3:
		bit.chainConf = &chaincfg.TestNet3Params
	case coin.NetworkTypeRegTestNet:
		bit.chainConf = &chaincfg.RegressionNetParams
	default:
		return nil, errors.Errorf("bitcoin network type is invalid in config")
	}

	// set bitcoin version
	netInfo, err := bit.GetNetworkInfo()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call bit.GetNetworkInfo()")
	}
	if coin.RequiredVersion > netInfo.Version {
		return nil, errors.Errorf("bitcoin core version should be %d +, but version %d is detected", coin.RequiredVersion, netInfo.Version)
	}
	bit.version = netInfo.Version
	bit.logger.Info("bitcoin rpc server", zap.Int("version", netInfo.Version.Int()))

	// set other information from config
	bit.confirmationBlock = conf.Block.ConfirmationNum
	bit.feeRange.max = conf.Fee.AdjustmentMax
	bit.feeRange.min = conf.Fee.AdjustmentMin

	return &bit, nil
}

// Close disconnect from bitcoin core server
func (b *Bitcoin) Close() {
	if b.client != nil {
		b.client.Shutdown()
	}
}

// GetChainConf returns chain conf
func (b *Bitcoin) GetChainConf() *chaincfg.Params {
	return b.chainConf
}

// SetChainConf sets chain conf
func (b *Bitcoin) SetChainConf(conf *chaincfg.Params) {
	b.chainConf = conf
}

// SetChainConfNet sets conf.Net
func (b *Bitcoin) SetChainConfNet(btcNet wire.BitcoinNet) {
	b.chainConf.Net = btcNet
}

// ConfirmationBlock returns confirmation block count
func (b *Bitcoin) ConfirmationBlock() uint64 {
	return b.confirmationBlock
}

// FeeRangeMax return maximum fee rate for adjustment
func (b *Bitcoin) FeeRangeMax() float64 {
	return b.feeRange.max
}

// FeeRangeMin returns minimum fee rate for adjustment
func (b *Bitcoin) FeeRangeMin() float64 {
	return b.feeRange.min
}

// Version returns core version
func (b *Bitcoin) Version() coin.BTCVersion {
	return b.version
}

// CoinTypeCode returns CoinTypeCode
func (b *Bitcoin) CoinTypeCode() coin.CoinTypeCode {
	return b.coinTypeCode
}
