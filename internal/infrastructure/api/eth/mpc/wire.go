package mpc

// MPCWireMessage is the internal JSON envelope used for coordinator ↔ node message relay.
//
// Both directions of the relay use this struct so nodes always receive consistent metadata:
//
//   - node → coordinator: {From, To, IsBroadcast, Data} where Data is the tss.Message wire bytes.
//     The special IsSignature flag signals that Data is the assembled 65-byte ECDSA signature
//     and the signing session is complete.
//
//   - coordinator → node: {From, IsBroadcast, Data}; To is omitted because the recipient
//     does not need routing information — it only needs to know the sender (From) and whether
//     the message was broadcast (IsBroadcast) to correctly call tss.ParseWireMessage.
type MPCWireMessage struct {
	// From is the party ID of the node that generated this message.
	// The coordinator uses this to implement broadcast-except-sender semantics.
	// Receiving nodes pass it to tss.ParseWireMessage to identify the sender.
	From string `json:"from"`

	// To contains the target party IDs.
	// An empty (or nil) slice means broadcast to all nodes except From.
	// A non-empty slice means unicast/multicast to exactly those parties.
	// Omitted by the coordinator when forwarding to a recipient (recipient does not need it).
	To []string `json:"to,omitempty"`

	// IsBroadcast is true when the TSS round message was broadcast to all parties.
	// Nodes must pass this flag verbatim to tss.ParseWireMessage so that the library
	// can correctly route the message inside the local state machine.
	IsBroadcast bool `json:"is_broadcast,omitempty"`

	// IsSignature indicates that Data contains the assembled 65-byte ECDSA signature
	// and the signing session is complete. The coordinator stops relaying and returns
	// the signature to the caller when this flag is true.
	IsSignature bool `json:"is_sig,omitempty"`

	// Data is the TSS round message bytes from tss.Message.WireBytes(), or the
	// 65-byte r||s||v signature when IsSignature is true.
	Data []byte `json:"data"`
}
