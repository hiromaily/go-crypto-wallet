package di_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

// TestNotImplementedAddMultisigSignatureUseCase verifies that the no-op use case
// returned by NewNotImplementedAddMultisigSignatureUseCase satisfies the interface
// and returns a clear "not yet implemented" error on Execute.
func TestNotImplementedAddMultisigSignatureUseCase_Execute_ReturnsNotImplementedError(t *testing.T) {
	t.Parallel()

	uc := di.NewNotImplementedAddMultisigSignatureUseCase()

	require.NotNil(t, uc, "use case must not be nil — nil interface causes silent panics at call sites")

	output, err := uc.Execute(context.Background(), watchusecase.AddMultisigSignatureInput{
		TxUUID:        "test-uuid",
		SignerAccount: "rSigner1",
		SignedTxBlob:  "DEADBEEF",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
	assert.Empty(t, output.CombinedTxBlob)
	assert.False(t, output.IsReady)
}
