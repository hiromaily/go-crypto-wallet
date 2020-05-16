package eth_test

import (
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestSyncing is test for Syncing
func TestSyncing(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	res, isSyncing, err := et.Syncing()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("resMap:", res)
	t.Log("isSyncing:", isSyncing)
}

// TestProtocolVersion is test for ProtocolVersion
func TestProtocolVersion(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	protocolVer, err := et.ProtocolVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ProtocolVersion:", protocolVer)
}

// TestCoinbase is test for Coinbase
func TestCoinbase(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	addr, err := et.Coinbase()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("coinbase address:", addr)
}

// TestAccounts is test for Accounts
func TestAccounts(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	accounts, err := et.Accounts()
	if err != nil {
		t.Fatal(err)
	}
	for _, account := range accounts {
		t.Log("address:", account)
	}
}

// TestBlockNumber is test for BlockNumber
func TestBlockNumber(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	blockNum, err := et.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("blockNumber:", blockNum)

	blockNum2, err := et.EnsureBlockNumber(100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("blockNumber2:", blockNum2)
}

// TestGetBalance is test for GetBalance
func TestGetBalance(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	tags := []eth.QuantityTag{
		eth.QuantityTagLatest,
		eth.QuantityTagPending,
		//eth.QuantityTagEarliest,
	}

	type args struct {
		addr string
	}
	type want struct {
		balance uint64
		isErr   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{addr: "0xEA247646137F74C38e04f3012db483d77F3dEc59"},
			want: want{50000000000000000, false},
		},
		{
			name: "happy path",
			args: args{addr: "0x6B91314E559D3FB5b40f9F51582631a7b5C610ef"},
			want: want{50000000000000000, false},
		},
		{
			name: "happy path",
			args: args{addr: "0x0aC5c95EB979C41CFa2C6BdF8e5515F966fEc103"},
			want: want{50000000000000000, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tag := range tags {
				balance, err := et.GetBalance(tt.args.addr, tag)
				if (err == nil) == tt.want.isErr {
					t.Errorf("GetBalance() = %v, want error = %v", err, tt.want.isErr)
					return
				}
				if balance != nil {
					t.Logf("quantityTag: %s, balance: %d", tag, balance.Uint64())
				}
			}
		})
	}
}

// TestGetStoreageAt is test for GetStoreageAt
//func TestGetStoreageAt(t *testing.T) {
//	//t.SkipNow()
//	et := testutil.GetETH()
//
//	tags := []eth.QuantityTag{
//		eth.QuantityTagLatest,
//		eth.QuantityTagPending,
//		//eth.QuantityTagEarliest,
//	}
//
//	type args struct {
//		addr string
//	}
//	type want struct {
//		isErr bool
//	}
//	tests := []struct {
//		name string
//		args args
//		want want
//	}{
//		{
//			name: "happy path",
//			args: args{addr: "0xEA247646137F74C38e04f3012db483d77F3dEc59"},
//			want: want{false},
//		},
//		{
//			name: "happy path",
//			args: args{addr: "0x6B91314E559D3FB5b40f9F51582631a7b5C610ef"},
//			want: want{false},
//		},
//		{
//			name: "happy path",
//			args: args{addr: "0x0aC5c95EB979C41CFa2C6BdF8e5515F966fEc103"},
//			want: want{false},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			for _, tag := range tags {
//				position, err := et.GetStoreageAt(tt.args.addr, tag)
//				if (err == nil) == tt.want.isErr {
//					t.Errorf("GetStoreageAt() = %v, want error = %v", err, tt.want.isErr)
//					return
//				}
//				t.Logf("quantityTag: %s, position: %s", tag, position)
//			}
//		})
//	}
//}

// TestGetTransactionCount is test for GetTransactionCount
func TestGetTransactionCount(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	tags := []eth.QuantityTag{
		eth.QuantityTagLatest,
		eth.QuantityTagPending,
		//eth.QuantityTagEarliest,
	}

	type args struct {
		addr string
	}
	type want struct {
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{addr: "0xEA247646137F74C38e04f3012db483d77F3dEc59"},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{addr: "0x6B91314E559D3FB5b40f9F51582631a7b5C610ef"},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{addr: "0x0aC5c95EB979C41CFa2C6BdF8e5515F966fEc103"},
			want: want{false},
		},
		{
			name: "happy path, outer address",
			args: args{addr: "0x8cED5ad0d8dA4Ec211C17355Ed3DBFEC4Cf0E5b9"},
			want: want{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tag := range tags {
				count, err := et.GetTransactionCount(tt.args.addr, tag)
				if (err == nil) == tt.want.isErr {
					t.Errorf("GetTransactionCount() = %v, want error = %v", err, tt.want.isErr)
					return
				}
				if count != nil {
					t.Logf("quantityTag: %s, count: %d", tag, count.Uint64())
				}
			}
		})
	}
}

// TestGetBlockTransactionCountByBlockHash is test for GetBlockTransactionCountByBlockHash
//func TestGetBlockTransactionCountByBlockHash(t *testing.T) {
//	t.SkipNow()
//	et := testutil.GetETH()
//
//	type args struct {
//		txHash string
//	}
//	type want struct {
//		isErr bool
//	}
//	tests := []struct {
//		name string
//		args args
//		want want
//	}{
//		{
//			name: "happy path 2706436",
//			args: args{txHash: "0x5a3d8c98fe321bf59806c45f56060f90e0e23070afb5bc4edfc593f6d9e485bc"},
//			want: want{false},
//		},
//		{
//			name: "happy path 2706435",
//			args: args{txHash: "0x77a67dc9ed3b08dc91f593ea1fdde621f0c0654dbb5e3f28a29cf8390f3e61a5"},
//			want: want{false},
//		},
//		{
//			name: "happy path 2606434",
//			args: args{txHash: "0xdaf6d232c009af5f68788d969ac3b501a1581325fed40e8c5472659c56b04cda"},
//			want: want{false},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			count, err := et.GetBlockTransactionCountByBlockHash(tt.args.txHash)
//			if (err == nil) == tt.want.isErr {
//				t.Errorf("GetBlockTransactionCountByBlockHash() = %v, want error = %v", err, tt.want.isErr)
//				return
//			}
//			if count != nil {
//				t.Logf("transactionCount: %d", count.Uint64())
//			}
//		})
//	}
//}

// TestGetBlockTransactionCountByNumber is test for GetBlockTransactionCountByNumber
func TestGetBlockTransactionCountByNumber(t *testing.T) {
	et := testutil.GetETH()

	type args struct {
		txNum uint64
	}
	type want struct {
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 2706436",
			args: args{txNum: 2706436},
			want: want{false},
		},
		{
			name: "happy path 2706435",
			args: args{txNum: 2706435},
			want: want{false},
		},
		{
			name: "happy path 2706434",
			args: args{txNum: 2606434},
			want: want{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := et.GetBlockTransactionCountByNumber(tt.args.txNum)
			if (err == nil) == tt.want.isErr {
				t.Errorf("GetBlockTransactionCountByNumber() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if count != nil {
				t.Logf("transactionCount: %d", count.Uint64())
			}
		})
	}
}
