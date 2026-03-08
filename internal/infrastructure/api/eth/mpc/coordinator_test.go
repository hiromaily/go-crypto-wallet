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

// reactiveTestNode is a coordinatorTestNode whose RelaySession sends a configured
// reply as soon as it receives the first inbound message. This makes it possible
// to verify routing without relying on post-Close() delivery timing.
type reactiveTestNode struct {
	protogen.UnimplementedMPCNodeServiceServer

	mu       sync.Mutex
	received chan []byte
	reply    *mpcinfra.MPCWireMessage // sent after the first inbound message arrives
}

func newReactiveTestNode(reply mpcinfra.MPCWireMessage) *reactiveTestNode {
	return &reactiveTestNode{
		received: make(chan []byte, 64),
		reply:    &reply,
	}
}

func (*reactiveTestNode) InitSigning(
	_ context.Context,
	req *protogen.MPCSigningRequest,
) (*protogen.MPCSigningResponse, error) {
	return (&protogen.MPCSigningResponse_builder{SessionId: req.GetSessionId()}).Build(), nil
}

func (n *reactiveTestNode) RelaySession(
	stream grpc.BidiStreamingServer[protogen.MPCMessage, protogen.MPCMessage],
) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return nil
		}
		n.received <- msg.GetPayload()

		// On first message, send the pre-configured reply if set.
		n.mu.Lock()
		reply := n.reply
		n.reply = nil
		n.mu.Unlock()
		if reply != nil {
			data, err := json.Marshal(reply)
			if err != nil {
				return err
			}
			if err := stream.Send((&protogen.MPCMessage_builder{Payload: data}).Build()); err != nil {
				return err
			}
		}
	}
}

func TestMPCCoordinator_SignTransaction_RoutesBroadcast(t *testing.T) {
	t.Parallel()

	sessionID := "test-session-route"
	hash := make([]byte, 32)
	sig := make([]byte, 65)
	sig[64] = 28

	// Node 1: sends a broadcast round message (no To → broadcast to all except sender).
	roundMsg := []byte("tss-round-1-from-party-1")
	node1 := newCoordinatorTestNode(
		mpcinfra.MPCWireMessage{
			From: "party-1",
			Data: roundMsg,
		},
	)

	// Node 2: reactive — sends the signature after it receives node 1's routed broadcast.
	// This guarantees that the broadcast was delivered BEFORE coord.SignTransaction returns,
	// avoiding any post-Close() delivery race.
	node2 := newReactiveTestNode(mpcinfra.MPCWireMessage{
		From:        "party-2",
		IsSignature: true,
		Data:        sig,
	})

	addr1 := startCoordinatorTestNode(t, node1)

	// Start node2's gRPC server (reactive node implements the same interface).
	lis2, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv2 := grpc.NewServer()
	protogen.RegisterMPCNodeServiceServer(srv2, node2)
	go func() { _ = srv2.Serve(lis2) }()
	t.Cleanup(func() { srv2.GracefulStop() })
	addr2 := lis2.Addr().String()

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
	assert.Equal(t, sig, result.Signature)

	// node2 received the routed broadcast message (it must have, since it replied with the sig).
	select {
	case payload := <-node2.received:
		var relayed mpcinfra.MPCWireMessage
		require.NoError(t, json.Unmarshal(payload, &relayed), "routed message must be valid MPCWireMessage JSON")
		assert.Equal(t, "party-1", relayed.From)
		assert.Equal(t, roundMsg, relayed.Data)
	default:
		t.Fatal("node2 did not record any received message")
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
