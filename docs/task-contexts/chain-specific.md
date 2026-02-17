---
task_type: chain-specific
description: 暗号通貨固有のタスク処理用コンテキスト
version: 1.0.0
---

# Task: Chain-Specific (暗号通貨固有タスク)

## Overview

このプロジェクトは複数の暗号通貨（BTC, BCH, ETH, XRP）をサポートするウォレットアプリケーションです。多くのタスクは特定の暗号通貨に対するものであり、正しいチェーンを識別し、そのチェーン固有の仕様に従って実装する必要があります。

## When to Use

- アドレス生成、トランザクション作成、署名、送信に関するタスク
- 特定の暗号通貨のAPIやUse Caseを実装する場合
- チェーン固有の仕様や制約を考慮する必要がある場合
- 既存のチェーン実装を参照して新機能を追加する場合

## Chain Identification

### Step 1: タスクからチェーンを特定

ユーザーの依頼から対象の暗号通貨を特定します：

| キーワード | チェーン | 備考 |
|-----------|---------|------|
| Bitcoin, BTC, ビットコイン | BTC | UTXO型、Taproot/SegWit対応 |
| Bitcoin Cash, BCH, ビットコインキャッシュ | BCH | UTXO型、BTCフォーク |
| Ethereum, ETH, イーサリアム | ETH | アカウント型、ERC-20対応 |
| Ripple, XRP, リップル | XRP | アカウント型、gRPC通信 |
| ERC-20, トークン | ETH | ETHスマートコントラクト |

### Step 2: チェーン固有のコンテキストをロード

特定されたチェーンに応じて、追加のコンテキストを読み込みます：

| Chain | Reference Document |
|-------|-------------------|
| BTC | [BTC Documentation](../chains/btc/README.md) |
| BCH | [BCH Documentation](../chains/bch/README.md) |
| ETH | [ETH Documentation](../chains/eth/README.md) |
| XRP | [XRP Documentation](../chains/xrp/README.md) |

## Required Context Documents

### Must Read (必須)

| Document | Path | Purpose |
|----------|------|---------|
| Multi-Chain Support | `docs/guidelines/multi-chain.md` | マルチチェーンアーキテクチャ |
| Architecture | `docs/guidelines/architecture.md` | Clean Architecture原則 |

### Chain-Specific (チェーン特定後)

| Chain | Document | Path |
|-------|----------|------|
| BTC | BTC Documentation | `docs/chains/btc/README.md` |
| BCH | BCH Documentation | `docs/chains/bch/README.md` |
| ETH | ETH Documentation | `docs/chains/eth/README.md` |
| XRP | XRP Documentation | `docs/chains/xrp/README.md` |

## Wallet Types

このプロジェクトは3種類のウォレットを実装しています：

### Watch Wallet (オンライン)

- **役割**: トランザクション作成・送信、アドレス監視
- **鍵**: 公開鍵のみ（秘密鍵なし）
- **セキュリティ**: 秘密鍵がオンラインに露出しない

```
internal/application/usecase/watch/{coin}/
  - create_transaction.go    # トランザクション作成
  - send_transaction.go      # トランザクション送信
  - monitor_transaction.go   # トランザクション監視
  - import_address.go        # アドレスインポート
```

### Keygen Wallet (オフライン)

- **役割**: 鍵生成、最初の署名（マルチシグ）
- **鍵**: 秘密鍵を生成・保管
- **セキュリティ**: 最高（オフライン、エアギャップ推奨）

```
internal/application/usecase/keygen/{coin}/
  - generate_*.go            # 鍵生成
  - sign_transaction.go      # トランザクション署名
  - export_*.go              # エクスポート
  - import_*.go              # インポート
```

### Sign Wallet (オフライン)

- **役割**: 2番目以降の署名（マルチシグ）
- **鍵**: 署名用の秘密鍵を保管
- **セキュリティ**: 高（オフライン）

```
internal/application/usecase/sign/{coin}/
  - sign_transaction.go      # トランザクション署名
  - export_fullpubkey.go     # 公開鍵エクスポート
```

## Directory Structure by Chain

### Use Case Layer

```
internal/application/usecase/
├── keygen/
│   ├── btc/          # BTC鍵生成・署名
│   ├── eth/          # ETH鍵生成・署名
│   ├── xrp/          # XRP鍵生成・署名
│   └── shared/       # 共通ロジック（HD Wallet, Seed生成）
├── sign/
│   ├── btc/          # BTC署名
│   ├── eth/          # ETH署名
│   ├── xrp/          # XRP署名
│   └── shared/       # 共通ロジック
└── watch/
    ├── btc/          # BTCトランザクション
    ├── eth/          # ETHトランザクション
    ├── xrp/          # XRPトランザクション
    └── shared/       # 共通ロジック
```

### Infrastructure Layer (API)

```
internal/infrastructure/api/
├── bitcoin/
│   ├── btc/          # Bitcoin Core RPC
│   ├── bch/          # Bitcoin Cash RPC
│   └── connection.go # 共通接続
├── ethereum/
│   ├── eth/          # Ethereum JSON-RPC
│   ├── erc20/        # ERC-20トークン
│   └── connection.go # 共通接続
└── ripple/
    └── xrp/          # Ripple gRPC
```

### CLI Layer

```
internal/interface-adapters/cli/
├── keygen/
│   ├── api/btc/      # BTC API呼び出し
│   ├── api/eth/      # ETH API呼び出し
│   ├── btc/          # BTC固有コマンド
│   └── ...
├── sign/
│   └── ...
└── watch/
    ├── api/btc/      # BTC API呼び出し
    ├── api/eth/      # ETH API呼び出し
    ├── api/xrp/      # XRP API呼び出し
    ├── btc/          # BTC固有コマンド
    └── ...
```

## Chain-Specific Considerations

### BTC (Bitcoin)

| 特徴 | 詳細 |
|------|------|
| トランザクションモデル | UTXO型 |
| アドレスタイプ | P2PKH, P2SH, P2WPKH, P2WSH, P2TR (Taproot) |
| マルチシグ | P2SH, P2WSH, MuSig2 (Taproot) |
| 特殊機能 | Descriptor, PSBT, Taproot, MuSig2 |

### BCH (Bitcoin Cash)

| 特徴 | 詳細 |
|------|------|
| トランザクションモデル | UTXO型（BTCフォーク） |
| アドレスタイプ | P2PKH, P2SH (CashAddr形式) |
| マルチシグ | P2SH |
| 注意 | Taproot非対応、BTCとの違いに注意 |

### ETH (Ethereum)

| 特徴 | 詳細 |
|------|------|
| トランザクションモデル | アカウント型 |
| トランザクション | Nonce管理、Gas Price/Limit |
| 特殊機能 | ERC-20トークン、スマートコントラクト |
| 署名 | Keystore、秘密鍵署名 |

### XRP (Ripple)

| 特徴 | 詳細 |
|------|------|
| トランザクションモデル | アカウント型 |
| 特殊機能 | Sequence番号、Destination Tag |
| 注意 | アカウントリザーブ要件 |

## Implementation Workflow

### 新しい機能をチェーンに追加する場合

1. **既存実装の確認**

   ```bash
   # 同じチェーンの既存Use Caseを確認
   ls internal/application/usecase/{wallet-type}/{coin}/

   # 他のチェーンの同等機能を確認（参考）
   ls internal/application/usecase/{wallet-type}/*/
   ```

2. **レイヤー順に実装**

   ```
   1. Domain (必要な場合)
   2. Ports (インターフェース定義)
   3. Use Case (ビジネスロジック)
   4. Infrastructure (API呼び出し)
   5. CLI (コマンド)
   6. DI (依存注入)
   ```

3. **チェーン固有のテスト**
   - 単体テスト: モック使用
   - 統合テスト: テストネット使用

## Pre-Task Checklist

- [ ] 対象の暗号通貨を正しく特定したか
- [ ] チェーン固有のドキュメントを読み込んだか
- [ ] 既存の同等実装を確認したか
- [ ] 正しいディレクトリに実装しようとしているか
- [ ] チェーン固有の仕様（UTXO vs アカウント）を理解しているか
- [ ] ウォレットタイプ（watch/keygen/sign）を正しく選択したか

## Verification Commands

```bash
# ビルド確認（全チェーン）
make check-build

# テスト実行
make gotest

# 特定チェーンのテスト（例: BTC）
go test ./internal/application/usecase/watch/btc/...
go test ./internal/infrastructure/api/btc/btc/...
```

## Examples

### Example 1: BTCのトランザクション作成機能

```
User: "BTCのトランザクション作成ロジックを修正して"

Agent Actions:
1. チェーン特定: BTC
2. chains/btc.md を読み込み
3. internal/application/usecase/watch/btc/create_transaction.go を確認
4. UTXO型トランザクションの仕様を確認
5. 修正を実装
6. テスト実行
```

### Example 2: ETHの新しいUse Case追加

```
User: "ETHにトークン残高確認のUse Caseを追加して"

Agent Actions:
1. チェーン特定: ETH
2. chains/eth.md を読み込み
3. internal/application/usecase/watch/eth/ の既存実装を確認
4. ERC-20関連のInfrastructure確認
5. Use Case実装
6. CLI実装
7. テスト追加
```

### Example 3: XRPの署名機能修正

```
User: "XRPの署名処理でエラーが発生する"

Agent Actions:
1. チェーン特定: XRP
2. chains/xrp.md を読み込み
3. internal/application/usecase/sign/xrp/sign_transaction.go を確認
4. gRPC通信の問題か、署名ロジックの問題かを特定
5. 修正を実装
6. テスト実行
```

## Common Mistakes to Avoid

### ❌ チェーンの混同

```
問題: BTCの実装をETHにそのままコピー
理由: UTXO型とアカウント型で根本的に異なる

対策: チェーン固有の仕様を確認してから実装
```

### ❌ 間違ったディレクトリへの実装

```
問題: watch/btc/ に署名ロジックを実装
理由: watch walletは秘密鍵を持たない

対策: ウォレットタイプの役割を確認
```

### ❌ 共通化すべきでないものの共通化

```
問題: チェーン固有のロジックを shared/ に配置
理由: 各チェーンの仕様が異なる

対策: 共通化は慎重に、チェーン固有は各ディレクトリへ
```

## Related Documents

- [Multi-Chain Support](../guidelines/multi-chain.md) - マルチチェーンアーキテクチャ詳細
- [Architecture Guidelines](../guidelines/architecture.md) - レイヤー構造
- [Chain Documentation](../chains/) - 各チェーンの詳細ドキュメント
