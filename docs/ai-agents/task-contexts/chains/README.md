# Chain-Specific References

このディレクトリには、各暗号通貨固有の実装リファレンスが含まれています。

## Quick Reference

| Chain | File | Transaction Model | Communication |
|-------|------|------------------|---------------|
| BTC | [btc.md](./btc.md) | UTXO | Bitcoin Core RPC |
| BCH | [bch.md](./bch.md) | UTXO | Bitcoin Cash RPC |
| ETH | [eth.md](./eth.md) | Account | JSON-RPC |
| XRP | [xrp.md](./xrp.md) | Account | gRPC |

## Usage

タスクを実行する際、対象の暗号通貨を特定したら、該当するリファレンスファイルを読み込んでください。

```
1. ユーザーの依頼からチェーンを特定
2. 該当するリファレンスファイル (*.md) を読み込む
3. ディレクトリ構造と既存実装を確認
4. チェーン固有の仕様に従って実装
```

## Directory Structure Overview

```
internal/
├── application/usecase/
│   ├── keygen/{btc,eth,xrp}/    # 鍵生成・署名
│   ├── sign/{btc,eth,xrp}/      # 署名
│   └── watch/{btc,eth,xrp}/     # トランザクション
├── infrastructure/api/
│   ├── bitcoin/{btc,bch}/       # BTC/BCH RPC
│   ├── ethereum/{eth,erc20}/    # ETH JSON-RPC
│   └── ripple/xrp/              # XRP gRPC
└── interface-adapters/cli/
    ├── keygen/api/{btc,eth}/
    ├── sign/
    └── watch/api/{btc,eth,xrp}/
```

## Comparison Table

| Feature | BTC | BCH | ETH | XRP |
|---------|-----|-----|-----|-----|
| Address Types | P2PKH, P2SH, Bech32, Taproot | P2PKH, P2SH (CashAddr) | 0x Hex | r... |
| Multisig | ✅ P2SH, P2WSH, MuSig2 | ✅ P2SH | ❌ (Contract) | ❌ |
| UTXO/Account | UTXO | UTXO | Account | Account |
| Fee Type | sat/vB | sat/B | Gas | Drops |
| Descriptor | ✅ | ❌ | ❌ | ❌ |
| Taproot | ✅ | ❌ | ❌ | ❌ |
| ERC-20 | ❌ | ❌ | ✅ | ❌ |

## Related Documents

- [Chain-Specific Task Context](../chain-specific.md) - チェーン固有タスクの処理方法
- [Multi-Chain Support](../../guidelines/multi-chain.md) - マルチチェーンアーキテクチャ

