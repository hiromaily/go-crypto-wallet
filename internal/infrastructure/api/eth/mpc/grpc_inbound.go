package mpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// Compile-time interface check.
var _ apieth.MPCInboundTransport = (*GRPCInboundTransport)(nil)

const (
	inboundRecvBuf = 256
	inboundSendBuf = 256
)

// GRPCInboundTransport implements MPCInboundTransport as a gRPC server.
//
// When the Watch wallet coordinator calls InitSigning, the expected session ID is registered.
// Subsequent RelaySession messages are validated against this session ID; messages with an
// unknown or mismatched session ID are silently discarded to prevent cross-session routing.
//
// The MPCNodeServer (Task 4.2) calls EnqueueOutbound to send messages back to the coordinator
// over the open bidirectional RelaySession stream.
type GRPCInboundTransport struct {
	protogen.UnimplementedMPCNodeServiceServer

	mu        sync.RWMutex
	sessionID string // set by InitSigning; empty means no session registered

	server  *grpc.Server
	started atomic.Bool

	recvCh chan []byte // coordinator → node messages delivered to MPCNodeServer
	sendCh chan []byte // node → coordinator messages queued by MPCNodeServer
}

// NewGRPCInboundTransport creates a transport ready for Listen to be called.
func NewGRPCInboundTransport() *GRPCInboundTransport {
	return &GRPCInboundTransport{
		recvCh: make(chan []byte, inboundRecvBuf),
		sendCh: make(chan []byte, inboundSendBuf),
	}
}

// Listen starts the gRPC server on listenAddr. It binds the TCP port and starts accepting
// connections in the background. Listen returns immediately after the port is bound.
// Returns an error if the transport is already started or the address cannot be bound.
func (t *GRPCInboundTransport) Listen(_ context.Context, listenAddr string) error {
	if !t.started.CompareAndSwap(false, true) {
		return errors.New("mpc inbound transport: already started")
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		t.started.Store(false)
		return fmt.Errorf("mpc inbound transport: listen on %s: %w", listenAddr, err)
	}

	t.server = grpc.NewServer()
	protogen.RegisterMPCNodeServiceServer(t.server, t)

	go func() { _ = t.server.Serve(lis) }()
	return nil
}

// Receive returns the channel on which coordinator-to-node messages arrive after session ID
// validation. The channel is buffered; the caller should drain it continuously.
func (t *GRPCInboundTransport) Receive(_ context.Context) (<-chan []byte, error) {
	return t.recvCh, nil
}

// Close gracefully stops the gRPC server, stopping acceptance of new RPCs.
// In-flight RPCs are allowed to complete. Safe to call multiple times.
func (t *GRPCInboundTransport) Close() error {
	if t.server != nil {
		t.server.GracefulStop()
	}
	return nil
}

// EnqueueOutbound queues msg to be sent to the coordinator over the active RelaySession stream.
// This method is NOT part of the MPCInboundTransport interface; it is called by MPCNodeServer
// (infrastructure, Task 4.2) to relay outbound TSS protocol messages back to the coordinator.
func (t *GRPCInboundTransport) EnqueueOutbound(msg []byte) error {
	select {
	case t.sendCh <- msg:
		return nil
	default:
		return errors.New("mpc inbound transport: outbound send buffer full")
	}
}

// =============================================================================
// MPCNodeServiceServer implementation (gRPC server handlers)
// =============================================================================

// InitSigning registers the session ID for this signing session.
// The coordinator calls this before opening the RelaySession stream.
func (t *GRPCInboundTransport) InitSigning(
	_ context.Context,
	req *protogen.MPCSigningRequest,
) (*protogen.MPCSigningResponse, error) {
	t.mu.Lock()
	t.sessionID = req.GetSessionId()
	t.mu.Unlock()

	return (&protogen.MPCSigningResponse_builder{SessionId: req.GetSessionId()}).Build(), nil
}

// RelaySession handles the bidirectional stream used for TSS round message relay.
//
// Inbound messages (coordinator → node) are validated against the registered session ID;
// messages with an unknown or mismatched ID are silently discarded. Valid messages are
// forwarded to recvCh for consumption by MPCNodeServer.
//
// Outbound messages queued via EnqueueOutbound are sent back to the coordinator
// over the same stream, allowing the node to push round messages without a separate connection.
func (t *GRPCInboundTransport) RelaySession(
	stream grpc.BidiStreamingServer[protogen.MPCMessage, protogen.MPCMessage],
) error {
	recvDone := make(chan struct{})

	// Goroutine: receive from coordinator, validate session ID, forward to recvCh.
	go func() {
		defer close(recvDone)
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			t.mu.RLock()
			expected := t.sessionID
			t.mu.RUnlock()

			if expected == "" || msg.GetSessionId() != expected {
				continue // discard messages with unknown/mismatched session ID
			}
			t.recvCh <- msg.GetPayload()
		}
	}()

	// Send loop: drain sendCh and write outbound messages back to the coordinator.
	for {
		select {
		case payload := <-t.sendCh:
			out := (&protogen.MPCMessage_builder{Payload: payload}).Build()
			if err := stream.Send(out); err != nil {
				<-recvDone
				return nil
			}
		case <-recvDone:
			return nil
		}
	}
}
