package testutil

import (
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

type BTCTestSuite struct {
	suite.Suite
	BTC btcgrp.Bitcoiner
}

func (suite *BTCTestSuite) SetupTest() {
	btc, err := GetBTC()
	suite.NoError(err)
	suite.BTC = btc
}

func (suite *BTCTestSuite) TearDownTest() {
	suite.BTC.Close()
}
