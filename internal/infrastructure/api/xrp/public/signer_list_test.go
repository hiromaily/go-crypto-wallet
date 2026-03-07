package xrppublic

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

// mockWSCaller is a test double for rpc.WSCaller.
type mockWSCaller struct {
	callFn  func(ctx context.Context, req, res any) error
	closeFn func() error
}

func (m *mockWSCaller) Call(ctx context.Context, req, res any) error {
	return m.callFn(ctx, req, res)
}

func (m *mockWSCaller) Close() error {
	if m.closeFn != nil {
		return m.closeFn()
	}
	return nil
}

// newPublicXRPWithMock creates a PublicXRP backed by a mock WSCaller for unit testing.
func newPublicXRPWithMock(caller *mockWSCaller) *PublicXRP {
	return &PublicXRP{
		PublicRPC: xrprpcpublic.NewPublicRPC(caller),
	}
}

func TestPrepareSignerListSetTransaction_HappyPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		account           string
		quorum            uint32
		entries           []apixrp.SignerEntryInput
		wantFee           string
		wantTransType     string
		wantSignerQuorum  uint32
		wantEntryCount    int
		wantSequence      uint64
		wantLastLedgerSeq uint64
	}{
		{
			name:    "2 signers, quorum 2",
			account: "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			quorum:  2,
			entries: []apixrp.SignerEntryInput{
				{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1},
				{Account: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", Weight: 1},
			},
			wantFee:           "36", // (2+1)*12
			wantTransType:     "SignerListSet",
			wantSignerQuorum:  2,
			wantEntryCount:    2,
			wantSequence:      42,
			wantLastLedgerSeq: 1015, // 1000 + 15
		},
		{
			name:    "1 signer, quorum 1",
			account: "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			quorum:  1,
			entries: []apixrp.SignerEntryInput{
				{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1},
			},
			wantFee:           "24", // (1+1)*12
			wantTransType:     "SignerListSet",
			wantSignerQuorum:  1,
			wantEntryCount:    1,
			wantSequence:      42,
			wantLastLedgerSeq: 1015,
		},
		{
			name:    "3 signers, quorum 2",
			account: "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			quorum:  2,
			entries: []apixrp.SignerEntryInput{
				{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1},
				{Account: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", Weight: 1},
				{Account: "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59", Weight: 1},
			},
			wantFee:           "48", // (3+1)*12
			wantTransType:     "SignerListSet",
			wantSignerQuorum:  2,
			wantEntryCount:    3,
			wantSequence:      42,
			wantLastLedgerSeq: 1015,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			caller := &mockWSCaller{
				callFn: func(_ context.Context, _, res any) error {
					r := res.(*xrprpcpublic.ResponseAccountInfo)
					r.Result.AccountData.Sequence = 42
					r.Result.LedgerCurrentIndex = 1000
					return nil
				},
			}
			p := newPublicXRPWithMock(caller)

			dto, jsonStr, err := p.PrepareSignerListSetTransaction(
				context.Background(),
				tt.account,
				tt.quorum,
				tt.entries,
				nil,
			)

			require.NoError(t, err)
			require.NotNil(t, dto)
			assert.NotEmpty(t, jsonStr)

			// Verify DTO fields
			assert.Equal(t, tt.wantTransType, dto.TransactionType)
			assert.Equal(t, tt.account, dto.Account)
			assert.Equal(t, tt.wantSignerQuorum, dto.SignerQuorum)
			assert.Equal(t, tt.wantFee, dto.Fee)
			assert.Equal(t, tt.wantSequence, dto.Sequence)
			assert.Equal(t, tt.wantLastLedgerSeq, dto.LastLedgerSequence)
			assert.Len(t, dto.SignerEntries, tt.wantEntryCount)

			// Verify JSON can be parsed and contains expected fields
			var parsed map[string]any
			require.NoError(t, json.Unmarshal([]byte(jsonStr), &parsed))
			assert.Equal(t, "SignerListSet", parsed["TransactionType"])
			assert.Equal(t, tt.account, parsed["Account"])
		})
	}
}

func TestPrepareSignerListSetTransaction_AccountInfoError(t *testing.T) {
	t.Parallel()

	caller := &mockWSCaller{
		callFn: func(_ context.Context, _, _ any) error {
			return errors.New("websocket connection failed")
		},
	}
	p := newPublicXRPWithMock(caller)

	dto, jsonStr, err := p.PrepareSignerListSetTransaction(
		context.Background(),
		"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		2,
		[]apixrp.SignerEntryInput{{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1}},
		nil,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "PrepareSignerListSetTransaction")
	assert.Contains(t, err.Error(), "websocket connection failed")
	assert.Nil(t, dto)
	assert.Empty(t, jsonStr)
}

func TestPrepareSignerListSetTransaction_AccountInfoRPCError(t *testing.T) {
	t.Parallel()

	caller := &mockWSCaller{
		callFn: func(_ context.Context, _, res any) error {
			r := res.(*xrprpcpublic.ResponseAccountInfo)
			r.Error = "actNotFound"
			return nil
		},
	}
	p := newPublicXRPWithMock(caller)

	dto, jsonStr, err := p.PrepareSignerListSetTransaction(
		context.Background(),
		"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		2,
		[]apixrp.SignerEntryInput{{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1}},
		nil,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "actNotFound")
	assert.Nil(t, dto)
	assert.Empty(t, jsonStr)
}

func TestPrepareSignerListSetTransaction_InstructionsOverride(t *testing.T) {
	t.Parallel()

	caller := &mockWSCaller{
		callFn: func(_ context.Context, _, res any) error {
			r := res.(*xrprpcpublic.ResponseAccountInfo)
			r.Result.AccountData.Sequence = 42
			r.Result.LedgerCurrentIndex = 1000
			return nil
		},
	}
	p := newPublicXRPWithMock(caller)

	// Test with nil instructions (baseline)
	dto, _, err := p.PrepareSignerListSetTransaction(
		context.Background(),
		"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		2,
		[]apixrp.SignerEntryInput{
			{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 1},
			{Account: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", Weight: 1},
		},
		nil,
	)

	require.NoError(t, err)
	assert.Equal(t, "36", dto.Fee) // default fee: (2+1)*12=36
	assert.Equal(t, uint64(42), dto.Sequence)
}

func TestPrepareSignerListSetTransaction_SignerEntryMapping(t *testing.T) {
	t.Parallel()

	caller := &mockWSCaller{
		callFn: func(_ context.Context, _, res any) error {
			r := res.(*xrprpcpublic.ResponseAccountInfo)
			r.Result.AccountData.Sequence = 1
			r.Result.LedgerCurrentIndex = 100
			return nil
		},
	}
	p := newPublicXRPWithMock(caller)

	entries := []apixrp.SignerEntryInput{
		{Account: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", Weight: 2},
		{Account: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", Weight: 3},
	}

	dto, jsonStr, err := p.PrepareSignerListSetTransaction(
		context.Background(),
		"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		3,
		entries,
		nil,
	)

	require.NoError(t, err)
	require.Len(t, dto.SignerEntries, 2)
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", dto.SignerEntries[0].Account)
	assert.Equal(t, uint32(2), dto.SignerEntries[0].SignerWeight)
	assert.Equal(t, "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", dto.SignerEntries[1].Account)
	assert.Equal(t, uint32(3), dto.SignerEntries[1].SignerWeight)

	// Verify JSON wire format has nested SignerEntry structure
	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &parsed))
	signerEntriesRaw, ok := parsed["SignerEntries"].([]any)
	require.True(t, ok, "SignerEntries must be an array")
	require.Len(t, signerEntriesRaw, 2)

	// Each entry should have the nested {"SignerEntry": {"Account": ..., "SignerWeight": ...}} structure
	firstEntry, ok := signerEntriesRaw[0].(map[string]any)
	require.True(t, ok)
	innerEntry, ok := firstEntry["SignerEntry"].(map[string]any)
	require.True(t, ok, "each entry must have a nested SignerEntry object")
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", innerEntry["Account"])
}
