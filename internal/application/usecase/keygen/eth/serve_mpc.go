package eth

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type serveMPCUseCase struct {
	shardStorage apieth.MPCKeyShardStorage
	transport    apieth.MPCInboundTransport
	signingNode  apieth.MPCSigningNodePort
}

// NewServeMPCUseCase constructs the use case for running the MPC node gRPC signing daemon.
func NewServeMPCUseCase(
	shardStorage apieth.MPCKeyShardStorage,
	transport apieth.MPCInboundTransport,
	signingNode apieth.MPCSigningNodePort,
) keygenusecase.ServeMPCUseCase {
	return &serveMPCUseCase{
		shardStorage: shardStorage,
		transport:    transport,
		signingNode:  signingNode,
	}
}

// Serve starts the MPC node daemon. It loads the key shard, starts the inbound transport
// listener, waits for the coordinator's InitSigning RPC, and runs the TSS signing protocol
// until it completes or ctx is cancelled.
func (u *serveMPCUseCase) Serve(ctx context.Context, input keygenusecase.ServeMPCInput) error {
	// 1. Load key shard at startup.
	shardJSON, err := u.shardStorage.LoadShard(ctx, input.ShardPath, input.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load key shard from %s: %w", input.ShardPath, err)
	}

	// 2. Start gRPC listener.
	if err = u.transport.Listen(ctx, input.ListenAddr); err != nil {
		return fmt.Errorf("failed to start MPC inbound listener on %s: %w", input.ListenAddr, err)
	}
	defer func() { _ = u.transport.Close() }()

	logger.Info("MPC node server started", "listenAddr", input.ListenAddr, "partyID", input.PartyID)

	// 3. Wait for the coordinator's InitSigning RPC to deliver session parameters.
	sessionInfo, err := u.transport.AwaitSessionInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to await signing session info: %w", err)
	}

	logger.Info("MPC signing session received",
		"sessionID", sessionInfo.SessionID,
		"threshold", sessionInfo.Threshold,
		"parties", sessionInfo.PartyIDs,
	)

	// 4. Verify that signer.Hash(raw_tx) matches the hash from InitSigning to prevent blind signing.
	if err = verifySigningHash(sessionInfo.RawTxHex, sessionInfo.Hash); err != nil {
		return fmt.Errorf("tx hash verification failed, refusing to sign: %w", err)
	}

	// 5. Run the TSS signing protocol.
	return u.signingNode.RunSigning(ctx, apieth.MPCSigningNodeParams{
		SessionID:   sessionInfo.SessionID,
		Hash:        sessionInfo.Hash,
		PartyID:     input.PartyID,
		AllPartyIDs: sessionInfo.PartyIDs,
		Threshold:   sessionInfo.Threshold,
		ShardJSON:   shardJSON,
	})
}

// verifySigningHash decodes rawTxHex and verifies that the EIP-1559 signing pre-image
// (signer.Hash) matches the expectedHash provided by the coordinator. This prevents
// blind signing: a compromised coordinator cannot make the node sign an arbitrary hash.
func verifySigningHash(rawTxHex string, expectedHash []byte) error {
	tx, err := pkgeth.DecodeTx(rawTxHex)
	if err != nil {
		return fmt.Errorf("decode raw tx: %w", err)
	}
	chainID := tx.ChainId()
	signer := types.LatestSignerForChainID(chainID)
	computed := signer.Hash(tx)
	if !bytes.Equal(computed.Bytes(), expectedHash) {
		return fmt.Errorf("hash mismatch: coordinator sent %x but raw_tx hashes to %x",
			expectedHash, computed.Bytes())
	}
	return nil
}
