### Blockchain Communication

#### BTC/BCH (Bitcoin Core RPC API)

**Implementation**: Bitcoin Core RPC API

**Location**: `internal/infrastructure/api/btc/`

**Features**:

- Transaction creation and broadcasting
- UTXO management
- Address generation
- Block and transaction queries

**API Clients**:

- `btc/`: Bitcoin Core RPC client
- `bch/`: Bitcoin Cash Core RPC client

#### ETH (Ethereum JSON-RPC API)

**Implementation**: Ethereum JSON-RPC API

**Location**: `internal/infrastructure/api/eth/`

**Features**:

- Transaction creation and broadcasting
- Smart contract interaction
- ERC-20 token support
- Gas price estimation
- Nonce management

**API Clients**:

- `eth/`: Ethereum JSON-RPC client
- `erc20/`: ERC-20 token client

#### XRP

TODO: updating

**Implementation**: Communication via

**Location**: `internal/infrastructure/api/xrp/`

**Features**:

- Transaction creation and submission
- Account management
- Balance queries
- Payment tracking

**API Client**:

- `xrp/`: Ripple gRPC client [Deprecated]

**Note**: XRP uses protocol buffers for gRPC communication. See [Code Generation Guidelines](../../../../docs/guidelines/code-generation.md) for protobuf code generation.
