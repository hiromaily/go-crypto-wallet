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
│   └── descriptor                 (BTC only) Descriptor operations
│       ├── import                 Import descriptors from file into Bitcoin Core
│       └── validate               Validate descriptors without importing
├── create
│   ├── deposit                    Create unsigned deposit transaction
│   ├── payment                    Create unsigned payment transaction
│   ├── transfer                   Create unsigned transfer transaction between accounts
│   └── db                         Create dummy payment_request data (dev only)
├── send
│   ├── tx                         Send signed transaction to the network
│   └── multisig                   Multi-signature operations
│       ├── collect-nonces         (BTC) Collect and aggregate MuSig2 nonces from all signers
│       ├── aggregate              (BTC) Aggregate MuSig2 partial signatures into final tx
│       ├── set-regular-key        (XRP) Set or remove a regular key for an account
│       ├── set-signer-list        (XRP) Configure the multisig signer list
│       ├── create-multisig-tx     (XRP) Create a pending multisig transaction
│       ├── add-multisig-signature (XRP) Add a signer's signature to a pending tx
│       └── submit-multisig-tx     (XRP) Submit a quorum-complete multisig tx
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
        └── sendcoin               Send XRP from faucet
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
│   └── multisig                   Create multisig address (traditional P2SH/P2WSH or MuSig2 Taproot)
├── export
│   ├── address                    Export generated public keys as CSV
│   └── descriptor                 Export output descriptors to file
├── import
│   ├── privkey                    Import private keys from database into Bitcoin Core wallet
│   └── fullpubkey                 Import full public key exported by sign wallet
├── sign
│   ├── signature                  Sign an unsigned transaction
│   └── musig2                     (BTC only) MuSig2 signing operations
│       ├── nonce                  Generate MuSig2 nonce (Round 1)
│       └── sign                   Create MuSig2 partial signature (Round 2)
└── api                            Chain node RPC commands (coin-specific)
    ├── btc / bch
    │   ├── encryptwallet          Encrypt the wallet with a passphrase
    │   ├── walletpassphrase       Unlock the wallet for N seconds
    │   ├── walletpassphrasechange Change the wallet passphrase
    │   ├── walletlock             Lock the wallet immediately
    │   ├── dumpwallet             Dump all wallet keys to a file
    │   └── importwallet           Import keys from a wallet dump file
    └── eth
        └── importrawkey           Import an Ethereum raw private key
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
│   └── fullpubkey                 Export full public key for import by keygen wallet
├── import
│   └── privkey                    Import authorization private key into database
├── sign
│   ├── signature                  Sign a multisig transaction
│   └── musig2                     (BTC only) MuSig2 signing operations
│       ├── nonce                  Generate MuSig2 nonce (Round 1)
│       └── sign                   Create MuSig2 partial signature (Round 2)
└── api                            Chain node RPC commands (coin-specific)
    ├── btc / bch
    │   ├── encryptwallet          Encrypt the wallet with a passphrase
    │   ├── walletpassphrase       Unlock the wallet for N seconds
    │   ├── walletpassphrasechange Change the wallet passphrase
    │   └── walletlock             Lock the wallet immediately
    └── eth
        └── importrawkey           Import an Ethereum raw private key
```

---

## Coin Support Matrix

| Command group            | BTC | BCH | ETH | XRP |
|--------------------------|:---:|:---:|:---:|:---:|
| watch import address     |  ✓  |  ✓  |  ✓  |  ✓  |
| watch import descriptor  |  ✓  |  –  |  –  |  –  |
| watch create *           |  ✓  |  ✓  |  ✓  |  ✓  |
| watch send tx            |  ✓  |  ✓  |  ✓  |  ✓  |
| watch send multisig      |  ✓  |  –  |  –  |  ✓  |
| watch monitor *          |  ✓  |  ✓  |  ✓  |  ✓  |
| watch api                |  ✓  |  ✓  |  ✓  |  ✓  |
| keygen create descriptor |  ✓  |  –  |  –  |  –  |
| keygen create multisig   |  ✓  |  –  |  –  |  –  |
| keygen sign musig2       |  ✓  |  –  |  –  |  –  |
| sign sign musig2         |  ✓  |  –  |  –  |  –  |

---

## Related Documents

| Document | Description |
|----------|-------------|
| [`docs/commands.md`](../../../docs/commands.md) | Detailed flags, options, and usage examples |
| [`.claude/skills/wallet-cli/SKILL.md`](../../../.claude/skills/wallet-cli/SKILL.md) | How to run wallet CLI commands |
| [`.claude/rules/internal/cli-structure.md`](../../../.claude/rules/internal/cli-structure.md) | CLI command hierarchy rules |
| [`docs/chains/btc/operations/wallet-flow.md`](../../../docs/chains/btc/operations/wallet-flow.md) | BTC end-to-end transaction workflow |
