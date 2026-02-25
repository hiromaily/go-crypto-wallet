package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkTypeETHConstants(t *testing.T) {
	assert.Equal(t, NetworkTypeETH("mainnet"), NetworkTypeETHMainNet)
	assert.Equal(t, NetworkTypeETH("sepolia"), NetworkTypeETHSepolia)
	assert.Equal(t, NetworkTypeETH("holesky"), NetworkTypeETHHolesky)
	assert.Equal(t, NetworkTypeETH("anvil"), NetworkTypeETHAnvil)
	assert.Equal(t, NetworkTypeETH("local"), NetworkTypeETHLocal)
}

func TestETHNodeTypeConstants(t *testing.T) {
	assert.Equal(t, ETHNodeType("anvil"), ETHNodeTypeAnvil)
	assert.Equal(t, ETHNodeType("geth"), ETHNodeTypeGeth)
}

func TestChainIDForNetwork(t *testing.T) {
	tests := []struct {
		network  NetworkTypeETH
		wantID   uint64
		wantZero bool
	}{
		{NetworkTypeETHMainNet, 1, false},
		{NetworkTypeETHSepolia, 11155111, false},
		{NetworkTypeETHHolesky, 17000, false},
		{NetworkTypeETHAnvil, 1337, false},
		{NetworkTypeETHLocal, 1337, false},
		{"unknown", 0, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.network), func(t *testing.T) {
			got := ChainIDForNetwork(tt.network)
			if tt.wantZero {
				assert.Equal(t, uint64(0), got)
			} else {
				assert.Equal(t, tt.wantID, got)
			}
		})
	}
}
