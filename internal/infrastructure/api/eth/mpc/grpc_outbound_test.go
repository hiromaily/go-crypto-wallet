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

// testRelayServer is a minimal in-process gRPC server for testing GRPCOutboundTransport.
type testRelayServer struct {
	protogen.UnimplementedMPCNodeServiceServer
	received chan []byte
	toSend   chan []byte
}

func newTestRelayServer() *testRelayServer {
	return &testRelayServer{
		received: make(chan []byte, 32),
		toSend:   make(chan []byte, 32),
	}
}

func (*testRelayServer) InitSigning(
	_ context.Context, req *protogen.MPCSigningRequest,
) (*protogen.MPCSigningResponse, error) {
	return (&protogen.MPCSigningResponse_builder{SessionId: req.GetSessionId()}).Build(), nil
}

func (s *testRelayServer) RelaySession(
	stream grpc.BidiStreamingServer[protogen.MPCMessage, protogen.MPCMessage],
) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			s.received <- msg.GetPayload()
			// Echo any queued server-side message back to the client.
			select {
			case payload := <-s.toSend:
				_ = stream.Send((&protogen.MPCMessage_builder{Payload: payload}).Build())
			default:
			}
		}
	}()
	<-done
	return nil
}

// startTestServer starts an in-process gRPC server and returns its address.
func startTestServer(t *testing.T, srv *testRelayServer) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	s := grpc.NewServer()
	protogen.RegisterMPCNodeServiceServer(s, srv)
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(func() { s.GracefulStop() })
	return lis.Addr().String()
}

func TestGRPCOutboundTransport_Send(t *testing.T) {
	t.Parallel()

	srv := newTestRelayServer()
	addr := startTestServer(t, srv)

	transport := mpcinfra.NewGRPCOutboundTransport()
	t.Cleanup(func() { _ = transport.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := []byte("test-tss-round-msg")
	err := transport.Send(ctx, addr, payload)
	require.NoError(t, err)

	select {
	case got := <-srv.received:
		assert.Equal(t, payload, got)
	case <-time.After(3 * time.Second):
		t.Fatal("server did not receive the message within timeout")
	}
}

func TestGRPCOutboundTransport_Receive(t *testing.T) {
	t.Parallel()

	srv := newTestRelayServer()
	addr := startTestServer(t, srv)

	transport := mpcinfra.NewGRPCOutboundTransport()
	t.Cleanup(func() { _ = transport.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Queue a response the server will echo back on the first Send.
	responsePayload := []byte("node-response-msg")
	srv.toSend <- responsePayload

	err := transport.Send(ctx, addr, []byte("trigger"))
	require.NoError(t, err)

	recvCh, err := transport.Receive(ctx)
	require.NoError(t, err)

	select {
	case got := <-recvCh:
		assert.Equal(t, responsePayload, got)
	case <-time.After(3 * time.Second):
		t.Fatal("Receive did not deliver server response within timeout")
	}
}

func TestGRPCOutboundTransport_Send_UnreachablePeer(t *testing.T) {
	t.Parallel()

	transport := mpcinfra.NewGRPCOutboundTransport()
	t.Cleanup(func() { _ = transport.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := transport.Send(ctx, "127.0.0.1:1", []byte("msg"))
	assert.Error(t, err)
}

func TestGRPCOutboundTransport_Close(t *testing.T) {
	t.Parallel()

	srv := newTestRelayServer()
	addr := startTestServer(t, srv)

	transport := mpcinfra.NewGRPCOutboundTransport()

	ctx := context.Background()
	err := transport.Send(ctx, addr, []byte("ping"))
	require.NoError(t, err)

	err = transport.Close()
	assert.NoError(t, err)
}

func TestGRPCOutboundTransport_ReusesConnection(t *testing.T) {
	t.Parallel()

	srv := newTestRelayServer()
	addr := startTestServer(t, srv)

	transport := mpcinfra.NewGRPCOutboundTransport()
	t.Cleanup(func() { _ = transport.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messages := [][]byte{[]byte("msg-1"), []byte("msg-2"), []byte("msg-3")}
	for _, m := range messages {
		require.NoError(t, transport.Send(ctx, addr, m))
	}

	for i := range messages {
		select {
		case got := <-srv.received:
			assert.Equal(t, messages[i], got)
		case <-time.After(3 * time.Second):
			t.Fatalf("message %d not received within timeout", i)
		}
	}
}

func TestGRPCOutboundTransport_InsecureDial(t *testing.T) {
	t.Parallel()

	srv := newTestRelayServer()
	addr := startTestServer(t, srv)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	client := protogen.NewMPCNodeServiceClient(conn)
	req := (&protogen.MPCSigningRequest_builder{SessionId: "test-session"}).Build()
	resp, err := client.InitSigning(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "test-session", resp.GetSessionId())
}
