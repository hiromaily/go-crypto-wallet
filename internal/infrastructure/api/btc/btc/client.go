package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Bitcoin includes btcdClient to call Json-RPC
type Bitcoin struct {
	pkgrpc            btcrpc.BTCRPC
	btcdClient        *rpcclient.Client
	chainConf         *chaincfg.Params
	coinTypeCode      domainCoin.CoinTypeCode // btc
	version           btcpkg.NodeVersion      // 280000
	confirmationBlock uint64
	feeRange          FeeAdjustmentRate
}

// FeeAdjustmentRate range of fee adjustment rate
type FeeAdjustmentRate struct {
	min float64
	max float64
}

// NewBitcoin creates bitcoin object
func NewBitcoin(
	btcdClient *rpcclient.Client,
	conf *config.Bitcoin,
	coinTypeCode domainCoin.CoinTypeCode,
) (*Bitcoin, error) {
	bit := Bitcoin{
		pkgrpc:     btcrpc.NewRPCClient(btcdClient),
		btcdClient: btcdClient,
	}

	bit.coinTypeCode = coinTypeCode

	// check network consistency between config and bitcoind
	blockInfo, err := bit.pkgrpc.GetBlockchainInfo()
	if err != nil {
		return nil, fmt.Errorf("fail to call bit.GetBlockchainInfo(): %w", err)
	}

	switch btcpkg.NetworkType(conf.NetworkType) {
	case btcpkg.NetworkTypeMainNet:
		bit.chainConf = &chaincfg.MainNetParams
		if blockInfo.Chain != "main" {
			return nil, fmt.Errorf(
				"connecting %s on bitcoind, but config file defines as %s",
				blockInfo.Chain, btcpkg.NetworkTypeMainNet)
		}
	case btcpkg.NetworkTypeTestNet3:
		bit.chainConf = &chaincfg.TestNet3Params
		if blockInfo.Chain != "test" {
			return nil, fmt.Errorf(
				"connecting %s on bitcoind, but config file defines as %s",
				blockInfo.Chain, btcpkg.NetworkTypeTestNet3)
		}
	case btcpkg.NetworkTypeRegTest:
		bit.chainConf = &chaincfg.RegressionNetParams
		if blockInfo.Chain != "regtest" {
			return nil, fmt.Errorf(
				"connecting %s on bitcoind, but config file defines as %s",
				blockInfo.Chain, btcpkg.NetworkTypeRegTest)
		}
	case btcpkg.NetworkTypeSigNet:
		bit.chainConf = &chaincfg.SigNetParams
		if blockInfo.Chain != "signet" {
			return nil, fmt.Errorf(
				"connecting %s on bitcoind, but config file defines as %s",
				blockInfo.Chain, btcpkg.NetworkTypeSigNet)
		}
	default:
		return nil, errors.New("bitcoin network type is invalid in config")
	}

	// set bitcoin version
	netInfo, err := bit.pkgrpc.GetNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("fail to call bit.GetNetworkInfo(): %w", err)
	}
	if btcpkg.RequiredNodeVersion > btcpkg.NodeVersion(netInfo.Version) {
		return nil, fmt.Errorf(
			"bitcoin core version should be %d +, but version %d is detected",
			btcpkg.RequiredNodeVersion, netInfo.Version)
	}
	bit.version = btcpkg.NodeVersion(netInfo.Version)
	logger.Info("bitcoin rpc server", "version", netInfo.Version)

	// set other information from config
	bit.confirmationBlock = conf.Block.ConfirmationNum
	bit.feeRange.max = conf.Fee.AdjustmentMax
	bit.feeRange.min = conf.Fee.AdjustmentMin

	return &bit, nil
}

// Close disconnect from bitcoin core server
func (b *Bitcoin) Close() {
	if b.btcdClient != nil {
		b.btcdClient.Shutdown()
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
func (b *Bitcoin) Version() domainBTC.Version {
	return domainBTC.Version(b.version)
}

// CoinTypeCode returns CoinTypeCode
func (b *Bitcoin) CoinTypeCode() domainCoin.CoinTypeCode {
	return b.coinTypeCode
}

// GetPkgRPC returns the underlying btcrpc.BTCRPC btcdClient for direct RPC access
func (b *Bitcoin) GetPkgRPC() btcrpc.BTCRPC {
	return b.pkgrpc
}
