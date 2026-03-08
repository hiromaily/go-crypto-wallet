package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ErrTxHashMismatch is returned when the hash of the raw transaction in the signing request
// does not match the tx_hash field. This indicates a tampered or corrupt request.
var ErrTxHashMismatch = errors.New("tx_hash does not match hash(raw_tx_hex): signing request rejected")

type serveMPCUseCase struct {
	shardStorage apieth.MPCKeyShardStorage
	transport    apieth.MPCInboundTransport
}

// NewServeMPCUseCase constructs the use case for running the MPC node gRPC signing daemon.
func NewServeMPCUseCase(
	shardStorage apieth.MPCKeyShardStorage,
	transport apieth.MPCInboundTransport,
) keygenusecase.ServeMPCUseCase {
	return &serveMPCUseCase{
		shardStorage: shardStorage,
		transport:    transport,
	}
}

// Serve starts the MPC node daemon. It loads the key shard, starts the inbound transport
// listener, waits for the first signing session request, verifies the tx hash, and then
// blocks until ctx is cancelled.
//
// Returns ErrTxHashMismatch if the coordinator sends a signing request where
// hash(raw_tx_hex) != tx_hash.
func (u *serveMPCUseCase) Serve(ctx context.Context, input keygenusecase.ServeMPCInput) error {
	// 1. Load key shard at startup.
	_, err := u.shardStorage.LoadShard(ctx, input.ShardPath, input.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to load key shard from %s: %w", input.ShardPath, err)
	}

	// 2. Start gRPC listener.
	if err = u.transport.Listen(ctx, input.ListenAddr); err != nil {
		return fmt.Errorf("failed to start MPC inbound listener on %s: %w", input.ListenAddr, err)
	}
	defer func() { _ = u.transport.Close() }()

	logger.Info("MPC node server started", "listenAddr", input.ListenAddr, "partyID", input.PartyID)

	// 3. Get inbound message channel.
	msgCh, err := u.transport.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to open receive channel: %w", err)
	}

	// 4. Wait for the first message (session request) and verify the tx hash.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg, ok := <-msgCh:
		if !ok {
			return nil // channel closed before any session
		}
		if verifyErr := u.verifySessionRequest(msg); verifyErr != nil {
			return verifyErr
		}
	}

	// 5. Block until context is cancelled (graceful shutdown).
	<-ctx.Done()
	return ctx.Err()
}

// verifySessionRequest unmarshals a session init message and verifies that
// hash(RawTxHex) matches TxHash. Returns ErrTxHashMismatch on mismatch.
func (*serveMPCUseCase) verifySessionRequest(msg []byte) error {
	var req apieth.MPCSessionRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		return fmt.Errorf("failed to unmarshal session request: %w", err)
	}

	tx, err := pkgeth.DecodeTx(req.RawTxHex)
	if err != nil {
		return fmt.Errorf("failed to decode raw transaction: %w", err)
	}

	computedHash := tx.Hash().Hex()
	if computedHash != req.TxHash {
		return fmt.Errorf("%w: session=%s computed=%s provided=%s",
			ErrTxHashMismatch, req.SessionID, computedHash, req.TxHash)
	}

	logger.Info("MPC signing session request verified",
		"sessionID", req.SessionID,
		"txHash", req.TxHash,
	)
	return nil
}
