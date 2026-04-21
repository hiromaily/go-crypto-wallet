## ERC-20 Token Support

### Overview

ERC-20 token transfers are supported via the `erc20` infrastructure package.

**Key Difference from ETH transfers:**

| | ETH Transfer | ERC-20 Transfer |
|-|-------------|-----------------|
| **To address** | Recipient | Contract address |
| **Value** | Transfer amount | 0 (always) |
| **Data** | Empty | ABI-encoded `transfer()` call |
| **Fee payer** | Sender | Master address |

### Function Selector

ERC-20 `transfer(address,uint256)` is encoded as:

```go
// Function selector = first 4 bytes of Keccak256("transfer(address,uint256)")
selector := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]

// Call data = selector + padded address (32 bytes) + padded amount (32 bytes)
data = append(selector, paddedAddress...)
data = append(data, paddedAmount...)
```

### Limitations

- Requires prior `approve()` transaction for some flows
- Master address pays ETH gas fees (not token holder)
- No ERC-721 or ERC-1155 support

---
