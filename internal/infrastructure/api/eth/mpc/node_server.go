package mpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// Compile-time interface check.
var _ apieth.MPCNodeDeps = (*MPCNodeServer)(nil)

// MPCNodeServer implements MPCKeyGeneratorPort, providing DKG ceremony execution and
// Paillier pre-parameter generation for MPC node operators.
//
// It drives the tss-lib/v2 keygen state machine using messages delivered by the
// MPCInboundTransport (via the Watch wallet coordinator relay). Outbound round messages
// are marshalled into MPCWireMessage envelopes and enqueued back to the coordinator.
type MPCNodeServer struct {
	transport    apieth.MPCInboundTransport
	shardStorage apieth.MPCKeyShardStorage
	logger       *slog.Logger
}

// NewMPCNodeServer creates an MPCNodeServer ready to run DKG ceremonies and generate
// pre-parameters.
func NewMPCNodeServer(
	transport apieth.MPCInboundTransport,
	shardStorage apieth.MPCKeyShardStorage,
	logger *slog.Logger,
) *MPCNodeServer {
	if logger == nil {
		logger = slog.Default()
	}
	return &MPCNodeServer{
		transport:    transport,
		shardStorage: shardStorage,
		logger:       logger,
	}
}

// =============================================================================
// MPCKeyGeneratorPort implementation
// =============================================================================

// GeneratePreParams generates the Paillier key pre-computation parameters required for DKG
// and writes them to outputPath as JSON. The operation is CPU-intensive (minutes of wall
// time) and respects ctx cancellation.
//
// The output file is world-owner-readable only (mode 0600) to protect key material.
func (*MPCNodeServer) GeneratePreParams(ctx context.Context, outputPath string) error {
	preParams, err := keygen.GeneratePreParamsWithContext(ctx)
	if err != nil {
		return fmt.Errorf("mpc node: generate pre-params: %w", err)
	}

	data, err := json.Marshal(preParams)
	if err != nil {
		return fmt.Errorf("mpc node: marshal pre-params: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("mpc node: write pre-params to %s: %w", outputPath, err)
	}
	return nil
}

// RunDKG runs the DKG ceremony for this node. It loads pre-computed Paillier parameters,
// creates a keygen.LocalParty, and drives all keygen rounds by exchanging messages with
// peer nodes.
//
// When params.ListenAddr is non-empty, P2P mode is used: this node starts a gRPC server on
// ListenAddr and connects directly to each peer in PeerAddrs (index-matched to AllPartyIDs
// minus self). This allows nodes to communicate without an external coordinator.
//
// When params.ListenAddr is empty, coordinator relay mode is used (backward-compatible):
// outbound messages are enqueued via EnqueueOutbound for delivery by the Watch wallet.
//
// On success the joint ECDSA public key is computed, the joint Ethereum address is derived,
// and the encrypted shard is written to params.ShardOutputPath via MPCKeyShardStorage.
//
// RunDKG blocks until all rounds complete, ctx is cancelled, or an error occurs.
// If DKG is aborted by any peer, RunDKG returns an error and does NOT save a partial shard.
func (s *MPCNodeServer) RunDKG(ctx context.Context, params apieth.DKGParams) (*apieth.DKGResult, error) {
	// 1. Build a deterministic, sorted party ID set from the string IDs.
	//    Each party key is derived as keccak256(partyID) so it is stable across restarts.
	//    Done first so we can fail fast on invalid party configuration before loading key material.
	sorted := buildSortedPartyIDs(params.AllPartyIDs)

	// 2. Locate this node's own PartyID in the sorted list.
	myPartyID, err := findPartyID(sorted, params.PartyID)
	if err != nil {
		return nil, fmt.Errorf("mpc node: %w", err)
	}

	// 3. Load and unmarshal pre-params.
	rawPreParams, err := os.ReadFile(params.PreParamsPath)
	if err != nil {
		return nil, fmt.Errorf("mpc node: read pre_params from %s: %w", params.PreParamsPath, err)
	}
	var preParams keygen.LocalPreParams
	if err := json.Unmarshal(rawPreParams, &preParams); err != nil {
		return nil, fmt.Errorf("mpc node: unmarshal pre_params: %w", err)
	}

	// 4. Validate pre-params before passing to tss-lib to avoid a panic inside NewLocalParty.
	if !preParams.ValidateWithProof() {
		return nil, fmt.Errorf("mpc node: pre-params failed validation; "+
			"re-generate with 'keygen pre-params' or check the file at %s", params.PreParamsPath)
	}

	// 5. Create tss Parameters.
	//    Note: tss-lib's threshold is T-1 (i.e. shares needed beyond one), matching GG18.
	//    A 2-of-2 setup uses threshold=1; a 2-of-3 setup uses threshold=1 as well.
	peerCtx := tss.NewPeerContext(sorted)
	tssParams := tss.NewParameters(tss.EC(), peerCtx, myPartyID, len(sorted), params.Threshold-1)

	// 6. Initialise keygen channels and party.
	outCh := make(chan tss.Message, len(sorted)*10)
	endCh := make(chan *keygen.LocalPartySaveData, 1)
	party := keygen.NewLocalParty(tssParams, outCh, endCh, preParams)

	if tssErr := party.Start(); tssErr != nil {
		return nil, fmt.Errorf("mpc node: start keygen party: %w", tssErr.Cause())
	}

	// 7. Build a lookup map for inbound message routing.
	partyByID := buildPartyByIDMap(sorted)

	// 8. P2P mode: start gRPC server and connect to peers directly.
	if params.ListenAddr != "" {
		return s.runDKGP2P(ctx, party, partyByID, outCh, endCh, params)
	}

	// 9. Coordinator relay mode (backward-compatible): get the inbound channel.
	recvCh, err := s.transport.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("mpc node: get receive channel: %w", err)
	}

	// 10. Run the DKG message loop with coordinator relay.
	return s.runDKGLoop(ctx, party, partyByID, outCh, endCh, recvCh, params)
}

// runDKGP2P runs the DKG ceremony in peer-to-peer mode.
//
// Each node starts a gRPC server on params.ListenAddr, waits for peer servers to be ready,
// then connects to all peers via GRPCOutboundTransport. TSS round messages are routed
// directly to the appropriate peer(s) without an external coordinator.
func (s *MPCNodeServer) runDKGP2P(
	ctx context.Context,
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	outCh <-chan tss.Message,
	endCh <-chan *keygen.LocalPartySaveData,
	params apieth.DKGParams,
) (*apieth.DKGResult, error) {
	// Start the gRPC server so peers can connect and deliver messages to this node.
	if err := s.transport.Listen(ctx, params.ListenAddr); err != nil {
		return nil, fmt.Errorf("mpc node: P2P listen on %s: %w", params.ListenAddr, err)
	}
	defer func() { _ = s.transport.Close() }()

	// Build a map from party ID → peer gRPC address.
	// PeerAddrs is index-matched with AllPartyIDs (self entry skipped).
	peerAddrByPartyID := buildPeerAddrMap(params.AllPartyIDs, params.PartyID, params.PeerAddrs)

	// Wait for all peer gRPC servers to be reachable before connecting.
	// This avoids "connection refused" races when nodes start concurrently.
	const peerReadyTimeout = 90 * time.Second
	for pid, addr := range peerAddrByPartyID {
		if err := waitForTCPAddr(ctx, addr, peerReadyTimeout); err != nil {
			return nil, fmt.Errorf("mpc node: peer %s (%s) not ready: %w", pid, addr, err)
		}
	}

	// Connect to all peers. Use a fixed DKG session ID so that each peer's gRPC inbound
	// transport accepts messages without requiring a signed signing request.
	const dkgSessionID = "dkg-ceremony"
	initReq := (&protogen.MPCSigningRequest_builder{SessionId: dkgSessionID}).Build()
	peerTransport := NewGRPCOutboundTransportForSession(dkgSessionID, initReq)
	defer func() { _ = peerTransport.Close() }()

	for pid, addr := range peerAddrByPartyID {
		if err := peerTransport.Connect(ctx, addr); err != nil {
			return nil, fmt.Errorf("mpc node: connect to DKG peer %s (%s): %w", pid, addr, err)
		}
	}

	// Receive inbound messages from peers via this node's gRPC server.
	recvCh, err := s.transport.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("mpc node: P2P get receive channel: %w", err)
	}

	return s.runDKGLoopP2P(ctx, party, partyByID, outCh, endCh, recvCh, params, peerTransport, peerAddrByPartyID)
}

// runDKGLoopP2P drives the tss-lib keygen state machine in peer-to-peer mode.
//
// Outbound TSS messages are routed directly to peer nodes via peerTransport.
// Inbound messages arrive via recvCh (delivered by peers connecting to this node's gRPC server).
func (s *MPCNodeServer) runDKGLoopP2P(
	ctx context.Context,
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	outCh <-chan tss.Message,
	endCh <-chan *keygen.LocalPartySaveData,
	recvCh <-chan []byte,
	params apieth.DKGParams,
	peerTransport *GRPCOutboundTransport,
	peerAddrByPartyID map[string]string,
) (*apieth.DKGResult, error) {
	sendFn := func(ctx context.Context, msg tss.Message) error {
		return s.sendToPeers(ctx, msg, peerTransport, peerAddrByPartyID)
	}
	return s.runDKGMessageLoop(ctx, party, partyByID, outCh, endCh, recvCh, params, sendFn)
}

// sendToPeers routes a TSS outbound message to the appropriate peer node(s).
// Broadcast messages are sent to all peers; unicast messages go to the specified party.
func (s *MPCNodeServer) sendToPeers(
	ctx context.Context,
	msg tss.Message,
	transport *GRPCOutboundTransport,
	peerAddrByPartyID map[string]string,
) error {
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

	if routing.IsBroadcast || len(to) == 0 {
		// Broadcast: send to all peers.
		for pid, addr := range peerAddrByPartyID {
			if err := transport.Send(ctx, addr, data); err != nil {
				return fmt.Errorf("broadcast to peer %s (%s): %w", pid, addr, err)
			}
		}
		return nil
	}

	// Unicast / multicast: send only to specified parties.
	for _, pid := range to {
		addr, ok := peerAddrByPartyID[pid]
		if !ok {
			s.logger.Warn("mpc node: unknown target party ID, skipping", "to", pid)
			continue
		}
		if err := transport.Send(ctx, addr, data); err != nil {
			return fmt.Errorf("unicast to peer %s (%s): %w", pid, addr, err)
		}
	}
	return nil
}

// buildPeerAddrMap constructs a map from party ID to gRPC address for all peers of selfID.
// PeerAddrs is index-matched with AllPartyIDs after removing selfID.
func buildPeerAddrMap(allPartyIDs []string, selfID string, peerAddrs []string) map[string]string {
	result := make(map[string]string, len(peerAddrs))
	idx := 0
	for _, pid := range allPartyIDs {
		if pid == selfID {
			continue
		}
		if idx < len(peerAddrs) {
			result[pid] = peerAddrs[idx]
		}
		idx++
	}
	return result
}

// waitForTCPAddr polls addr until a TCP connection can be established or timeout expires.
func waitForTCPAddr(ctx context.Context, addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("address %s not reachable after %s: %w", addr, timeout, err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

// runDKGLoop drives the tss-lib keygen state machine in coordinator relay mode.
//
// It selects over three channels simultaneously:
//   - outCh: TSS library has produced an outbound round message → encode and send to coordinator.
//   - recvCh: coordinator has delivered an inbound round message → parse and update party.
//   - endCh: DKG is complete → save shard and return result.
//   - ctx.Done(): context expired → abort and return error.
func (s *MPCNodeServer) runDKGLoop(
	ctx context.Context,
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	outCh <-chan tss.Message,
	endCh <-chan *keygen.LocalPartySaveData,
	recvCh <-chan []byte,
	params apieth.DKGParams,
) (*apieth.DKGResult, error) {
	sendFn := func(_ context.Context, msg tss.Message) error {
		return s.sendOutbound(msg)
	}
	return s.runDKGMessageLoop(ctx, party, partyByID, outCh, endCh, recvCh, params, sendFn)
}

// runDKGMessageLoop is the shared DKG event loop used by both coordinator relay and P2P modes.
// sendFn is called for every outbound TSS round message and encapsulates transport-specific routing.
func (s *MPCNodeServer) runDKGMessageLoop(
	ctx context.Context,
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	outCh <-chan tss.Message,
	endCh <-chan *keygen.LocalPartySaveData,
	recvCh <-chan []byte,
	params apieth.DKGParams,
	sendFn func(ctx context.Context, msg tss.Message) error,
) (*apieth.DKGResult, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("mpc node: DKG cancelled: %w", ctx.Err())

		case msg := <-outCh:
			if err := sendFn(ctx, msg); err != nil {
				return nil, fmt.Errorf("mpc node: send DKG message: %w", err)
			}

		case payload, ok := <-recvCh:
			if !ok {
				return nil, errors.New("mpc node: inbound channel closed unexpectedly")
			}
			s.handleInbound(party, partyByID, payload)

		case saveData := <-endCh:
			return s.finalizeDKG(ctx, saveData, params)
		}
	}
}

// sendOutbound encodes an outbound tss.Message into an MPCWireMessage and enqueues it via
// the transport for delivery to the coordinator.
func (s *MPCNodeServer) sendOutbound(msg tss.Message) error {
	wireBytes, routing, err := msg.WireBytes()
	if err != nil {
		return fmt.Errorf("get wire bytes: %w", err)
	}

	// Collect target party IDs (nil means broadcast to all).
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

// handleInbound decodes an inbound MPCWireMessage from the coordinator, resolves the
// sender's *tss.PartyID, and feeds the message to the local party's state machine.
func (s *MPCNodeServer) handleInbound(
	party tss.Party,
	partyByID map[string]*tss.PartyID,
	payload []byte,
) {
	var wm MPCWireMessage
	if err := json.Unmarshal(payload, &wm); err != nil {
		s.logger.Warn("mpc node: malformed inbound wire message, skipping", "err", err)
		return
	}

	fromParty, ok := partyByID[wm.From]
	if !ok {
		s.logger.Warn("mpc node: unknown sender party ID, skipping", "from", wm.From)
		return
	}

	parsed, err := tss.ParseWireMessage(wm.Data, fromParty, wm.IsBroadcast)
	if err != nil {
		s.logger.Warn("mpc node: parse wire message failed, skipping", "err", err)
		return
	}

	if updated, tssErr := party.Update(parsed); !updated {
		s.logger.Warn("mpc node: party update returned error", "err", tssErr)
	}
}

// finalizeDKG computes the joint Ethereum address, saves the encrypted shard, and returns
// the DKGResult. It does NOT write a partial shard — if any step fails the shard is not saved.
func (s *MPCNodeServer) finalizeDKG(
	ctx context.Context,
	saveData *keygen.LocalPartySaveData,
	params apieth.DKGParams,
) (*apieth.DKGResult, error) {
	if saveData.ECDSAPub == nil {
		return nil, errors.New("mpc node: DKG save data missing ECDSA public key")
	}

	// Derive the Ethereum address from the joint public key.
	ecPub := &ecdsa.PublicKey{
		Curve: saveData.ECDSAPub.Curve(),
		X:     saveData.ECDSAPub.X(),
		Y:     saveData.ECDSAPub.Y(),
	}
	ethAddr := ethcrypto.PubkeyToAddress(*ecPub)

	// Serialise save data to JSON for encrypted storage.
	saveJSON, err := json.Marshal(saveData)
	if err != nil {
		return nil, fmt.Errorf("mpc node: marshal DKG save data: %w", err)
	}

	// Encrypt and write shard. Only called on success — no partial shard is written on error.
	if err := s.shardStorage.SaveShard(ctx, params.ShardOutputPath, params.Passphrase, saveJSON); err != nil {
		return nil, fmt.Errorf("mpc node: save shard to %s: %w", params.ShardOutputPath, err)
	}

	// Encode uncompressed public key: 0x04 || X[32] || Y[32].
	pubBytes := make([]byte, 65)
	pubBytes[0] = 0x04
	xb := saveData.ECDSAPub.X().Bytes()
	yb := saveData.ECDSAPub.Y().Bytes()
	copy(pubBytes[1+32-len(xb):33], xb)
	copy(pubBytes[33+32-len(yb):65], yb)

	return &apieth.DKGResult{
		EthAddress: ethAddr.Hex(),
		PublicKey:  pubBytes,
		ShardPath:  params.ShardOutputPath,
	}, nil
}

// =============================================================================
// Party ID helpers
// =============================================================================

// buildSortedPartyIDs converts string party IDs to tss.SortedPartyIDs using a deterministic
// key derived from keccak256(partyID), ensuring stable ordering across nodes.
func buildSortedPartyIDs(allPartyIDs []string) tss.SortedPartyIDs {
	unsorted := make(tss.UnSortedPartyIDs, len(allPartyIDs))
	for i, pid := range allPartyIDs {
		key := new(big.Int).SetBytes(ethcrypto.Keccak256([]byte(pid)))
		unsorted[i] = tss.NewPartyID(pid, pid, key)
	}
	return tss.SortPartyIDs(unsorted)
}

// findPartyID locates a *tss.PartyID in sorted by its string ID.
func findPartyID(sorted tss.SortedPartyIDs, partyID string) (*tss.PartyID, error) {
	for _, pid := range sorted {
		if pid.Id == partyID {
			return pid, nil
		}
	}
	return nil, fmt.Errorf("party ID %q not found in all_party_ids", partyID)
}

// buildPartyByIDMap returns a map from party ID string to *tss.PartyID for O(1) lookup
// when routing inbound messages.
func buildPartyByIDMap(sorted tss.SortedPartyIDs) map[string]*tss.PartyID {
	m := make(map[string]*tss.PartyID, len(sorted))
	for _, pid := range sorted {
		m[pid.Id] = pid
	}
	return m
}
