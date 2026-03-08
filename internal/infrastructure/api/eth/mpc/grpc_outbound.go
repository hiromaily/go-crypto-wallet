// Package mpc provides gRPC-based transport implementations for the MPC/TSS signing relay.
package mpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// Compile-time interface check.
var _ apieth.MPCOutboundTransport = (*GRPCOutboundTransport)(nil)

const outboundRecvBuf = 256

// peerConn holds an open gRPC connection and RelaySession bidirectional stream for one peer.
type peerConn struct {
	conn   *grpc.ClientConn
	stream grpc.BidiStreamingClient[protogen.MPCMessage, protogen.MPCMessage]
}

// GRPCOutboundTransport implements MPCOutboundTransport using gRPC bidirectional streaming.
//
// One RelaySession stream is opened per peer on the first Send call and reused for subsequent
// messages. Inbound messages from all peer streams are aggregated into the channel returned by
// Receive. Call Close when the signing session completes.
type GRPCOutboundTransport struct {
	mu     sync.Mutex
	peers  map[string]*peerConn
	recvCh chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewGRPCOutboundTransport creates a new transport ready to relay TSS messages.
// Call Close when the session is complete to release all gRPC connections.
func NewGRPCOutboundTransport() *GRPCOutboundTransport {
	ctx, cancel := context.WithCancel(context.Background())
	return &GRPCOutboundTransport{
		peers:  make(map[string]*peerConn),
		recvCh: make(chan []byte, outboundRecvBuf),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Send transmits msg to the MPC node at peerAddr.
//
// If no stream exists for peerAddr, a new gRPC connection and RelaySession stream are
// established with exponential backoff. The connection is reused on subsequent calls to the
// same peer. ctx is accepted for interface compatibility but is not used for stream lifetime
// (the stream lives until Close is called).
func (t *GRPCOutboundTransport) Send(ctx context.Context, peerAddr string, msg []byte) error {
	pc, err := t.getOrCreatePeer(ctx, peerAddr)
	if err != nil {
		return fmt.Errorf("mpc outbound transport: connect to %s: %w", peerAddr, err)
	}
	out := (&protogen.MPCMessage_builder{Payload: msg}).Build()
	if err := pc.stream.Send(out); err != nil {
		return fmt.Errorf("mpc outbound transport: send to %s: %w", peerAddr, err)
	}
	return nil
}

// Receive returns a channel that delivers inbound messages sent back by all connected peers.
// Messages arrive as raw payload bytes. The channel is buffered; drain it continuously.
func (t *GRPCOutboundTransport) Receive(_ context.Context) (<-chan []byte, error) {
	return t.recvCh, nil
}

// Close cancels the internal context, waits for all receive goroutines to exit, then closes
// every open gRPC connection. It is safe to call Close multiple times.
func (t *GRPCOutboundTransport) Close() error {
	t.cancel()
	t.wg.Wait()

	t.mu.Lock()
	defer t.mu.Unlock()

	var lastErr error
	for addr, pc := range t.peers {
		if err := pc.conn.Close(); err != nil {
			lastErr = fmt.Errorf("mpc outbound transport: close connection to %s: %w", addr, err)
		}
	}
	t.peers = make(map[string]*peerConn)
	return lastErr
}

// getOrCreatePeer returns an existing peerConn for peerAddr or establishes a new one.
// A background goroutine is started to read from the stream and forward messages to recvCh.
func (t *GRPCOutboundTransport) getOrCreatePeer(_ context.Context, peerAddr string) (*peerConn, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if pc, ok := t.peers[peerAddr]; ok {
		return pc, nil
	}

	bc := backoff.Config{
		BaseDelay:  100 * time.Millisecond,
		Multiplier: 1.6,
		Jitter:     0.2,
		MaxDelay:   10 * time.Second,
	}
	conn, err := grpc.NewClient(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: bc}),
	)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", peerAddr, err)
	}

	client := protogen.NewMPCNodeServiceClient(conn)

	// The stream's lifetime is bound to the transport (t.ctx).
	// Connection to unreachable peers fails fast via "connection refused" from the OS.
	stream, err := client.RelaySession(t.ctx)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open relay session to %s: %w", peerAddr, err)
	}

	pc := &peerConn{conn: conn, stream: stream}
	t.peers[peerAddr] = pc

	// Background goroutine: read responses from this peer and fan into recvCh.
	t.wg.Go(func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			select {
			case t.recvCh <- msg.GetPayload():
			case <-t.ctx.Done():
				return
			}
		}
	})

	return pc, nil
}
