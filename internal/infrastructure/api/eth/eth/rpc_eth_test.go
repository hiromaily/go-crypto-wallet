//go:build integration
// +build integration

package eth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/testutil"
)

type ethTest struct {
	testutil.ETHTestSuite
}

// TestSyncing is test for Syncing
func (et *ethTest) TestSyncing() {
	ctx := context.Background()
	res, isSyncing, err := et.ETH.Syncing(ctx)
	et.NoError(err)
	et.T().Log("resMap:", res)
	et.T().Log("isSyncing:", isSyncing)
}

// TestProtocolVersion is test for ProtocolVersion
func (et *ethTest) TestProtocolVersion() {
	ctx := context.Background()
	protocolVer, err := et.ETH.ProtocolVersion(ctx)
	et.NoError(err)
	et.T().Log("ProtocolVersion:", protocolVer)
}

// TestCoinbase is test for Coinbase
func (et *ethTest) TestCoinbase() {
	ctx := context.Background()
	addr, err := et.ETH.Coinbase(ctx)
	et.NoError(err)
	et.T().Log("coinbase address:", addr)
}

// TestAccounts is test for Accounts
func (et *ethTest) TestAccounts() {
	ctx := context.Background()
	accounts, err := et.ETH.Accounts(ctx)
	et.NoError(err)
	for _, account := range accounts {
		et.T().Log("address:", account)
	}
}

// TestBlockNumber is test for BlockNumber
func (et *ethTest) TestBlockNumber() {
	ctx := context.Background()
	blockNum, err := et.ETH.BlockNumber(ctx)
	et.NoError(err)
	et.T().Log("BlockNumber:", blockNum)

	blockNum, err = et.ETH.EnsureBlockNumber(ctx, 100)
	et.NoError(err)
	et.T().Log("EnsureBlockNumber:", blockNum)
}

// TestGetBalance is test for GetBalance
func (et *ethTest) TestGetBalance() {
	tags := []domainEthereum.QuantityTag{
		domainEthereum.QuantityTagLatest,
		domainEthereum.QuantityTagPending,
		// domainEthereum.QuantityTagEarliest,
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
			ctx := context.Background()
			for _, tag := range tags {
				balance, err := et.ETH.GetBalance(ctx, tt.args.addr, tag)
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
	tags := []domainEthereum.QuantityTag{
		domainEthereum.QuantityTagLatest,
		domainEthereum.QuantityTagPending,
		// domainEthereum.QuantityTagEarliest,
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
			ctx := context.Background()
			for _, tag := range tags {
				count, err := et.ETH.GetTransactionCount(ctx, tt.args.addr, tag)
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
			ctx := context.Background()
			count, err := et.ETH.GetBlockTransactionCountByNumber(ctx, tt.args.txNum)
			et.NoError(err)
			if err == nil {
				t.Logf("GetBlockTransactionCountByNumber: %d", count.Uint64())
			}

			count, err = et.ETH.GetUncleCountByBlockNumber(ctx, tt.args.txNum)
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
			ctx := context.Background()
			blockInfo, err := et.ETH.GetBlockByNumber(ctx, tt.args.txNum)
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
