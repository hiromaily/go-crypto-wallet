package xrp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	xrpapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/mocks"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// setSignerListDeps holds all mock dependencies for SetSignerListUseCase tests.
type setSignerListDeps struct {
	preparer        *xrpapiamocks.MockSignerListPreparer
	signerListRepo  *coldmocks.MockXRPSignerListRepositorier
	signerEntryRepo *coldmocks.MockXRPSignerEntryRepositorier
	txFileRepo      *storagemocks.MockTransactionFileRepositorier
	uuidHandler     pkguuid.UUIDHandler
}

func newSetSignerListDeps(t *testing.T) *setSignerListDeps {
	t.Helper()
	return &setSignerListDeps{
		preparer:        xrpapiamocks.NewMockSignerListPreparer(t),
		signerListRepo:  coldmocks.NewMockXRPSignerListRepositorier(t),
		signerEntryRepo: coldmocks.NewMockXRPSignerEntryRepositorier(t),
		txFileRepo:      storagemocks.NewMockTransactionFileRepositorier(t),
		uuidHandler:     pkguuid.NewGoogleUUIDHandler(),
	}
}

func createSetSignerListUseCase(deps *setSignerListDeps) watchusecase.SetSignerListUseCase {
	return xrp.NewSetSignerListUseCase(
		deps.preparer,
		deps.uuidHandler,
		deps.signerListRepo,
		deps.signerEntryRepo,
		deps.txFileRepo,
	)
}

func TestSetSignerListUseCase_Execute_EmptyAccount(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "",
		SignerQuorum:   2,
		SignerEntries: []watchusecase.SignerEntry{
			{Account: "rSigner1", Weight: 1},
			{Account: "rSigner2", Weight: 1},
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account address is required")
	assert.Empty(t, output.FileName)
}

func TestSetSignerListUseCase_Execute_EmptySignerEntries(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "rAccount123",
		SignerQuorum:   1,
		SignerEntries:  []watchusecase.SignerEntry{}, // empty — at least one required
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signer entries")
	assert.Empty(t, output.FileName)
}

func TestSetSignerListUseCase_Execute_QuorumExceedsWeight(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "rAccount123",
		SignerQuorum:   5, // exceeds total weight of 2
		SignerEntries: []watchusecase.SignerEntry{
			{Account: "rSigner1", Weight: 1},
			{Account: "rSigner2", Weight: 1},
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid signer entries")
	assert.Empty(t, output.FileName)
}

func TestSetSignerListUseCase_Execute_PrepareTransactionError(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	deps.preparer.EXPECT().
		PrepareSignerListSetTransaction(mock.Anything, "rAccount123", uint32(2), mock.Anything, mock.Anything).
		Return(nil, "", errors.New("rpc error"))

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "rAccount123",
		SignerQuorum:   2,
		SignerEntries: []watchusecase.SignerEntry{
			{Account: "rSigner1", Weight: 1},
			{Account: "rSigner2", Weight: 1},
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to prepare SignerListSet transaction")
	assert.Contains(t, err.Error(), "rpc error")
	assert.Empty(t, output.FileName)
}

func TestSetSignerListUseCase_Execute_DeactivateError(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	signerListSetDTO := &dtoxrp.SignerListSetTxInput{
		TransactionType: "SignerListSet",
		Account:         "rAccount123",
		SignerQuorum:    2,
	}

	deps.preparer.EXPECT().
		PrepareSignerListSetTransaction(mock.Anything, "rAccount123", uint32(2), mock.Anything, mock.Anything).
		Return(signerListSetDTO, `{"TransactionType":"SignerListSet"}`, nil)
	deps.signerListRepo.EXPECT().
		DeactivateByAccountID(mock.Anything, "rAccount123").
		Return(errors.New("db error"))

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "rAccount123",
		SignerQuorum:   2,
		SignerEntries: []watchusecase.SignerEntry{
			{Account: "rSigner1", Weight: 1},
			{Account: "rSigner2", Weight: 1},
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to deactivate existing signer lists")
	assert.Empty(t, output.FileName)
}

func TestSetSignerListUseCase_Execute_Success(t *testing.T) {
	t.Parallel()
	deps := newSetSignerListDeps(t)
	uc := createSetSignerListUseCase(deps)

	signerListSetDTO := &dtoxrp.SignerListSetTxInput{
		TransactionType: "SignerListSet",
		Account:         "rAccount123",
		SignerQuorum:    2,
	}
	txJSON := `{"TransactionType":"SignerListSet","Account":"rAccount123","SignerQuorum":2}`

	deps.preparer.EXPECT().
		PrepareSignerListSetTransaction(mock.Anything, "rAccount123", uint32(2), mock.Anything, mock.Anything).
		Return(signerListSetDTO, txJSON, nil)
	deps.signerListRepo.EXPECT().
		DeactivateByAccountID(mock.Anything, "rAccount123").
		Return(nil)
	deps.signerListRepo.EXPECT().
		Insert(mock.Anything, mock.Anything).
		Return(int64(7), nil)
	deps.signerEntryRepo.EXPECT().
		Insert(mock.Anything, mock.MatchedBy(func(e any) bool { return e != nil })).
		Return(int64(1), nil).Times(2)
	deps.txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeTransfer, domainTx.TxTypeUnsigned, int64(0), 0).
		Return("/tmp/set_signer_list_0.txt")
	deps.txFileRepo.EXPECT().
		WriteFileSlice("/tmp/set_signer_list_0.txt", mock.MatchedBy(func(lines []string) bool {
			return len(lines) == 2 && lines[0] == "rAccount123"
		})).
		Return("/tmp/set_signer_list_0.txt", nil)

	output, err := uc.Execute(context.Background(), watchusecase.SetSignerListInput{
		AccountAddress: "rAccount123",
		SignerQuorum:   2,
		SignerEntries: []watchusecase.SignerEntry{
			{Account: "rSigner1", Weight: 1},
			{Account: "rSigner2", Weight: 1},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "/tmp/set_signer_list_0.txt", output.FileName)
	assert.Equal(t, txJSON, output.TxJSON)
	assert.Equal(t, int64(7), output.SignerListID)
}
