//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

type ethTest struct {
	testutil.ETHTestSuite
}

// TestSyncing is test for Syncing
func (et *ethTest) TestSyncing() {
	res, isSyncing, err := et.ETH.Syncing()
	et.NoError(err)
	et.T().Log("resMap:", res)
	et.T().Log("isSyncing:", isSyncing)
}

// TestProtocolVersion is test for ProtocolVersion
func (et *ethTest) TestProtocolVersion() {
	protocolVer, err := et.ETH.ProtocolVersion()
	et.NoError(err)
	et.T().Log("ProtocolVersion:", protocolVer)
}

// TestCoinbase is test for Coinbase
func (et *ethTest) TestCoinbase() {
	addr, err := et.ETH.Coinbase()
	et.NoError(err)
	et.T().Log("coinbase address:", addr)
}

// TestAccounts is test for Accounts
func (et *ethTest) TestAccounts() {
	accounts, err := et.ETH.Accounts()
	et.NoError(err)
	for _, account := range accounts {
		et.T().Log("address:", account)
	}
}

// TestBlockNumber is test for BlockNumber
func (et *ethTest) TestBlockNumber() {
	blockNum, err := et.ETH.BlockNumber()
	et.NoError(err)
	et.T().Log("BlockNumber:", blockNum)

	blockNum, err = et.ETH.EnsureBlockNumber(100)
	et.NoError(err)
	et.T().Log("EnsureBlockNumber:", blockNum)
}

// TestGetBalance is test for GetBalance
func (et *ethTest) TestGetBalance() {
	tags := []eth.QuantityTag{
		eth.QuantityTagLatest,
		eth.QuantityTagPending,
		// eth.QuantityTagEarliest,
	}

	type args struct {
		addr string
	}
	type want struct {
		balance uint64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{addr: "0xEA247646137F74C38e04f3012db483d77F3dEc59"},
			want: want{50000000000000000},
		},
		{
			name: "happy path",
			args: args{addr: "0x6B91314E559D3FB5b40f9F51582631a7b5C610ef"},
			want: want{50000000000000000},
		},
		{
			name: "happy path",
			args: args{addr: "0x0aC5c95EB979C41CFa2C6BdF8e5515F966fEc103"},
			want: want{50000000000000000},
		},
	}
	for _, tt := range tests {
		et.T().Run(tt.name, func(t *testing.T) {
			for _, tag := range tags {
				balance, err := et.ETH.GetBalance(tt.args.addr, tag)
				et.NoError(err)
				if err == nil {
					t.Logf("quantityTag: %s, balance: %d", tag, balance.Uint64())
				}
			}
		})
	}
}

// TestGetTransactionCount is test for GetTransactionCount
func (et *ethTest) TestGetTransactionCount() {
	tags := []eth.QuantityTag{
		eth.QuantityTagLatest,
		eth.QuantityTagPending,
		// eth.QuantityTagEarliest,
	}

	type args struct {
		addr string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{addr: "0xEA247646137F74C38e04f3012db483d77F3dEc59"},
			want: want{},
		},
		{
			name: "happy path",
			args: args{addr: "0x6B91314E559D3FB5b40f9F51582631a7b5C610ef"},
			want: want{},
		},
		{
			name: "happy path",
			args: args{addr: "0x0aC5c95EB979C41CFa2C6BdF8e5515F966fEc103"},
			want: want{},
		},
		{
			name: "happy path, outer address",
			args: args{addr: "0x8cED5ad0d8dA4Ec211C17355Ed3DBFEC4Cf0E5b9"},
			want: want{},
		},
	}
	for _, tt := range tests {
		et.T().Run(tt.name, func(t *testing.T) {
			for _, tag := range tags {
				count, err := et.ETH.GetTransactionCount(tt.args.addr, tag)
				et.NoError(err)
				if err == nil {
					t.Logf("quantityTag: %s, count: %d", tag, count.Uint64())
				}
			}
		})
	}
}

// TestGetBlockTransactionCountByNumber is test for GetBlockTransactionCountByNumber
func (et *ethTest) TestGetBlockTransactionCountByNumber() {
	type args struct {
		txNum uint64
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 2706436",
			args: args{txNum: 2706436},
			want: want{},
		},
		{
			name: "happy path 2706435",
			args: args{txNum: 2706435},
			want: want{},
		},
		{
			name: "happy path 2706434",
			args: args{txNum: 2606434},
			want: want{},
		},
	}
	for _, tt := range tests {
		et.T().Run(tt.name, func(t *testing.T) {
			count, err := et.ETH.GetBlockTransactionCountByNumber(tt.args.txNum)
			et.NoError(err)
			if err == nil {
				t.Logf("GetBlockTransactionCountByNumber: %d", count.Uint64())
			}

			count, err = et.ETH.GetUncleCountByBlockNumber(tt.args.txNum)
			et.NoError(err)
			if err == nil {
				t.Logf("GetUncleCountByBlockNumber: %d", count.Uint64())
			}
		})
	}
}

// TestGetBlockByNumber is test for GetBlockByNumber
func (et *ethTest) TestGetBlockByNumber() {
	type args struct {
		txNum uint64
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 2706436",
			args: args{txNum: 2706436},
			want: want{},
		},
		{
			name: "happy path 2706435",
			args: args{txNum: 2706435},
			want: want{},
		},
		{
			name: "happy path 2706434",
			args: args{txNum: 2606434},
			want: want{},
		},
	}
	for _, tt := range tests {
		et.T().Run(tt.name, func(t *testing.T) {
			blockInfo, err := et.ETH.GetBlockByNumber(tt.args.txNum)
			et.NoError(err)
			if err == nil {
				t.Logf("blockInfo: %v", blockInfo)
			}
		})
	}
}

func TestRPCEthTestSuite(t *testing.T) {
	suite.Run(t, new(ethTest))
}
