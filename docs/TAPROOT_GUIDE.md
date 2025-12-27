# Taproot User Guide

This guide explains how to use Taproot (BIP341/BIP86) addresses in the go-crypto-wallet system.

## Table of Contents

- [What is Taproot?](#what-is-taproot)
- [Why Use Taproot?](#why-use-taproot)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Getting Started](#getting-started)
- [Workflow Examples](#workflow-examples)
- [Migration Guide](#migration-guide)
- [Best Practices](#best-practices)
- [FAQ](#faq)
- [Troubleshooting](#troubleshooting)

## What is Taproot?

Taproot is Bitcoin's latest major upgrade (activated November 2021) that introduces a new address format and signature scheme:

- **Address Format:** P2TR (Pay-to-Taproot) using bech32m encoding
- **Address Prefix:** `bc1p` for mainnet, `tb1p` for testnet/signet
- **Signature Scheme:** Schnorr signatures (BIP340) instead of ECDSA
- **Derivation Path:** BIP86 (`m/86'/0'/account'/0/index`)

### Technical Improvements

1. **Privacy:** All Taproot spends look identical on-chain (single-sig and multisig indistinguishable)
2. **Efficiency:** Smaller transaction sizes = lower fees (30-50% reduction)
3. **Smart Contracts:** Script path spending enables complex conditions
4. **Signature Aggregation:** Multiple signatures can be combined (MuSig2)

## Why Use Taproot?

### Benefits

✅ **Lower Transaction Fees** - 30-50% smaller transaction sizes
✅ **Enhanced Privacy** - All spends look the same on-chain
✅ **Future-Proof** - Latest Bitcoin address standard
✅ **Faster Validation** - Schnorr signatures verify faster than ECDSA
✅ **Better Multisig** - More efficient than traditional multisig

### When to Use Taproot

- ✅ New wallets and addresses
- ✅ When transaction fee savings matter
- ✅ When privacy is important
- ✅ For multisig setups (more efficient)

### When NOT to Use Taproot

- ⚠️ If you need compatibility with services that don't support Taproot yet
- ⚠️ If you're using Bitcoin Core older than v22.0
- ⚠️ If you need to support very old wallet software

## Prerequisites

### Required Software Versions

1. **Bitcoin Core v22.0 or later** (for Taproot/Schnorr support)
   ```bash
   bitcoin-cli --version
   # Should output: Bitcoin Core version v22.0.0 or higher
   ```

2. **go-crypto-wallet v5.0.0 or later**
   ```bash
   ./keygen --version
   ./watch --version
   ./sign --version
   ```

3. **MySQL 5.7 or later** (for address/transaction storage)

### Recommended Setup

- **Operating System:** Linux, macOS, or Windows
- **Network:** Bitcoin Mainnet, Testnet3, or Signet
- **Disk Space:** Varies by network (Mainnet requires significant space)

## Configuration

### Step 1: Configure Key Type and Address Type

Edit your wallet configuration file to enable Taproot:

**Keygen Wallet** (`data/config/btc_keygen.toml`):
```toml
# BIP86 is required for Taproot
key_type = "bip86"  # bip44, bip49, bip84, bip86, musig2

# Set address type to taproot
address_type = "taproot"  # legacy, p2sh-segwit, bech32, taproot

[bitcoin]
host = "127.0.0.1:18332"
user = "your_rpc_user"
pass = "your_rpc_password"
network_type = "testnet3"  # mainnet, testnet3, regtest, signet
```

**Watch Wallet** (`data/config/btc_watch.toml`):
```toml
address_type = "taproot"

[bitcoin]
host = "127.0.0.1:18332"
user = "your_rpc_user"
pass = "your_rpc_password"
network_type = "testnet3"
```

**Sign Wallet** (`data/config/btc_sign.toml`):
```toml
address_type = "taproot"

[bitcoin]
host = "127.0.0.1:20332"
user = "your_rpc_user"
pass = "your_rpc_password"
network_type = "testnet3"
```

### Step 2: Verify Bitcoin Core Configuration

Ensure Bitcoin Core is configured correctly:

```bash
# Check Bitcoin Core version (must be v22.0+)
bitcoin-cli --version

# Check if Taproot is activated (mainnet)
bitcoin-cli getblockchaininfo | grep -A 5 taproot

# Create wallet with descriptor support (recommended)
bitcoin-cli createwallet "taproot_wallet" false false "" false true
```

### Step 3: Verify Database Schema

Ensure your database has the `taproot_address` column:

```sql
-- Check if taproot_address column exists
DESCRIBE account_key;

-- Should show:
-- taproot_address varchar(255) NULL
```

If missing, run the migration script (already included in Phase 1a).

## Getting Started

### Quick Start Example

This example walks through creating a Taproot wallet from scratch:

```bash
# 1. Set environment variables
export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen.toml
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml
export BTC_ACCOUNT_CONF=./data/config/account.toml

# 2. Generate seed (Keygen wallet - OFFLINE)
./keygen --coin btc create seed

# 3. Generate Taproot keys (Keygen wallet - OFFLINE)
./keygen --coin btc create hdkey --account client --count 10

# 4. Export addresses (Keygen wallet - OFFLINE)
./keygen --coin btc export address --account client
# Output: ./data/address/btc/client_1234567890.csv

# 5. Import addresses to Watch wallet (Watch wallet - ONLINE)
./watch --coin btc import address --account client \
  --filepath ./data/address/btc/client_1234567890.csv

# 6. Check imported addresses
./watch --coin btc api getaddressesbylabel --label client
```

### Expected Output

After generating Taproot keys, you should see addresses like:

```
Mainnet: bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr
Testnet: tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c
```

**Taproot Address Characteristics:**
- Prefix: `bc1p` (mainnet) or `tb1p` (testnet/signet)
- Length: 62 characters
- Encoding: bech32m (lowercase only)

## Workflow Examples

### Example 1: Receiving Funds to Taproot Address

**Scenario:** Customer deposits funds to your Taproot address

```bash
# 1. Generate client account Taproot addresses (Keygen - OFFLINE)
export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen.toml
./keygen --coin btc create hdkey --account client --count 100

# 2. Export addresses (Keygen - OFFLINE)
./keygen --coin btc export address --account client
# Output: ./data/address/btc/client_1234567890.csv

# 3. Transfer address file to Watch wallet (secure transfer)
# Copy client_1234567890.csv to Watch wallet system

# 4. Import addresses (Watch - ONLINE)
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml
./watch --coin btc import address --account client \
  --filepath ./data/address/btc/client_1234567890.csv

# 5. Monitor for incoming transactions
./watch --coin btc monitor transaction

# 6. Create deposit transaction when funds are confirmed
./watch --coin btc create transaction --account deposit
# Output: ./data/tx/btc/deposit_1_unsigned_0_1234567890.tx
```

### Example 2: Sending Funds from Taproot Address (Single-Sig)

**Scenario:** Send funds from deposit account to payment account

```bash
# 1. Create unsigned transaction (Watch - ONLINE)
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml
./watch --coin btc create transaction --account deposit
# Output: ./data/tx/btc/deposit_1_unsigned_0_1234567890.tx

# 2. Transfer transaction file to Keygen wallet (secure transfer)
# Copy deposit_1_unsigned_0_1234567890.tx to Keygen system

# 3. Sign transaction with Schnorr signature (Keygen - OFFLINE)
export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen.toml
./keygen --coin btc sign \
  --file ./data/tx/btc/deposit_1_unsigned_0_1234567890.tx
# Output: ./data/tx/btc/deposit_1_signed_0_1234567890.tx

# 4. Transfer signed transaction back to Watch wallet

# 5. Send transaction (Watch - ONLINE)
./watch --coin btc send transaction --account deposit \
  --file ./data/tx/btc/deposit_1_signed_0_1234567890.tx
```

### Example 3: Multisig Taproot Transaction

**Scenario:** Send funds from payment account (requires 2-of-3 signatures)

```bash
# 1. Create unsigned transaction (Watch - ONLINE)
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml
./watch --coin btc create transaction --account payment
# Output: ./data/tx/btc/payment_5_unsigned_0_1234567890.tx

# 2. First signature (Keygen - OFFLINE)
export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen.toml
./keygen --coin btc sign \
  --file ./data/tx/btc/payment_5_unsigned_0_1234567890.tx
# Output: ./data/tx/btc/payment_5_unsigned_1_1234567890.tx (still unsigned - needs more sigs)

# 3. Second signature (Sign wallet - OFFLINE)
export BTC_SIGN_WALLET_CONF=./data/config/btc_sign.toml
./sign --coin btc sign \
  --file ./data/tx/btc/payment_5_unsigned_1_1234567890.tx
# Output: ./data/tx/btc/payment_5_signed_0_1234567890.tx (now fully signed)

# 4. Send transaction (Watch - ONLINE)
./watch --coin btc send transaction --account payment \
  --file ./data/tx/btc/payment_5_signed_0_1234567890.tx
```

### Example 4: Creating Payment Request with Taproot

```bash
# 1. Create payment request (Watch - ONLINE)
export BTC_WATCH_WALLET_CONF=./data/config/btc_watch.toml
./watch --coin btc create payment-request \
  --address bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr \
  --amount 0.001

# 2. Monitor payment status
./watch --coin btc monitor transaction
```

## Migration Guide

### Migrating from Legacy Addresses to Taproot

**Important:** You cannot convert existing addresses to Taproot. You must generate new addresses and migrate funds.

#### Step 1: Generate New Taproot Addresses

```bash
# Update configuration
sed -i 's/address_type = "bech32"/address_type = "taproot"/' \
  data/config/btc_keygen.toml
sed -i 's/key_type = "bip84"/key_type = "bip86"/' \
  data/config/btc_keygen.toml

# Generate new Taproot addresses
./keygen --coin btc create hdkey --account client --count 100
./keygen --coin btc export address --account client
```

#### Step 2: Import to Watch Wallet

```bash
# Import new Taproot addresses
./watch --coin btc import address --account client \
  --filepath ./data/address/btc/client_1234567890.csv
```

#### Step 3: Gradual Migration Strategy

**Option A: Immediate Migration**
- Stop using old addresses immediately
- Transfer all funds from old addresses to new Taproot addresses
- Update all systems to use new addresses

**Option B: Gradual Migration**
1. Start issuing new Taproot addresses to new users
2. Keep accepting payments to old addresses
3. Periodically consolidate old address balances to Taproot addresses
4. Phase out old addresses over 6-12 months

#### Step 4: Update External Systems

```bash
# Export new Taproot addresses for integration
./keygen --coin btc export address --account client

# Update:
# - Payment processors
# - Customer databases
# - Invoicing systems
# - Exchange integrations
```

### Compatibility Matrix

| Address Type | Your Wallet | Bitcoin Core | Network | Status |
|--------------|-------------|--------------|---------|--------|
| Taproot (bc1p) | v5.0.0+ | v22.0+ | All | ✅ Full Support |
| Bech32 (bc1q) | All | v0.16.0+ | All | ✅ Compatible |
| P2SH-SegWit (3) | All | v0.13.1+ | All | ✅ Compatible |
| Legacy (1) | All | All | All | ✅ Compatible |

## Best Practices

### Security

1. **Keep Keygen/Sign Wallets Offline**
   - Never connect Keygen or Sign wallets to the internet
   - Use air-gapped systems or hardware security modules (HSMs)
   - Transfer files via USB or QR codes only

2. **Secure File Transfers**
   ```bash
   # Encrypt transaction files
   gpg --encrypt --recipient your@email.com transaction.tx

   # Verify checksums
   sha256sum transaction.tx
   ```

3. **Backup Seed Phrases**
   - Store seed phrases in multiple secure locations
   - Use hardware wallets for seed storage
   - Test wallet recovery before storing funds

### Performance

1. **Batch Operations**
   ```bash
   # Generate addresses in batches
   ./keygen --coin btc create hdkey --account client --count 1000
   ```

2. **Monitor Transaction Pool**
   ```bash
   # Check mempool before sending
   bitcoin-cli getmempoolinfo
   ```

3. **Optimize Fee Estimation**
   - Taproot transactions are smaller = lower fees
   - Use dynamic fee estimation based on mempool

### Operational

1. **Test on Testnet First**
   ```bash
   # Always test new workflows on testnet
   export BTC_KEYGEN_WALLET_CONF=./data/config/btc_keygen_testnet.toml
   ```

2. **Keep Audit Logs**
   ```bash
   # Log all wallet operations
   ./keygen --coin btc create hdkey --account client --count 10 \
     2>&1 | tee -a keygen_audit.log
   ```

3. **Regular Backups**
   ```bash
   # Backup database
   mysqldump -u user -p wallet_db > backup_$(date +%Y%m%d).sql

   # Backup address files
   tar -czf addresses_backup_$(date +%Y%m%d).tar.gz data/address/
   ```

## FAQ

### General Questions

**Q: Can I receive funds to Taproot addresses from any wallet?**
A: Yes, most modern wallets support sending to Taproot addresses. However, some older wallets may not recognize bc1p addresses.

**Q: Are Taproot addresses case-sensitive?**
A: No, Taproot addresses use bech32m encoding which is case-insensitive. However, they should always be displayed in lowercase.

**Q: Can I convert my existing bech32 (bc1q) addresses to Taproot (bc1p)?**
A: No, you must generate new Taproot addresses. Each address type has a different derivation path and key structure.

**Q: Do I need to upgrade my entire system at once?**
A: No, you can run Taproot and legacy address types side-by-side. Migrate gradually.

### Technical Questions

**Q: What is the difference between bech32 and bech32m?**
A: Bech32m is an improved version that fixes a checksum issue in the original bech32 encoding. Taproot uses bech32m.

**Q: Can I use Taproot for multisig?**
A: Yes, but traditional multisig (multiple keys, threshold signing) works the same way. True Taproot multisig (MuSig2 key aggregation) is more efficient but not yet implemented in this wallet.

**Q: How much smaller are Taproot transactions?**
A: Typical savings:
- Single-sig: 30-40% smaller
- 2-of-3 multisig: 40-50% smaller
- Larger multisig: Even greater savings

**Q: Do I need to change my backup procedures?**
A: No, seed phrases work the same way. However, ensure you note the derivation path (BIP86) for recovery.

**Q: Can I mix Taproot and legacy UTXOs in one transaction?**
A: Yes, you can spend from multiple address types in a single transaction.

### Troubleshooting Questions

**Q: Why do I get "fail to call NewAddressTaproot()" errors?**
A: This usually means:
1. Bitcoin Core version is older than v22.0
2. Bitcoin Core is not properly configured
3. RPC connection issue

**Q: Why are my Taproot addresses not showing up?**
A: Check:
1. Configuration: `address_type = "taproot"` and `key_type = "bip86"`
2. Database schema has `taproot_address` column
3. Address import completed successfully

**Q: Transaction signatures failing?**
A: Ensure:
1. Bitcoin Core v22.0+ is running
2. Transaction file includes previous transaction data
3. Correct wallet is unlocked in Bitcoin Core

## Troubleshooting

### Issue: "Address type not supported"

**Symptoms:**
```
Error: address type 'taproot' not recognized
```

**Solutions:**
1. Verify configuration file:
   ```bash
   grep "address_type" data/config/btc_keygen.toml
   # Should show: address_type = "taproot"
   ```

2. Check wallet version:
   ```bash
   ./keygen --version
   # Should be v5.0.0 or later
   ```

### Issue: "Bitcoin Core RPC connection failed"

**Symptoms:**
```
Error: could not connect to Bitcoin Core RPC
```

**Solutions:**
1. Verify Bitcoin Core is running:
   ```bash
   bitcoin-cli getblockchaininfo
   ```

2. Check RPC credentials in config:
   ```toml
   [bitcoin]
   host = "127.0.0.1:18332"
   user = "your_rpc_user"
   pass = "your_rpc_password"
   ```

3. Test RPC connection:
   ```bash
   curl --user your_rpc_user:your_rpc_password \
     --data-binary '{"jsonrpc": "1.0", "id": "test", "method": "getblockchaininfo", "params": []}' \
     -H 'content-type: text/plain;' \
     http://127.0.0.1:18332/
   ```

### Issue: "Taproot address not generated"

**Symptoms:**
- Keys generated but no Taproot addresses in output
- Database shows NULL for `taproot_address` column

**Solutions:**
1. Verify `key_type` is set to `bip86`:
   ```bash
   grep "key_type" data/config/btc_keygen.toml
   # Should show: key_type = "bip86"
   ```

2. Check database schema:
   ```sql
   DESCRIBE account_key;
   -- Should include: taproot_address varchar(255) NULL
   ```

3. Regenerate keys with correct configuration

### Issue: "Transaction signature invalid"

**Symptoms:**
```
Error: non-mandatory-script-verify-flag (Signature must be zero for failed CHECK(MULTI)SIG operation)
```

**Solutions:**
1. Ensure Bitcoin Core v22.0+:
   ```bash
   bitcoin-cli --version
   ```

2. Verify transaction file includes previous outputs:
   ```bash
   cat ./data/tx/btc/payment_5_unsigned_0_*.tx
   # Should contain: hex,encoded_prevs_addrs
   ```

3. Check wallet is properly configured for Taproot

### Issue: "Database migration required"

**Symptoms:**
```
Error: column 'taproot_address' does not exist
```

**Solutions:**
1. Run database migration:
   ```sql
   ALTER TABLE account_key ADD COLUMN taproot_address VARCHAR(255) NULL
     AFTER bech32_address;
   ```

2. Restart wallet services

### Getting Help

If you encounter issues not covered here:

1. **Check Logs:**
   ```bash
   # Keygen wallet
   tail -f keygen.log

   # Watch wallet
   tail -f watch.log
   ```

2. **Enable Debug Logging:**
   ```toml
   [logger]
   level = "debug"  # change from "info" to "debug"
   ```

3. **Verify Configuration:**
   ```bash
   # Validate TOML syntax
   python3 -c "import toml; toml.load('data/config/btc_keygen.toml')"
   ```

4. **Test with Signet:**
   - Signet is easier to test than testnet
   - Free test coins available
   - Faster block times

5. **Report Issues:**
   - GitHub Issues: https://github.com/hiromaily/go-crypto-wallet/issues
   - Include: version, config (redact secrets), error messages, steps to reproduce

## Additional Resources

### Documentation
- [Taproot Testing Guide](./testing/TAPROOT_TESTING.md) - For developers running tests
- [BIP341 - Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki) - Technical specification
- [BIP86 - Key Derivation](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki) - Derivation path standard
- [BIP340 - Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki) - Signature scheme

### Tools
- [Bitcoin Core](https://bitcoincore.org/) - Full node software
- [Bitcoin Explorer](https://www.blockchain.com/explorer) - View transactions on-chain
- [Mempool.space](https://mempool.space/) - Fee estimation and mempool monitoring

### Community
- Bitcoin-dev mailing list
- Bitcoin Stack Exchange
- Bitcoin Core development

---

**Document Version:** Phase 1h
**Last Updated:** 2025-12-27
**Related Issue:** #89 - Implement Taproot Support (Phase 1)
**Minimum Wallet Version:** v5.0.0
**Minimum Bitcoin Core Version:** v22.0.0
