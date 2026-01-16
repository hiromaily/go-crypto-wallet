//go:build integration
// +build integration

package eth_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/testutil"
)

type adminTest struct {
	testutil.ETHTestSuite
}

// TestAddPeer is test for AddPeer
// https://github.com/ethereum/go-ethereum/blob/master/params/bootnodes.go
func (at *adminTest) TestAddPeer() {
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
			args: args{params.SepoliaBootnodes[0]},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{params.SepoliaBootnodes[1]},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{params.SepoliaBootnodes[2]},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{"enode://foobar"},
			want: want{true}, // invalid enode: does not contain node ID
		},
		{
			name: "wrong schema",
			args: args{
				"xxxx://f4a9c6ee28586009fb5a96c8af13a58ed6d8315a9eee4772212c1d4d9cebe5a8b8a78ea4434f318726317d04a3f53" +
					"1a1ef0420cf9752605a562cfe858c46e263@213.186.16.82:30303",
			},
			want: want{true}, // invalid enode: missing 'enr:' prefix for base64-encoded record
		},
	}

	for _, tt := range tests {
		at.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := at.ETH.AddPeer(ctx, tt.args.addr)
			if (err != nil) != tt.want.isErr {
				t.Errorf("AddPeer() = %v, want error = %v", err, tt.want.isErr)
				return
			}
		})
	}
}

// TestAdminDataDir is test for AdminDataDir
func (at *adminTest) TestAdminDataDir() {
	ctx := context.Background()
	dirName, err := at.ETH.AdminDataDir(ctx)
	at.NoError(err)
	at.T().Log(dirName) // /Users/hy/Library/Ethereum/goerli
}

// TestNodeInfo is test for NodeInfo
func (at *adminTest) TestNodeInfo() {
	ctx := context.Background()
	nodeInfo, err := at.ETH.NodeInfo(ctx)
	at.NoError(err)

	t := at.T()
	t.Log("Name:", nodeInfo.Name)             // Geth/v1.9.13-stable/darwin-amd64/go1.14.2
	t.Log("ID:", nodeInfo.ID)                 // 2250fc365755468c831afcea6df37aca52754309060923daee832eb0d7cc49a4
	t.Log("IP:", nodeInfo.IP)                 // xx.xx.xx.xx
	t.Log("ListenAddr:", nodeInfo.ListenAddr) // [::]:30303
	t.Log("Ports:", nodeInfo.Ports)           // {30303 30303}
	t.Log("Protocols:", nodeInfo.Protocols)
	t.Log("Enode:", nodeInfo.Enode) // enode://xxxxx
	t.Log("ENR:", nodeInfo.ENR)     // ??
}

// TestAdminPeers is test for AdminPeers
func (at *adminTest) TestAdminPeers() {
	ctx := context.Background()
	adminPeers, err := at.ETH.AdminPeers(ctx)
	at.NoError(err)
	for _, peer := range adminPeers {
		at.T().Log(peer)
	}
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(adminTest))
}
