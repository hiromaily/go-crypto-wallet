package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
