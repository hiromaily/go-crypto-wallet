// Package mpc provides gRPC-based transport implementations for the MPC/TSS signing relay.
package mpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// Compile-time interface check.
var _ apieth.MPCCoordinatorDeps = (*MPCCoordinator)(nil)

// MPCCoordinator implements MPCTransactionSigner as the Watch-side TSS signing coordinator.
//
// The coordinator is a passive relay participant — it never holds or reconstructs the private key.
// It contacts all T participating MPC nodes, fans out the InitSigning request, and then acts as
// a message bus: routing TSS round messages between nodes until one node reports the assembled
// 65-byte ECDSA signature. A fresh GRPCOutboundTransport is created per signing session so
// that session state is fully isolated.
type MPCCoordinator struct {
	logger *slog.Logger
}

// NewMPCCoordinator creates a coordinator ready to handle TSS signing sessions.
func NewMPCCoordinator(logger *slog.Logger) *MPCCoordinator {
	return &MPCCoordinator{logger: logger}
}

// SignTransaction initiates a TSS signing session for req.Hash, contacts all T nodes listed in
// req.PeerAddrs, relays TSS round messages between them, and returns the assembled 65-byte ECDSA
// signature once a node signals completion.
//
// Validation:
//   - req.Hash must be exactly 32 bytes.
//   - len(req.PartyIDs) must be >= req.Threshold.
//   - req.PartyIDs and req.PeerAddrs must have the same length.
//
// The method blocks until the signature is received or ctx is cancelled. It returns a wrapped
// error containing the session ID on any failure so that callers can correlate log entries.
func (c *MPCCoordinator) SignTransaction(
	ctx context.Context,
	req apieth.MPCSigningRequest,
) (*apieth.MPCSigningResult, error) {
	// 1. Validate request fields.
	if len(req.Hash) != 32 {
		return nil, fmt.Errorf(
			"mpc coordinator [%s]: hash must be 32 bytes, got %d", req.SessionID, len(req.Hash),
		)
	}
	if len(req.PartyIDs) < req.Threshold {
		return nil, fmt.Errorf(
			"mpc coordinator [%s]: insufficient parties: got %d, threshold %d",
			req.SessionID, len(req.PartyIDs), req.Threshold,
		)
	}
	if len(req.PartyIDs) != len(req.PeerAddrs) {
		return nil, fmt.Errorf(
			"mpc coordinator [%s]: party_ids and peer_addrs must have the same length (%d != %d)",
			req.SessionID, len(req.PartyIDs), len(req.PeerAddrs),
		)
	}

	// 2. Build party-ID → gRPC-address lookup for message routing.
	partyToAddr := make(map[string]string, len(req.PartyIDs))
	for i, pid := range req.PartyIDs {
		partyToAddr[pid] = req.PeerAddrs[i]
	}

	// 3. Create a per-session transport that calls InitSigning on each peer before opening
	//    the RelaySession stream.  This registers the session ID on each peer's inbound
	//    transport so that relay messages are not silently discarded by the session ID filter.
	initReq := (&protogen.MPCSigningRequest_builder{
		SessionId: req.SessionID,
		Hash:      req.Hash,
		PartyIds:  req.PartyIDs,
		Threshold: int32(req.Threshold),
		RawTxHex:  req.RawTxHex,
	}).Build()

	transport := NewGRPCOutboundTransportForSession(req.SessionID, initReq)
	defer func() {
		if err := transport.Close(); err != nil {
			c.logger.Warn("mpc coordinator: transport close error", "session", req.SessionID, "err", err)
		}
	}()

	// 4. Connect to all peers. This triggers InitSigning + RelaySession on each peer so that
	//    their background receive goroutines start before the relay loop begins waiting.
	if err := c.connectAll(ctx, transport, req.PeerAddrs, req.SessionID); err != nil {
		return nil, err
	}

	// 5. Run the relay loop until a node reports the signature or ctx expires.
	return c.relayLoop(ctx, req, transport, partyToAddr)
}

// connectAll opens connections to all peers sequentially. Any connection failure aborts the
// session with a wrapped error.
func (*MPCCoordinator) connectAll(
	ctx context.Context,
	transport *GRPCOutboundTransport,
	peerAddrs []string,
	sessionID string,
) error {
	for _, addr := range peerAddrs {
		if err := transport.Connect(ctx, addr); err != nil {
			return fmt.Errorf("mpc coordinator [%s]: connect to %s: %w", sessionID, addr, err)
		}
	}
	return nil
}

// relayLoop reads MPCWireMessages from the transport receive channel and routes the embedded
// Data bytes to the target peer(s) indicated by To.  When a message with IsSignature == true
// arrives, it validates the signature length and returns the result.
func (c *MPCCoordinator) relayLoop(
	ctx context.Context,
	req apieth.MPCSigningRequest,
	transport *GRPCOutboundTransport,
	partyToAddr map[string]string,
) (*apieth.MPCSigningResult, error) {
	recvCh, err := transport.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("mpc coordinator [%s]: open receive channel: %w", req.SessionID, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"mpc coordinator [%s]: context cancelled: %w", req.SessionID, ctx.Err(),
			)

		case payload, ok := <-recvCh:
			if !ok {
				return nil, fmt.Errorf(
					"mpc coordinator [%s]: receive channel closed unexpectedly", req.SessionID,
				)
			}

			var wm MPCWireMessage
			if err := json.Unmarshal(payload, &wm); err != nil {
				c.logger.Warn("mpc coordinator: malformed wire message, skipping",
					"session", req.SessionID, "err", err)
				continue
			}

			// Signing complete — validate and return.
			if wm.IsSignature {
				if len(wm.Data) != 65 {
					return nil, fmt.Errorf(
						"mpc coordinator [%s]: invalid signature length %d (want 65)",
						req.SessionID, len(wm.Data),
					)
				}
				return &apieth.MPCSigningResult{
					SessionID: req.SessionID,
					Signature: wm.Data,
				}, nil
			}

			// Route TSS round message to target node(s).
			if err := c.routeMessage(ctx, transport, partyToAddr, wm, req.SessionID); err != nil {
				return nil, err
			}
		}
	}
}

// routeMessage forwards wm to the appropriate peer(s).
//
// The coordinator re-encodes a relay MPCWireMessage that preserves From and IsBroadcast
// (so the receiving node can call tss.ParseWireMessage correctly) but clears To (the
// recipient does not need routing metadata).
//
// If wm.To is empty, the message is broadcast to all peers except wm.From.
// If wm.To is non-empty, the message is sent only to the listed party IDs.
func (c *MPCCoordinator) routeMessage(
	ctx context.Context,
	transport *GRPCOutboundTransport,
	partyToAddr map[string]string,
	wm MPCWireMessage,
	sessionID string,
) error {
	// Build relay envelope: preserve sender and broadcast flag so the node can parse
	// the embedded TSS wire bytes with tss.ParseWireMessage(data, from, isBroadcast).
	relay := MPCWireMessage{
		From:        wm.From,
		IsBroadcast: wm.IsBroadcast,
		Data:        wm.Data,
	}
	relayData, err := json.Marshal(relay)
	if err != nil {
		return fmt.Errorf("mpc coordinator [%s]: marshal relay message: %w", sessionID, err)
	}

	if len(wm.To) == 0 {
		// Broadcast: send to all parties except the sender.
		for pid, addr := range partyToAddr {
			if pid == wm.From {
				continue
			}
			if err := transport.Send(ctx, addr, relayData); err != nil {
				return fmt.Errorf(
					"mpc coordinator [%s]: broadcast to %s: %w", sessionID, addr, err,
				)
			}
		}
		return nil
	}

	// Unicast / multicast to specified parties.
	for _, pid := range wm.To {
		addr, ok := partyToAddr[pid]
		if !ok {
			c.logger.Warn("mpc coordinator: unknown party ID in To field, skipping",
				"session", sessionID, "party", pid)
			continue
		}
		if err := transport.Send(ctx, addr, relayData); err != nil {
			return fmt.Errorf(
				"mpc coordinator [%s]: unicast to %s: %w", sessionID, addr, err,
			)
		}
	}
	return nil
}
