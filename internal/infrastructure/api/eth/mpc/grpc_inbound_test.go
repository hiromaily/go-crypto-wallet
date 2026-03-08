package mpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mpcinfra "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc/protogen"
)

// freeAddr returns an available TCP address by binding then closing immediately.
func freeAddr(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := lis.Addr().String()
	require.NoError(t, lis.Close())
	return addr
}

// newTestGRPCClient creates an insecure gRPC client to the given server address.
func newTestGRPCClient(t *testing.T, addr string) protogen.MPCNodeServiceClient {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return protogen.NewMPCNodeServiceClient(conn)
}

func TestGRPCInboundTransport_Listen_Receive(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	transport := mpcinfra.NewGRPCInboundTransport()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := transport.Listen(ctx, addr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = transport.Close() })

	time.Sleep(50 * time.Millisecond)

	client := newTestGRPCClient(t, addr)

	sessionID := "session-abc"
	req := (&protogen.MPCSigningRequest_builder{SessionId: sessionID}).Build()
	resp, err := client.InitSigning(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, sessionID, resp.GetSessionId())

	stream, err := client.RelaySession(ctx)
	require.NoError(t, err)

	payload := []byte("tss-round-payload")
	msg := (&protogen.MPCMessage_builder{SessionId: sessionID, Payload: payload}).Build()
	err = stream.Send(msg)
	require.NoError(t, err)

	recvCh, err := transport.Receive(ctx)
	require.NoError(t, err)

	select {
	case got := <-recvCh:
		assert.Equal(t, payload, got)
	case <-time.After(3 * time.Second):
		t.Fatal("Receive did not deliver the message within timeout")
	}
}

func TestGRPCInboundTransport_ValidatesSessionID(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	transport := mpcinfra.NewGRPCInboundTransport()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, transport.Listen(ctx, addr))
	t.Cleanup(func() { _ = transport.Close() })

	time.Sleep(50 * time.Millisecond)

	client := newTestGRPCClient(t, addr)

	req := (&protogen.MPCSigningRequest_builder{SessionId: "good-session"}).Build()
	_, err := client.InitSigning(ctx, req)
	require.NoError(t, err)

	stream, err := client.RelaySession(ctx)
	require.NoError(t, err)

	// Send a message with a wrong session ID — must be discarded.
	badMsg := (&protogen.MPCMessage_builder{SessionId: "bad-session", Payload: []byte("should-be-discarded")}).Build()
	err = stream.Send(badMsg)
	require.NoError(t, err)

	// Send a message with the correct session ID — must be delivered.
	validPayload := []byte("valid-payload")
	goodMsg := (&protogen.MPCMessage_builder{SessionId: "good-session", Payload: validPayload}).Build()
	err = stream.Send(goodMsg)
	require.NoError(t, err)

	recvCh, err := transport.Receive(ctx)
	require.NoError(t, err)

	select {
	case got := <-recvCh:
		assert.Equal(t, validPayload, got, "only the valid-session message should be delivered")
	case <-time.After(3 * time.Second):
		t.Fatal("Receive did not deliver the valid message within timeout")
	}

	// No extra messages should arrive.
	select {
	case extra := <-recvCh:
		t.Fatalf("unexpected extra message received: %s", string(extra))
	case <-time.After(200 * time.Millisecond):
		// Good.
	}
}

func TestGRPCInboundTransport_EnqueueOutbound(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	transport := mpcinfra.NewGRPCInboundTransport()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, transport.Listen(ctx, addr))
	t.Cleanup(func() { _ = transport.Close() })

	time.Sleep(50 * time.Millisecond)

	client := newTestGRPCClient(t, addr)

	sessionID := "session-out"
	req := (&protogen.MPCSigningRequest_builder{SessionId: sessionID}).Build()
	_, err := client.InitSigning(ctx, req)
	require.NoError(t, err)

	stream, err := client.RelaySession(ctx)
	require.NoError(t, err)

	// Send a trigger message to open the server-side stream.
	triggerMsg := (&protogen.MPCMessage_builder{SessionId: sessionID, Payload: []byte("trigger")}).Build()
	err = stream.Send(triggerMsg)
	require.NoError(t, err)

	// Give the server time to set up the stream handler.
	time.Sleep(100 * time.Millisecond)

	// Enqueue an outbound message from the node back to the coordinator.
	outPayload := []byte("node-outbound-response")
	err = transport.EnqueueOutbound(outPayload)
	require.NoError(t, err)

	// The client (coordinator) should receive it.
	reply, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, outPayload, reply.GetPayload())
}

func TestGRPCInboundTransport_Close(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	transport := mpcinfra.NewGRPCInboundTransport()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, transport.Listen(ctx, addr))

	err := transport.Close()
	assert.NoError(t, err)
}

func TestGRPCInboundTransport_Listen_AlreadyStarted(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	transport := mpcinfra.NewGRPCInboundTransport()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, transport.Listen(ctx, addr))
	t.Cleanup(func() { _ = transport.Close() })

	time.Sleep(50 * time.Millisecond)

	// Second Listen call must return an error.
	err := transport.Listen(ctx, addr)
	assert.Error(t, err)
}
