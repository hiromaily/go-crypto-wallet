package mpc_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	mpcinfra "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// =============================================================================
// Fake node server used by coordinator tests
// =============================================================================

// coordinatorTestNode is a minimal MPCNodeService gRPC server for coordinator tests.
//
// On RelaySession, it immediately sends each queued outbound message from toSend,
// then drains received messages into the received channel.
type coordinatorTestNode struct {
	protogen.UnimplementedMPCNodeServiceServer

	mu       sync.Mutex
	toSend   []mpcinfra.MPCWireMessage // queued messages to push to coordinator
	received chan []byte               // raw payloads received from coordinator
}

func newCoordinatorTestNode(toSend ...mpcinfra.MPCWireMessage) *coordinatorTestNode {
	return &coordinatorTestNode{
		toSend:   toSend,
		received: make(chan []byte, 64),
	}
}

func (*coordinatorTestNode) InitSigning(
	_ context.Context,
	req *protogen.MPCSigningRequest,
) (*protogen.MPCSigningResponse, error) {
	return (&protogen.MPCSigningResponse_builder{SessionId: req.GetSessionId()}).Build(), nil
}

func (n *coordinatorTestNode) RelaySession(
	stream grpc.BidiStreamingServer[protogen.MPCMessage, protogen.MPCMessage],
) error {
	n.mu.Lock()
	pending := n.toSend
	n.mu.Unlock()

	// Push pre-queued messages to the coordinator immediately.
	for _, wm := range pending {
		data, err := json.Marshal(wm)
		if err != nil {
			return err
		}
		if err := stream.Send((&protogen.MPCMessage_builder{Payload: data}).Build()); err != nil {
			return err
		}
	}

	// Drain inbound messages (routed from coordinator) so the stream doesn't block.
	for {
		msg, err := stream.Recv()
		if err != nil {
			return nil
		}
		n.received <- msg.GetPayload()
	}
}

// startCoordinatorTestNode starts an in-process gRPC server for coordinatorTestNode
// and returns its address.
func startCoordinatorTestNode(t *testing.T, node *coordinatorTestNode) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	s := grpc.NewServer()
	protogen.RegisterMPCNodeServiceServer(s, node)
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(func() { s.GracefulStop() })
	return lis.Addr().String()
}

// =============================================================================
// Validation tests (no gRPC required)
// =============================================================================

func TestMPCCoordinator_SignTransaction_Validation(t *testing.T) {
	t.Parallel()

	coord := mpcinfra.NewMPCCoordinator(slog.Default())
	ctx := context.Background()

	tests := []struct {
		name    string
		req     apieth.MPCSigningRequest
		wantErr string
	}{
		{
			name: "hash too short",
			req: apieth.MPCSigningRequest{
				SessionID: "sess-1",
				Hash:      make([]byte, 16),
				PartyIDs:  []string{"p1", "p2"},
				PeerAddrs: []string{"addr1", "addr2"},
				Threshold: 2,
			},
			wantErr: "hash must be 32 bytes",
		},
		{
			name: "hash too long",
			req: apieth.MPCSigningRequest{
				SessionID: "sess-2",
				Hash:      make([]byte, 64),
				PartyIDs:  []string{"p1", "p2"},
				PeerAddrs: []string{"addr1", "addr2"},
				Threshold: 2,
			},
			wantErr: "hash must be 32 bytes",
		},
		{
			name: "insufficient parties",
			req: apieth.MPCSigningRequest{
				SessionID: "sess-3",
				Hash:      make([]byte, 32),
				PartyIDs:  []string{"p1"},
				PeerAddrs: []string{"addr1"},
				Threshold: 2,
			},
			wantErr: "insufficient parties",
		},
		{
			name: "party_ids and peer_addrs length mismatch",
			req: apieth.MPCSigningRequest{
				SessionID: "sess-4",
				Hash:      make([]byte, 32),
				PartyIDs:  []string{"p1", "p2"},
				PeerAddrs: []string{"addr1"},
				Threshold: 2,
			},
			wantErr: "party_ids and peer_addrs must have the same length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := coord.SignTransaction(ctx, tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =============================================================================
// Integration tests — coordinator relay with fake nodes
// =============================================================================

func TestMPCCoordinator_SignTransaction_CollectsSignature(t *testing.T) {
	t.Parallel()

	// Fake 65-byte signature: r(32) || s(32) || v(1)
	sig := make([]byte, 65)
	sig[64] = 27

	sessionID := "test-session-collect"
	hash := make([]byte, 32)

	// Node 1: immediately sends the signature as soon as the relay session opens.
	node1 := newCoordinatorTestNode(mpcinfra.MPCWireMessage{
		From:        "party-1",
		IsSignature: true,
		Data:        sig,
	})
	// Node 2: passive — receives messages, sends nothing.
	node2 := newCoordinatorTestNode()

	addr1 := startCoordinatorTestNode(t, node1)
	addr2 := startCoordinatorTestNode(t, node2)

	coord := mpcinfra.NewMPCCoordinator(slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := coord.SignTransaction(ctx, apieth.MPCSigningRequest{
		SessionID: sessionID,
		Hash:      hash,
		PartyIDs:  []string{"party-1", "party-2"},
		PeerAddrs: []string{addr1, addr2},
		Threshold: 2,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, sig, result.Signature)
}

func TestMPCCoordinator_SignTransaction_RoutesBroadcast(t *testing.T) {
	t.Parallel()

	sessionID := "test-session-route"
	hash := make([]byte, 32)
	sig := make([]byte, 65)
	sig[64] = 28

	// Node 1: sends a broadcast round message, then (after coordinator routes it to node 2)
	// node 2 will send the signature.
	roundMsg := []byte("tss-round-1-from-party-1")
	node1 := newCoordinatorTestNode(
		// Broadcast round message (To is empty → send to all except sender)
		mpcinfra.MPCWireMessage{
			From: "party-1",
			Data: roundMsg,
		},
	)

	// Node 2: sends the signature after receiving node 1's round message from coordinator.
	node2 := newCoordinatorTestNode()

	addr1 := startCoordinatorTestNode(t, node1)
	addr2 := startCoordinatorTestNode(t, node2)

	// Give node 2 a way to react: it'll send the signature after receiving the routed message.
	// We add a helper goroutine to simulate node 2 responding.
	signalDone := make(chan struct{})
	go func() {
		defer close(signalDone)
		select {
		case payload := <-node2.received:
			// Verify node 2 received the routed round message.
			assert.Equal(t, roundMsg, payload)
		case <-time.After(5 * time.Second):
			return
		}
	}()

	coord := mpcinfra.NewMPCCoordinator(slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// For this test, node 1 sends sig immediately after the round msg to avoid
	// a complex reactive node implementation. We layer both messages in toSend.
	node1.mu.Lock()
	node1.toSend = append(node1.toSend, mpcinfra.MPCWireMessage{
		From:        "party-1",
		IsSignature: true,
		Data:        sig,
	})
	node1.mu.Unlock()

	result, err := coord.SignTransaction(ctx, apieth.MPCSigningRequest{
		SessionID: sessionID,
		Hash:      hash,
		PartyIDs:  []string{"party-1", "party-2"},
		PeerAddrs: []string{addr1, addr2},
		Threshold: 2,
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, sig, result.Signature)

	// Verify the broadcast was routed to node 2.
	select {
	case <-signalDone:
	case <-time.After(3 * time.Second):
		t.Fatal("node 2 did not receive the routed broadcast message in time")
	}
}

func TestMPCCoordinator_SignTransaction_ContextCancellation(t *testing.T) {
	t.Parallel()

	// No node sends anything — coordinator should time out via context cancellation.
	node := newCoordinatorTestNode() // sends nothing
	addr := startCoordinatorTestNode(t, node)

	coord := mpcinfra.NewMPCCoordinator(slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := coord.SignTransaction(ctx, apieth.MPCSigningRequest{
		SessionID: "timeout-session",
		Hash:      make([]byte, 32),
		PartyIDs:  []string{"party-1"},
		PeerAddrs: []string{addr},
		Threshold: 1,
	})

	require.Error(t, err)
}
