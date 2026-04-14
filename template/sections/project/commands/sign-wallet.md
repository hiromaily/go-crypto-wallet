### Sign Wallet Commands

Sign Wallet operates offline as a cold wallet. It provides subsequent signatures (2nd, 3rd, etc.) for multisig
transactions. Each authorization operator uses their own Sign Wallet instance.

#### Create Commands

##### `sign create seed`

Creates a seed for the authorization account wallet. If `--seed` is provided, it will be stored instead of
generating a new one.

**Options:**

- `--seed <string>` - Seed value to store (development use only)

**Example:**

```bash
sign create seed
```

##### `sign create hdkey`

Creates HD keys for the authorization account.

**Example:**

```bash
sign create hdkey
```

#### Export Commands

##### `sign export fullpubkey`

Exports full public key addresses as a CSV file for import into Keygen Wallet. These are used to create multisig addresses.

**Example:**

```bash
sign export fullpubkey
```

#### Import Commands

##### `sign import privkey`

Imports generated private keys for the authorization account into the database.

**Example:**

```bash
sign import privkey
```

#### Sign Commands

##### `sign sign signature`

Signs a transaction that has already been signed by Keygen Wallet (or other Sign Wallets). This provides
subsequent signatures for multisig addresses.

**Options:**

- `--file <path>` - Path to the signed transaction file

**Example:**

```bash
sign sign signature --file data/tx/btc/tx_signed1_1234567890.json
```

#### API Commands

API commands are coin-specific and dynamically configured based on the `--coin` flag.

##### Bitcoin/Bitcoin Cash API (`sign api` with `--coin btc` or `--coin bch`)

- `sign api encryptwallet` - Encrypts the wallet with a passphrase
  - `--passphrase <string>` - Passphrase for encryption
- `sign api walletpassphrase` - Stores the wallet decryption key in memory
  - `--passphrase <string>` - Passphrase to unlock the wallet
- `sign api walletpassphrasechange` - Changes the wallet passphrase
  - `--old <string>` - Old passphrase
  - `--new <string>` - New passphrase
- `sign api walletlock` - Locks the wallet by removing the encryption key from memory

##### Ethereum API (`sign api` with `--coin eth`)

- `sign api clientversion` - Get client version
- `sign api nodeinfo` - Get node information
- `sign api syncing` - Get synchronization status
- `sign api netversion` - Get network version
