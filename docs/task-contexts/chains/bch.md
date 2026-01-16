# BCH (Bitcoin Cash) Reference

## Overview

| 項目 | 値 |
|------|-----|
| トランザクションモデル | UTXO型 |
| 通信方式 | Bitcoin Cash RPC |
| アドレス形式 | P2PKH, P2SH (CashAddr形式: bitcoincash:q...) |
| 特殊機能 | CashAddr, BTCからのフォーク |

## ⚠️ BTCとの違い

BCHはBTCからのフォークですが、以下の重要な違いがあります：

| 項目 | BTC | BCH |
|------|-----|-----|
| アドレス形式 | Legacy/Bech32 | CashAddr |
| SegWit | ✅ | ❌ |
| Taproot | ✅ | ❌ |
| Descriptor | ✅ | ❌ |
| ブロックサイズ | 1MB (+ SegWit) | 32MB |
| 手数料単位 | sat/vB | sat/B |

## Directory Structure

### Use Case Layer

```
internal/application/usecase/
└── keygen/btc/          # BTC/BCH共通（一部）
    └── ...
```

**注意**: BCH固有のUse Caseは現在BTCと共有されている部分があります。実装時は違いに注意してください。

### Infrastructure Layer

```
internal/infrastructure/api/btc/bch/
├── bitcoin_cash.go      # クライアント初期化
├── account.go           # アカウント管理
└── address.go           # アドレス操作（CashAddr対応）
```

### CLI Layer

BCH固有のCLIコマンドはBTCと構造を共有していますが、`coin_type = "bch"` で区別されます。

## Key Concepts

### CashAddr Format

BCHはCashAddr形式のアドレスを使用：

```
# Legacy形式（互換性あり）
1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2

# CashAddr形式（推奨）
bitcoincash:qr95sy3j9xwd2ap32xkykttr4cvcu7as4y0qverfuy

# プレフィックスなし
qr95sy3j9xwd2ap32xkykttr4cvcu7as4y0qverfuy
```

### UTXO Model

BTCと同様にUTXO型トランザクションモデルを使用：

```
入力 (UTXO) → トランザクション → 出力 (新しいUTXO)

- 入力: 以前のトランザクションの未使用出力
- 出力: 送金先アドレス + お釣りアドレス
- 手数料: 入力合計 - 出力合計
```

### Address Types

| タイプ | CashAddr | Legacy | 説明 |
|--------|----------|--------|------|
| P2PKH | q... | 1... | 通常アドレス |
| P2SH | p... | 3... | Script Hash (マルチシグ等) |

**注意**: SegWit (bc1...) およびTaproot (bc1p...) はBCHでは使用不可。

## Implementation Patterns

### BitcoinCash構造体の設計

BCHのAPIクライアントは**BTCのインスタンスを埋め込んで拡張**する設計になっています：

```go
// internal/infrastructure/api/btc/bch/bitcoin_cash.go
type BitcoinCash struct {
    btc.Bitcoin  // BTCの実装を埋め込み
}
```

この設計により：
- BTCと共通のメソッドは自動的に継承される
- BCH固有の処理が必要な場合は**メソッドをオーバーライド**して対応

### ⚠️ 重要: BTC APIの問題修正パターン

BTCのAPI実装に問題がある場合、**BTCのコードを直接修正するのではなく**、BCH側でメソッドをオーバーライドして対応します：

```go
// ❌ やってはいけない: BTCのコードを直接修正
// internal/infrastructure/api/btc/btc/address.go を編集

// ✅ 正しいパターン: BCHでメソッドをオーバーライド
// internal/infrastructure/api/btc/bch/address.go
func (b *BitcoinCash) GetAddressInfo(addr string) (*dtobtc.AddressInfo, error) {
    // BCH固有の実装
    input, err := json.Marshal(addr)
    if err != nil {
        return nil, fmt.Errorf("fail to call json.Marchal() in bch: %w", err)
    }
    rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
    if err != nil {
        return nil, fmt.Errorf("fail to call json.RawRequest(getaddressinfo) %s in bch: %w", addr, err)
    }

    // BCH固有のレスポンス型を使用
    infoResult := GetAddressInfoResult{}
    err = json.Unmarshal(rawResult, &infoResult)
    if err != nil {
        return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult) in bch: %w", err)
    }

    // BCH型からBTC型に変換し、最終的にDTOに変換
    btcResult := &btc.GetAddressInfoResult{
        Address:      infoResult.Address,
        ScriptPubKey: infoResult.ScriptPubKey,
        // ... BCH固有のマッピング
        Iswitness:    false,  // BCHはSegWit非対応
    }

    return btc.ToAddressInfo(btcResult), nil
}
```

**この設計の理由**:
1. BTCの実装はBTC専用として安定させる
2. BCH固有の差異はBCHレイヤーで吸収
3. 両方のチェーンへの影響を分離

### BCH固有の実装が必要なケース

以下の場合はBCH側でオーバーライドが必要：

| ケース | 理由 |
|--------|------|
| レスポンス構造の違い | BCHノードのAPIレスポンスがBTCと異なる |
| フィールドの有無 | SegWit/Taproot関連フィールドがない |
| アドレス形式 | CashAddr形式の処理 |
| 手数料計算 | sat/B vs sat/vB |
| 署名方式 | リプレイプロテクション |

### BTCとの共通化

BCHとBTCは多くのロジックを共有できますが、以下の点で分離が必要：

```go
// ✅ 共通化可能（埋め込みで自動継承）
- UTXO選択ロジック
- トランザクション構造の基本部分
- 署名アルゴリズム（ECDSA）
- 多くのRPCメソッド

// ❌ BCH側でオーバーライドが必要
- アドレス関連API（GetAddressInfo等）
- 手数料関連API
- SegWit/Taproot関連機能（BCHでは無効化）
```

### アドレス変換

```go
// internal/infrastructure/api/btc/bch/address.go
// CashAddr ↔ Legacy変換が必要な場合の処理
```

## Configuration

```toml
# config/wallet/bch_watch.toml
coin_type = "bch"
network_type = "mainnet"  # mainnet, testnet

[bitcoin]  # Bitcoin Cash Node
host = "localhost"
port = 8332
user = "user"
pass = "password"
```

## Testing

```bash
# BCH Infrastructureテスト
go test ./internal/infrastructure/api/btc/bch/...

# BCH関連の設定ファイル確認
ls config/wallet/bch_*.toml
```

## Related Documentation

- [BCH README](../../../../docs/crypto/bch/README.md) - 詳細なBCHドキュメント
- [BTC/BCH Technical Guide](../../../../docs/crypto/btc/overview/technical-reference.md) - BTC/BCH技術比較

## Common Operations

| 操作 | Wallet | 実装場所 |
|------|--------|----------|
| アドレス生成 | Keygen | 共通ロジック + CashAddr変換 |
| トランザクション作成 | Watch | BTC共通 + BCH固有調整 |
| 署名 | Keygen/Sign | BTC共通 |
| 送信 | Watch | BTC共通 |

## ⚠️ Implementation Notes

1. **BTCコードを直接修正しない**: BCH固有の問題はBCH側でメソッドをオーバーライドして対応
2. **埋め込み構造の理解**: `BitcoinCash`は`btc.Bitcoin`を埋め込んでいる
3. **SegWit非対応**: BTCのSegWit/Taproot関連コードはBCHでは使用しない
4. **Descriptor非対応**: BTCのDescriptor機能はBCHでは使用しない
5. **アドレス形式**: CashAddr形式を適切に処理する
6. **手数料計算**: sat/B (バイト) で計算、vByte (仮想バイト) ではない
7. **チェーンID**: リプレイプロテクションのためBTCと異なる署名が必要

## Comparison with BTC Implementation

### 埋め込みによる継承関係

```
btc.Bitcoin (internal/infrastructure/api/btc/btc/)
    ↑ 埋め込み
BitcoinCash (internal/infrastructure/api/btc/bch/)
    → BTCのメソッドを継承
    → 必要に応じてオーバーライド
```

### 参照・修正パターン

```
BTC実装を参考にする場合:

✅ 自動継承される（そのまま使用可能）:
- internal/infrastructure/api/btc/btc/unspent.go (UTXO取得)
- internal/infrastructure/api/btc/btc/transaction.go (トランザクション基本)
- internal/infrastructure/api/btc/btc/balance.go (残高取得)

✅ 参考にできる（Use Case層）:
- internal/application/usecase/keygen/btc/sign_transaction.go (署名基本)

⚠️ BCHでオーバーライド済み:
- internal/infrastructure/api/btc/bch/address.go (GetAddressInfo)
- 必要に応じて追加のオーバーライドを実装

❌ BCHでは使用しない:
- internal/infrastructure/api/btc/btc/descriptor*.go (Descriptor非対応)
- internal/infrastructure/api/btc/btc/psbt.go (PSBT非対応)
- internal/infrastructure/api/btc/btc/musig2.go (MuSig2非対応)
- internal/application/usecase/*/btc/*musig2*.go
- internal/application/usecase/*/btc/*descriptor*.go
```

### BCH固有実装を追加する場合

```go
// 1. bch/ ディレクトリに新しいファイルを作成
// internal/infrastructure/api/btc/bch/new_feature.go

package bch

// 2. BitcoinCashのメソッドとして実装
func (b *BitcoinCash) SomeMethod() error {
    // BCH固有の実装
}

// または、BTCのメソッドをオーバーライド
func (b *BitcoinCash) ExistingBTCMethod() (*Result, error) {
    // BCH向けにカスタマイズした実装
}
```

