package mpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

// Compile-time interface check.
var _ apieth.MPCSigningNodePort = (*MPCSigningNode)(nil)

// MPCSigningNode implements MPCSigningNodePort. It drives the tss-lib/v2 ECDSA signing
// state machine, exchanging round messages with peer nodes via the coordinator relay.
type MPCSigningNode struct {
	transport apieth.MPCInboundTransport
	logger    *slog.Logger
}

// NewMPCSigningNode creates an MPCSigningNode ready to run TSS signing sessions.
func NewMPCSigningNode(
	transport apieth.MPCInboundTransport,
	logger *slog.Logger,
) *MPCSigningNode {
	if logger == nil {
		logger = slog.Default()
	}
	return &MPCSigningNode{
		transport: transport,
		logger:    logger,
	}
}

// RunSigning runs the TSS ECDSA signing protocol for the session described by params.
//
// It unmarshals the key shard, builds the TSS signing party, and drives all signing rounds
// by exchanging MPCWireMessages with the coordinator relay via the inbound transport.
// When the TSS protocol completes, the assembled 65-byte ECDSA signature is enqueued
// back to the coordinator with IsSignature=true.
func (s *MPCSigningNode) RunSigning(ctx context.Context, params apieth.MPCSigningNodeParams) error {
	// 1. Unmarshal the key shard from DKG.
	var saveData keygen.LocalPartySaveData
	if err := json.Unmarshal(params.ShardJSON, &saveData); err != nil {
		return fmt.Errorf("mpc signing node: unmarshal shard JSON: %w", err)
	}

	// 2. Build sorted party IDs from the T signing participants.
	sorted := buildSortedPartyIDs(params.AllPartyIDs)

	// 3. Locate this node's own PartyID in the sorted list.
	myPartyID, err := findPartyID(sorted, params.PartyID)
	if err != nil {
		return fmt.Errorf("mpc signing node: %w", err)
	}

	// 4. Create tss.Parameters.
	//    tss-lib uses threshold T-1 (GG18 convention: shares needed beyond one).
	peerCtx := tss.NewPeerContext(sorted)
	tssParams := tss.NewParameters(tss.EC(), peerCtx, myPartyID, len(sorted), params.Threshold-1)

	// 5. Convert the 32-byte hash to *big.Int for the signing party.
	msgBig := new(big.Int).SetBytes(params.Hash)

	// 6. Initialize signing channels and party.
	outCh := make(chan tss.Message, len(sorted)*10)
	endCh := make(chan *common.SignatureData, 1)
	party := signing.NewLocalParty(msgBig, tssParams, saveData, outCh, endCh)

	if tssErr := party.Start(); tssErr != nil {
		return fmt.Errorf("mpc signing node: start signing party: %w", tssErr.Cause())
	}

	s.logger.Info("mpc signing node: started TSS signing",
		"session", params.SessionID,
		"partyID", params.PartyID,
		"threshold", params.Threshold,
	)

	// 7. Get the inbound message channel from the transport.
	recvCh, err := s.transport.Receive(ctx)
	if err != nil {
		return fmt.Errorf("mpc signing node: get receive channel: %w", err)
	}

	// 8. Build a lookup map for inbound message routing.
	partyByID := buildPartyByIDMap(sorted)

	// 9. Run the signing message loop.
	return s.runSigningLoop(ctx, party, partyByID, outCh, endCh, recvCh, params.SessionID)
}

// runSigningLoop drives the TSS signing state machine. It selects over four channels:
//   - outCh: TSS round message ready to send → encode and enqueue to coordinator.
//   - recvCh: inbound round message from coordinator → parse and update party.
//   - endCh: signing complete → assemble 65-byte signature, enqueue, return.
//   - ctx.Done(): context cancelled → abort.
func (s *MPCSigningNode) runSigningLoop(
	ctx context.Context,
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	outCh <-chan tss.Message,
	endCh <-chan *common.SignatureData,
	recvCh <-chan []byte,
	sessionID string,
) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("mpc signing node [%s]: context cancelled: %w", sessionID, ctx.Err())

		case msg := <-outCh:
			if err := s.sendOutbound(msg); err != nil {
				return fmt.Errorf("mpc signing node [%s]: send outbound: %w", sessionID, err)
			}

		case payload, ok := <-recvCh:
			if !ok {
				return errors.New("mpc signing node: inbound channel closed unexpectedly")
			}
			s.handleInbound(party, partyByID, payload)

		case sigData := <-endCh:
			return s.finalizeSignature(sigData, sessionID)
		}
	}
}

// sendOutbound encodes an outbound TSS message into an MPCWireMessage and enqueues it via
// the transport for delivery to the coordinator.
func (s *MPCSigningNode) sendOutbound(msg tss.Message) error {
	wireBytes, routing, err := msg.WireBytes()
	if err != nil {
		return fmt.Errorf("get wire bytes: %w", err)
	}

	var to []string
	for _, pid := range routing.To {
		to = append(to, pid.Id)
	}

	wm := MPCWireMessage{
		From:        msg.GetFrom().Id,
		To:          to,
		IsBroadcast: routing.IsBroadcast,
		Data:        wireBytes,
	}
	data, err := json.Marshal(wm)
	if err != nil {
		return fmt.Errorf("marshal wire message: %w", err)
	}
	return s.transport.EnqueueOutbound(data)
}

// handleInbound decodes an inbound MPCWireMessage, resolves the sender's *tss.PartyID,
// and feeds the message to the local party's signing state machine.
func (s *MPCSigningNode) handleInbound(
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	payload []byte,
) {
	var wm MPCWireMessage
	if err := json.Unmarshal(payload, &wm); err != nil {
		s.logger.Warn("mpc signing node: malformed inbound wire message, skipping", "err", err)
		return
	}

	fromParty, ok := partyByID[wm.From]
	if !ok {
		s.logger.Warn("mpc signing node: unknown sender party ID, skipping", "from", wm.From)
		return
	}

	parsed, err := tss.ParseWireMessage(wm.Data, fromParty, wm.IsBroadcast)
	if err != nil {
		s.logger.Warn("mpc signing node: parse wire message failed, skipping", "err", err)
		return
	}

	if updated, tssErr := party.Update(parsed); !updated {
		s.logger.Warn("mpc signing node: party update returned error", "err", tssErr)
	}
}

// finalizeSignature assembles the 65-byte ECDSA signature from the TSS result and
// sends it to the coordinator as a terminal MPCWireMessage with IsSignature=true.
func (s *MPCSigningNode) finalizeSignature(sigData *common.SignatureData, sessionID string) error {
	sig := make([]byte, 65)
	// R and S from tss-lib are not zero-padded (big.Int.Bytes() omits leading zeros).
	// Use FillBytes to ensure each component is exactly 32 bytes.
	new(big.Int).SetBytes(sigData.R).FillBytes(sig[0:32])
	new(big.Int).SetBytes(sigData.S).FillBytes(sig[32:64])
	if len(sigData.SignatureRecovery) > 0 {
		sig[64] = sigData.SignatureRecovery[0]
	}

	wm := MPCWireMessage{
		IsSignature: true,
		Data:        sig,
	}
	data, err := json.Marshal(wm)
	if err != nil {
		return fmt.Errorf("mpc signing node [%s]: marshal signature: %w", sessionID, err)
	}

	s.logger.Info("mpc signing node: TSS signing complete, sending signature", "session", sessionID)
	return s.transport.EnqueueOutbound(data)
}
