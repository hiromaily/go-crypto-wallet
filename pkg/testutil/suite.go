package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

type BTCTestSuite struct {
	suite.Suite
	BTC btcgrp.Bitcoiner
}

func (bts *BTCTestSuite) SetupTest() {
	btc, err := GetBTC()
	bts.NoError(err)
	bts.BTC = btc
}

func (suite *BTCTestSuite) TearDownTest() {
	suite.BTC.Close()
}

type ETHTestSuite struct {
	suite.Suite
	ETH ethgrp.Ethereumer
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
	XRP xrpgrp.Rippler
}

func (xts *XRPTestSuite) SetupTest() {
	xrp, err := GetXRP()
	xts.NoError(err)
	xts.XRP = xrp
}

func (xts *XRPTestSuite) TearDownTest() {
	xts.XRP.Close()
}
