package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
)

type BTCTestSuite struct {
	suite.Suite
	BTC bitcoin.Bitcoiner
}

func (bts *BTCTestSuite) SetupTest() {
	btc, err := GetBTC()
	bts.NoError(err)
	bts.BTC = btc
}

func (bts *BTCTestSuite) TearDownTest() {
	bts.BTC.Close()
}

type ETHTestSuite struct {
	suite.Suite
	ETH ethereum.Ethereumer
}

func (ets *ETHTestSuite) SetupTest() {
	eth, err := GetETH()
	ets.NoError(err)
	ets.ETH = eth
}

func (ets *ETHTestSuite) TearDownTest() {
	ets.ETH.Close()
}

type XRPTestSuite struct {
	suite.Suite
	XRP ripple.Rippler
}

func (xts *XRPTestSuite) SetupTest() {
	xrp, err := GetXRP()
	xts.NoError(err)
	xts.XRP = xrp
}

func (xts *XRPTestSuite) TearDownTest() {
	_ = xts.XRP.Close() // Best effort cleanup
}
