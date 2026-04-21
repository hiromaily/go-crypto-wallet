## Address Creation

### Creating MuSig2 Addresses

#### Command

```bash
# On Keygen wallet
./keygen create musig2-address --account <account> [--count <number>]
```

#### Parameters

- `--account`: Account type (payment, stored, etc.)
- `--count`: Number of addresses to create (default: 10)

#### Example

```bash
# Create 20 MuSig2 addresses for payment account
./keygen create musig2-address --account payment --count 20
```

#### Output

```
✓ MuSig2 Taproot addresses created successfully
Account: payment
Addresses created: 20
Address type: P2TR (Taproot)
Multisig type: MuSig2 3-of-3
Public keys: 3 (Keygen + Sign1 + Sign2)
Sample addresses:
  tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxq...
  tb1pq8r3t5ys7whg3vz2nq4x6kj7h5pq9c8r3t5ys7whg3vz2nq4x6k...
Next: Export addresses to Watch wallet
```

### Address Properties

MuSig2 addresses have these characteristics:

- **Address Type**: P2TR (Taproot)
- **Public Key**: Aggregated from all signers
- **On-Chain**: Looks like single-sig address
- **Privacy**: Maximum (indistinguishable from single-sig)
- **Efficiency**: Smallest possible address type

### Exporting Addresses

```bash
# On Keygen wallet
./keygen export address --account payment --file payment_musig2_addresses.txt

# Transfer to Watch wallet
# On Watch wallet
./watch import address --file payment_musig2_addresses.txt
```

---
