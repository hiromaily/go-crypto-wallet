# CLI List

Complete command reference for the three wallet CLIs.
Command hierarchy follows the rule: **[wallet] → [verb] → [target]**.

> **For detailed usage examples, flags, and workflow guidance, see [`docs/commands.md`](../../../docs/commands.md).**

---

## Global Flags

All wallet CLIs share the following flags:

| Flag               | Short | Required | Description                                          |
|--------------------|-------|----------|------------------------------------------------------|
| `--config`         | `-c`  | Yes      | Path to the wallet configuration file                |
| `--coin`           |       | No       | Coin type: `btc`, `bch`, `eth`, `xrp` (default: `btc`) |
| `--account-config` |       | No       | Path to the account config file (multisig)           |
| `--wallet`         | `-w`  | No       | Bitcoin Core wallet name (BTC/BCH only)              |

---

## watch

Online wallet. Manages public keys, creates unsigned transactions, sends signed transactions, and monitors the blockchain.

```
watch
├── import
│   ├── address                    Import addresses exported by keygen wallet
│   └── descriptor                 (BTC only) Import descriptors into Bitcoin Core
│                                  Use --validate flag for dry-run validation
├── create
│   ├── deposit                    Create unsigned deposit transaction
│   ├── payment                    Create unsigned payment transaction
│   ├── transfer                   Create unsigned transfer transaction between accounts
│   └── multisig                   (ETH only) Create unsigned Safe multisig transaction proposal file
│                                  Flags: --safe <addr>, --to <addr>, --amount <ETH>,
│                                         --threshold <n>, --action-type <type>
├── send
│   ├── tx                         Send signed transaction to the network
│   └── multisig                   Multi-signature operations
│       ├── send-eth               (ETH) Submit a fully signed Safe multisig transaction
│       │                          Flags: --file <path>
│       ├── collect-nonces         (BTC) Collect and aggregate MuSig2 nonces from all signers
│       ├── aggregate              (BTC) Aggregate MuSig2 partial signatures into final tx
│       ├── set-regular-key        (XRP) Set or remove a regular key for an account
│       ├── set-signer-list        (XRP) Configure the multisig signer list
│       ├── create-multisig-tx     (XRP) Create a pending multisig transaction
│       ├── add-multisig-signature (XRP) Add a signer's signature to a pending tx
│       └── submit-multisig-tx     (XRP) Submit a quorum-complete multisig tx
├── safe                           (ETH only) Safe contract operations
│   └── info                       Retrieve on-chain Safe contract state (owners, threshold, nonce)
│                                  Flags: --safe <addr>
├── monitor
│   ├── senttx                     Monitor sent transaction confirmations
│   └── balance                    Monitor account balances
└── api                            Chain node RPC commands (coin-specific)
    ├── btc / bch
    │   ├── balance                Get balance for account
    │   ├── estimatefee            Estimate transaction fee
    │   ├── getnetworkinfo         Get network info from node
    │   ├── getaddressinfo         Get info for a specific address
    │   ├── listunspent            List unspent transaction outputs
    │   ├── logging                Set Bitcoin Core log level
    │   ├── unlocktx               Unlock locked UTXOs
    │   └── validateaddress        Validate an address
    ├── eth
    │   ├── clientversion          Get Ethereum client version
    │   ├── nodeinfo               Get node information
    │   ├── syncing                Get sync status
    │   └── netversion             Get network version
    └── xrp
        └── sendcoin               Send XRP from faucet (testing only)
```

---

## keygen

Offline wallet. Generates keys, creates descriptors and multisig addresses, signs MuSig2 rounds, and manages wallet encryption.

```
keygen
├── create
│   ├── key                        Create one key (debug use)
│   ├── hdkey                      Create HD wallet keys
│   ├── seed                       Create or store a seed
│   ├── descriptor                 (BTC only) Descriptor operations
│   │   ├── generate               Generate a descriptor for an account
│   │   ├── export                 Export descriptors for one account to file
│   │   └── export-all             Export all descriptors (receive + change) to file
│   └── multisig                   (BTC only) Create multisig address (traditional P2SH/P2WSH or MuSig2 Taproot)
├── export
│   ├── address                    Export generated public keys as CSV
│   ├── descriptor                 (BTC only) Export output descriptors to file
│   └── fullpubkey                 (ETH only) Export account xpub for Watch wallet address derivation
├── import
│   ├── privkey                    (BTC/BCH/ETH) Import private keys into node wallet / keystore
│   └── fullpubkey                 (BTC only) Import full public key exported by sign wallet
├── sign
│   ├── signature                  Sign an unsigned transaction (1st signature)
│   │                              ETH Safe multisig: use --file <multisig-json> --signer-address <addr>
│   └── musig2                     (BTC only) MuSig2 signing operations
│       ├── nonce                  Generate MuSig2 nonce (Round 1)
│       └── sign                   Create MuSig2 partial signature (Round 2)
└── api                            Chain node RPC commands (BTC/BCH only)
    ├── encryptwallet              Encrypt the Bitcoin Core wallet with a passphrase
    ├── walletpassphrase           Unlock the wallet for N seconds
    ├── walletpassphrasechange     Change the wallet passphrase
    ├── walletlock                 Lock the wallet immediately
    ├── dumpwallet                 Dump all wallet keys to a file
    └── importwallet               Import keys from a wallet dump file
```

---

## sign

Offline wallet (authorization signer). Generates authorization keys, signs transactions, and performs MuSig2 signing rounds.

```
sign
├── create
│   ├── key                        Create one key (debug use)
│   ├── hdkey                      Create HD wallet authorization keys
│   └── seed                       Create or store a seed
├── export
│   └── fullpubkey                 (BTC only) Export full public key for import by keygen wallet
├── import
│   └── privkey                    (BTC only) Import authorization private key into Bitcoin Core
├── sign
│   ├── signature                  Sign a multisig transaction (2nd+ signature)
│   │                              ETH Safe multisig: use --file <multisig-json> --signer-address <addr>
│   └── musig2                     (BTC only) MuSig2 signing operations
│       ├── nonce                  Generate MuSig2 nonce (Round 1)
│       └── sign                   Create MuSig2 partial signature (Round 2)
└── api                            Chain node RPC commands (BTC/BCH only)
    ├── encryptwallet              Encrypt the Bitcoin Core wallet with a passphrase
    ├── walletpassphrase           Unlock the wallet for N seconds
    ├── walletpassphrasechange     Change the wallet passphrase
    └── walletlock                 Lock the wallet immediately
```

---

## Command x Chain x UseCase Matrix

This section is the **single source of truth (SSOT)** for:

- Which commands are available per chain (`✓` supported / `–` not applicable)
- The use case interface each command invokes
- The role of each command in the wallet workflow

Interface references are in `internal/application/usecase/{wallet}/interfaces.go`.
`(node RPC)` = direct call to the chain node API, bypasses the use case layer.

### watch

| Command | BTC | BCH | ETH | XRP | Use Case Interface | Role |
|---------|:---:|:---:|:---:|:---:|---|---|
| `import address` | ✓ | ✓ | ✓ | ✓ | `watch.ImportAddressUseCase` | Import addresses exported by keygen wallet into watch DB |
| `import descriptor` | ✓ | – | – | – | `watch.ImportDescriptorUseCase` | Import output descriptors into Bitcoin Core; enables address derivation without importing individual keys |
| `create deposit` | ✓ | ✓ | ✓ | ✓ | `watch.CreateTransactionUseCase` | Create unsigned deposit tx: aggregate coins from client addresses into cold wallet |
| `create payment` | ✓ | ✓ | ✓ | ✓ | `watch.CreateTransactionUseCase` | Create unsigned payment tx: send coins from cold wallet to user-specified addresses |
| `create transfer` | ✓ | ✓ | ✓ | ✓ | `watch.CreateTransactionUseCase` | Create unsigned transfer tx: move coins between internal accounts (e.g. deposit → payment) |
| `create multisig` | – | – | ✓ | – | `watch.CreateETHMultisigTransactionUseCase` | (ETH) Propose a new Safe multisig tx: fetch on-chain nonce, compute EIP-712 hash, write unsigned JSON proposal file |
| `send tx` | ✓ | ✓ | ✓ | ✓ | `watch.SendTransactionUseCase` | Broadcast a signed transaction file to the blockchain network |
| `send multisig send-eth` | – | – | ✓ | – | `watch.SendETHMultisigTransactionUseCase` | (ETH) Submit a fully signed Safe multisig JSON file by calling `execTransaction` on-chain |
| `send multisig collect-nonces` | ✓ | – | – | – | `watch.AggregateMuSig2SignaturesUseCase` | Collect MuSig2 nonces from all signers and embed them in PSBT (Round 1 aggregation) |
| `send multisig aggregate` | ✓ | – | – | – | `watch.AggregateMuSig2SignaturesUseCase` | Aggregate MuSig2 partial signatures from all signers into the final transaction (Round 2) |
| `send multisig set-regular-key` | – | – | – | ✓ | `watch.SetRegularKeyUseCase` | Set or remove a regular key on an XRP account (allows signing without the master key) |
| `send multisig set-signer-list` | – | – | – | ✓ | `watch.SetSignerListUseCase` | Configure the XRP multi-sig signer list and quorum for an account |
| `send multisig create-multisig-tx` | – | – | – | ✓ | `watch.CreateMultisigTxUseCase` | Create a pending XRP multi-sig transaction stored in DB awaiting signatures |
| `send multisig add-multisig-signature` | – | – | – | ✓ | `watch.AddMultisigSignatureUseCase` | Record a signer's signed-tx blob; auto-combines when quorum is met |
| `send multisig submit-multisig-tx` | – | – | – | ✓ | `watch.SubmitMultisigTxUseCase` | Submit a quorum-complete combined XRP multi-sig transaction to the ledger |
| `monitor senttx` | ✓ | ✓ | ✓ | ✓ | `watch.MonitorTransactionUseCase` | Update confirmation status of sent transactions in watch DB |
| `monitor balance` | ✓ | ✓ | ✓ | ✓ | `watch.MonitorTransactionUseCase` | Query and record current account balances from the chain |
| `api balance` | ✓ | ✓ | – | – | _(node RPC)_ | Query account balance directly from Bitcoin Core node |
| `api estimatefee` | ✓ | ✓ | – | – | _(node RPC)_ | Estimate smart fee (sat/vB) from Bitcoin Core node |
| `api getnetworkinfo` | ✓ | ✓ | – | – | _(node RPC)_ | Retrieve version and network info from Bitcoin Core node |
| `api getaddressinfo` | ✓ | ✓ | – | – | _(node RPC)_ | Query full metadata for a specific address from Bitcoin Core |
| `api listunspent` | ✓ | ✓ | – | – | _(node RPC)_ | List UTXOs for an account with optional confirmation filter |
| `api logging` | ✓ | ✓ | – | – | _(node RPC)_ | Adjust Bitcoin Core log categories/level at runtime |
| `api unlocktx` | ✓ | ✓ | – | – | _(node RPC)_ | Unlock locked UTXOs so they can be spent again |
| `api validateaddress` | ✓ | ✓ | – | – | _(node RPC)_ | Validate whether a Bitcoin/BCH address is well-formed |
| `safe info` | – | – | ✓ | – | `watch.ETHSafeInfoUseCase` | (ETH) Query on-chain Safe state: owners, threshold, nonce |
| `api clientversion` | – | – | ✓ | – | _(node RPC)_ | Get Ethereum client (geth/etc.) version string |
| `api nodeinfo` | – | – | ✓ | – | _(node RPC)_ | Get Ethereum node peer and connection info |
| `api syncing` | – | – | ✓ | – | _(node RPC)_ | Get Ethereum node sync status and progress |
| `api netversion` | – | – | ✓ | – | _(node RPC)_ | Get Ethereum network/chain ID |
| `api sendcoin` | – | – | – | ✓ | _(node RPC)_ | Send XRP from faucet account to a target address (testing only) |

### keygen

| Command | BTC | BCH | ETH | XRP | Use Case Interface | Role |
|---------|:---:|:---:|:---:|:---:|---|---|
| `create seed` | ✓ | ✓ | ✓ | ✓ | `keygen.GenerateSeedUseCase` | Generate or store a BIP39 seed in the cold wallet DB |
| `create key` | ✓ | ✓ | ✓ | – | _(debug)_ | Create one key pair for debugging and verification |
| `create hdkey` | ✓ | ✓ | ✓ | ✓ | `keygen.GenerateHDWalletUseCase` / `keygen.GenerateKeyUseCase` | Generate HD keys (BIP32/44); use `--keypair` flag for XRP key pair generation |
| `create descriptor generate` | ✓ | – | – | – | `keygen.GenerateDescriptorUseCase` | Generate an output descriptor for a given account and address type |
| `create descriptor export` | ✓ | – | – | – | `keygen.ExportDescriptorUseCase` | Export descriptors for one account to file (bitcoin-core/json/text format) |
| `create descriptor export-all` | ✓ | – | – | – | `keygen.ExportDescriptorUseCase` | Export all receive + change descriptors for an account in one call |
| `create multisig` | ✓ | – | – | – | `keygen.CreateMultisigAddressUseCase` / `keygen.CreateMuSig2AddressUseCase` | Create P2SH/P2WSH (traditional) or MuSig2 Taproot multisig address |
| `export address` | ✓ | ✓ | ✓ | ✓ | `keygen.ExportAddressUseCase` | Export generated public key addresses as CSV for Watch wallet import |
| `export descriptor` | ✓ | – | – | – | `keygen.ExportDescriptorUseCase` | Export output descriptors to file (shorthand path under `export`) |
| `export fullpubkey` | – | – | ✓ | – | `keygen.ExportFullPubkeyUseCase` | Export account-level xpub (ETH only) so Watch wallet can derive and verify child addresses |
| `import privkey` | ✓ | ✓ | ✓ | – | `keygen.ImportPrivateKeyUseCase` | Import private keys into Bitcoin Core node wallet (BTC/BCH) or ETH keystore |
| `import fullpubkey` | ✓ | – | – | – | `keygen.ImportFullPubkeyUseCase` | Import sign wallet's public key to enable multisig address creation |
| `sign signature` | ✓ | ✓ | ✓ | ✓ | `keygen.SignTransactionUseCase` (BTC/BCH/XRP) `keygen.SignMultisigTransactionUseCase` (ETH Safe) | Sign unsigned transaction; provides the 1st signature in multisig flows. ETH: detects multisig JSON format and routes to EIP-712 offline signing; use `--signer-address` to select key |
| `sign musig2 nonce` | ✓ | – | – | – | `keygen.GenerateMuSig2NonceUseCase` | Generate MuSig2 nonce (Round 1); share output with Watch wallet for aggregation |
| `sign musig2 sign` | ✓ | – | – | – | `keygen.MuSig2SignUseCase` | Create MuSig2 partial signature (Round 2); send signed PSBT to Watch for aggregation |
| `api encryptwallet` | ✓ | ✓ | – | – | _(node RPC)_ | Encrypt the Bitcoin Core wallet with a passphrase |
| `api walletpassphrase` | ✓ | ✓ | – | – | _(node RPC)_ | Temporarily unlock the encrypted wallet for N seconds |
| `api walletpassphrasechange` | ✓ | ✓ | – | – | _(node RPC)_ | Change the wallet encryption passphrase |
| `api walletlock` | ✓ | ✓ | – | – | _(node RPC)_ | Lock the wallet immediately (remove decryption key from memory) |
| `api dumpwallet` | ✓ | ✓ | – | – | _(node RPC)_ | Dump all wallet keys to a human-readable backup file |
| `api importwallet` | ✓ | ✓ | – | – | _(node RPC)_ | Import keys from a wallet dump file |

### sign

| Command | BTC | BCH | ETH | XRP | Use Case Interface | Role |
|---------|:---:|:---:|:---:|:---:|---|---|
| `create seed` | ✓ | ✓ | ✓ | ✓ | `sign.GenerateSeedUseCase` / `sign.StoreSeedUseCase` | Generate or store the authorization seed in the cold wallet DB |
| `create key` | ✓ | ✓ | ✓ | – | _(debug)_ | Create one authorization key pair for debugging |
| `create hdkey` | ✓ | ✓ | ✓ | ✓ | `sign.GenerateAuthKeyUseCase` | Generate HD authorization keys from the stored seed |
| `export fullpubkey` | ✓ | – | – | – | `sign.ExportFullPubkeyUseCase` | Export this signer's public key as CSV; keygen wallet imports it to build multisig addresses |
| `import privkey` | ✓ | – | – | – | `sign.ImportPrivateKeyUseCase` | Import authorization private keys into Bitcoin Core wallet |
| `sign signature` | ✓ | ✓ | ✓ | ✓ | `sign.SignTransactionUseCase` (BTC/BCH/XRP) `keygen.SignMultisigTransactionUseCase` (ETH Safe) | Sign a transaction that already has the 1st signature (2nd+ in multisig). ETH: detects multisig JSON format and routes to EIP-712 offline signing; use `--signer-address` to select key |
| `sign musig2 nonce` | ✓ | – | – | – | `sign.GenerateMuSig2NonceUseCase` | Generate MuSig2 nonce (Round 1); share with Watch wallet for nonce aggregation |
| `sign musig2 sign` | ✓ | – | – | – | `sign.MuSig2SignUseCase` | Create MuSig2 partial signature (Round 2); send signed PSBT to Watch for final aggregation |
| `api encryptwallet` | ✓ | ✓ | – | – | _(node RPC)_ | Encrypt the Bitcoin Core wallet with a passphrase |
| `api walletpassphrase` | ✓ | ✓ | – | – | _(node RPC)_ | Temporarily unlock the encrypted wallet for N seconds |
| `api walletpassphrasechange` | ✓ | ✓ | – | – | _(node RPC)_ | Change the wallet encryption passphrase |
| `api walletlock` | ✓ | ✓ | – | – | _(node RPC)_ | Lock the wallet immediately (remove decryption key from memory) |

---

## Related Documents

| Document | Description |
|----------|-------------|
| [`docs/commands.md`](../../../docs/commands.md) | Detailed flags, options, and usage examples per command |
| [`.claude/skills/wallet-cli/SKILL.md`](../../../.claude/skills/wallet-cli/SKILL.md) | How to run wallet CLI commands |
| [`.claude/rules/internal/cli-structure.md`](../../../.claude/rules/internal/cli-structure.md) | CLI command hierarchy rules |
| [`docs/chains/btc/operations/wallet-flow.md`](../../../docs/chains/btc/operations/wallet-flow.md) | BTC end-to-end transaction workflow |
