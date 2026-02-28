package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWSCaller is a simple test double for WSCaller.
type mockWSCaller struct {
	callFn func(ctx context.Context, req, res any) error
}

func (m *mockWSCaller) Call(ctx context.Context, req, res any) error {
	return m.callFn(ctx, req, res)
}

// ─── AccountChannels ──────────────────────────────────────────────────────────

func TestAccountChannels_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*AccountChannelsRequest)
			assert.Equal(t, "account_channels", r.Command)
			assert.Equal(t, "rSender", r.Account)
			assert.Equal(t, "rReceiver", r.DestinationAccount)
			p := res.(*ResponseAccountChannels)
			p.Status = "success"
			p.Result.Account = "rSender"
			return nil
		},
	}
	res, err := AccountChannels(context.Background(), caller, "rSender", "rReceiver")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "rSender", res.Result.Account)
}

func TestAccountChannels_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("websocket closed")
		},
	}
	_, err := AccountChannels(context.Background(), caller, "rSender", "rReceiver")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_channels")
}

// ─── AccountInfo ──────────────────────────────────────────────────────────────

func TestAccountInfo_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*AccountInfoRequest)
			assert.Equal(t, "account_info", r.Command)
			assert.Equal(t, "rAddress", r.Account)
			p := res.(*ResponseAccountInfo)
			p.Status = "success"
			p.Result.AccountData.Account = "rAddress"
			p.Result.AccountData.Balance = "1000000"
			return nil
		},
	}
	res, err := AccountInfo(context.Background(), caller, "rAddress")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "rAddress", res.Result.AccountData.Account)
	assert.Equal(t, "1000000", res.Result.AccountData.Balance)
}

func TestAccountInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("connection reset")
		},
	}
	_, err := AccountInfo(context.Background(), caller, "rAddress")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_info")
}

// ─── ServerInfo ───────────────────────────────────────────────────────────────

func TestServerInfo_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*RequestCommand)
			assert.Equal(t, "server_info", r.Command)
			p := res.(*ResponseServerInfo)
			p.Status = "success"
			p.Result.Info.BuildVersion = "1.9.4"
			return nil
		},
	}
	res, err := ServerInfo(context.Background(), caller)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "1.9.4", res.Result.Info.BuildVersion)
}

func TestServerInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("node unavailable")
		},
	}
	_, err := ServerInfo(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server_info")
}

// ─── ValidationCreate ─────────────────────────────────────────────────────────

func TestValidationCreate_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*ValidationCreateRequest)
			assert.Equal(t, "validation_create", r.Command)
			assert.Equal(t, "mysecret", r.Secret)
			p := res.(*ResponseValidationCreate)
			p.Result.Status = "success"
			p.Result.ValidationKey = "BARD OGRE LULL DOCK ROME RISK LOOT CURE"
			p.Result.ValidationPublicKey = "n9KAa2zVWjPHgfzsE3iZ8HAbzJtPrnoh4H2M2HgE7dfqtvyEb1KJ"
			p.Result.ValidationSeed = "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"
			return nil
		},
	}
	res, err := ValidationCreate(context.Background(), caller, "mysecret")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Result.Status)
	assert.NotEmpty(t, res.Result.ValidationKey)
	assert.NotEmpty(t, res.Result.ValidationPublicKey)
}

func TestValidationCreate_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("admin api disabled")
		},
	}
	_, err := ValidationCreate(context.Background(), caller, "mysecret")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation_create")
}

// ─── WalletProposeWithKey ─────────────────────────────────────────────────────

func TestWalletProposeWithKey_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*WalletProposeWithKeyRequest)
			assert.Equal(t, "wallet_propose", r.Command)
			assert.Equal(t, "snoPBrXtMeMyMHUVTgbuqAfg1SUTb", r.Seed)
			assert.Equal(t, "secp256k1", r.KeyType)
			p := res.(*ResponseWalletPropose)
			p.Status = "success"
			p.Result.AccountID = "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
			return nil
		},
	}
	res, err := WalletProposeWithKey(context.Background(), caller, "snoPBrXtMeMyMHUVTgbuqAfg1SUTb", KeyTypeSECP256K1)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh", res.Result.AccountID)
}

func TestWalletProposeWithKey_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("admin api disabled")
		},
	}
	_, err := WalletProposeWithKey(context.Background(), caller, "seed", KeyTypeSECP256K1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wallet_propose")
}

// ─── WalletPropose ────────────────────────────────────────────────────────────

func TestWalletPropose_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, req, res any) error {
			r := req.(*WalletProposeRequest)
			assert.Equal(t, "wallet_propose", r.Command)
			assert.Equal(t, "mypassphrase", r.Passphrase)
			p := res.(*ResponseWalletPropose)
			p.Status = "success"
			p.Result.AccountID = "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59"
			p.Result.MasterSeed = "snoPBrXtMeMyMHUVTgbuqAfg1SUTb"
			return nil
		},
	}
	res, err := WalletPropose(context.Background(), caller, "mypassphrase")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59", res.Result.AccountID)
}

func TestWalletPropose_Error(t *testing.T) {
	t.Parallel()
	caller := &mockWSCaller{
		callFn: func(_ context.Context, _ any, _ any) error {
			return errors.New("admin api disabled")
		},
	}
	_, err := WalletPropose(context.Background(), caller, "passphrase")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wallet_propose")
}

// ─── KeyType ──────────────────────────────────────────────────────────────────

func TestKeyType_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "secp256k1", KeyTypeSECP256K1.String())
	assert.Equal(t, "ed25519", KeyTypeED25519.String())
}
