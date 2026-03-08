package mpc

// MPCWireMessage is the internal JSON envelope used for coordinator ↔ node message relay.
//
// Nodes encode their TSS round messages in this envelope so the coordinator can route
// them to the appropriate peer(s). The coordinator inspects From/To/IsSignature,
// then forwards Data (raw TSS bytes) to the target peer(s).
//
// Layout:
//   - node → coordinator: wrapped in MPCWireMessage so the coordinator knows how to route.
//   - coordinator → node: sends Data bytes directly (no wrapper); nodes receive raw TSS bytes.
type MPCWireMessage struct {
	// From is the party ID of the node that generated this message.
	// The coordinator uses this to implement broadcast-except-sender semantics.
	From string `json:"from"`

	// To contains the target party IDs.
	// An empty (or nil) slice means broadcast to all nodes except From.
	// A non-empty slice means unicast/multicast to exactly those parties.
	To []string `json:"to,omitempty"`

	// IsSignature indicates that Data contains the assembled 65-byte ECDSA signature
	// and the signing session is complete. The coordinator stops relaying and returns
	// the signature to the caller when this flag is true.
	IsSignature bool `json:"is_sig,omitempty"`

	// Data is the TSS round message bytes, or the 65-byte r||s||v signature when
	// IsSignature is true.
	Data []byte `json:"data"`
}
