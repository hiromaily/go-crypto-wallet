# BTC (Bitcoin) Reference

## Overview

| 項目 | 値 |
|------|-----|
| トランザクションモデル | UTXO型 |
| 通信方式 | Bitcoin Core RPC |
| アドレス形式 | P2PKH, P2SH, Bech32 (P2WPKH, P2WSH), Bech32m (P2TR) |
| 特殊機能 | Descriptor, PSBT, Taproot, MuSig2 |

## Directory Structure

### Use Case Layer

```
internal/application/usecase/
├── keygen/btc/
│   ├── create_multisig_address.go   # マルチシグアドレス作成
│   ├── create_musig2_address.go     # MuSig2アドレス作成
│   ├── export_address.go            # アドレスエクスポート
│   ├── export_descriptor.go         # Descriptorエクスポート
│   ├── generate_descriptor.go       # Descriptor生成
│   ├── import_fullpubkey.go         # 公開鍵インポート
│   ├── import_private_key.go        # 秘密鍵インポート
│   ├── musig2_nonce.go              # MuSig2 Nonce生成
│   ├── musig2_sign.go               # MuSig2署名
│   └── sign_transaction.go          # トランザクション署名
├── sign/btc/
│   ├── export_fullpubkey.go         # 公開鍵エクスポート
│   ├── import_private_key.go        # 秘密鍵インポート
│   ├── musig2_nonce.go              # MuSig2 Nonce
│   ├── musig2_sign.go               # MuSig2署名
│   └── sign_transaction.go          # トランザクション署名
└── watch/btc/
    ├── create_transaction.go        # トランザクション作成
    ├── import_address.go            # アドレスインポート
    ├── import_descriptor.go         # Descriptorインポート
    ├── monitor_transaction.go       # トランザクション監視
    ├── musig2_aggregate.go          # MuSig2署名集約
    └── send_transaction.go          # トランザクション送信
```

### Infrastructure Layer

```
internal/infrastructure/api/bitcoin/btc/
├── bitcoin.go           # クライアント初期化
├── account.go           # アカウント管理
├── address.go           # アドレス操作
├── amount.go            # 金額計算
├── balance.go           # 残高取得
├── block.go             # ブロック情報
├── descriptor*.go       # Descriptor関連（多数）
├── fee.go               # 手数料推定
├── import.go            # インポート機能
├── label.go             # ラベル管理
├── multisig.go          # マルチシグ
├── musig2.go            # MuSig2
├── psbt.go              # PSBT処理
├── transaction.go       # トランザクション
├── unspent.go           # UTXO管理
└── wallet.go            # ウォレット操作
```

### CLI Layer

```
internal/interface-adapters/cli/
├── keygen/
│   ├── api/btc/         # Bitcoin Core API呼び出し
│   ├── btc/descriptor.go
│   └── musig2.go
├── sign/
│   └── musig2.go
└── watch/
    ├── api/btc/         # Bitcoin Core API呼び出し
    ├── btc/descriptor.go
    └── musig2.go
```

## Key Concepts

### UTXO (Unspent Transaction Output)

BTCはUTXO型トランザクションモデルを使用：

```
入力 (UTXO) → トランザクション → 出力 (新しいUTXO)

- 入力: 以前のトランザクションの未使用出力
- 出力: 送金先アドレス + お釣りアドレス
- 手数料: 入力合計 - 出力合計
```

### Address Types

| タイプ | 形式 | 説明 |
|--------|------|------|
| P2PKH | 1... | Legacy (非推奨) |
| P2SH | 3... | Script Hash (マルチシグ等) |
| P2WPKH | bc1q... | Native SegWit (推奨) |
| P2WSH | bc1q... | SegWit Script Hash |
| P2TR | bc1p... | Taproot (最新) |

### Descriptor

Output Descriptorはアドレス/スクリプトを表現する標準形式：

```
# P2WPKH (SegWit single key)
wpkh([fingerprint/path]xpub/*)

# P2WSH multisig
wsh(multi(2,[fp1/path]xpub1/*,[fp2/path]xpub2/*))

# Taproot
tr([fingerprint/path]xpub/*)
```

### PSBT (Partially Signed Bitcoin Transaction)

マルチシグなどで部分署名を扱うフォーマット：

```
1. Watch Wallet: PSBT作成（未署名）
2. Keygen Wallet: 最初の署名を追加
3. Sign Wallet: 残りの署名を追加
4. Watch Wallet: 完了したPSBTをブロードキャスト
```

### MuSig2

Taprootでのn-of-nマルチシグ実装：

```
1. Nonce生成（各署名者）
2. Nonce集約
3. 部分署名生成
4. 署名集約
5. 最終署名完成
```

## Implementation Patterns

### トランザクション作成 (Watch Wallet)

```go
// internal/application/usecase/watch/btc/create_transaction.go
type CreateTransactionUseCase struct {
    bitcoinClient ports.BitcoinClient
    // ...
}

func (u *CreateTransactionUseCase) Execute(ctx context.Context, req *CreateTxRequest) (*CreateTxResponse, error) {
    // 1. UTXOを取得
    // 2. 入力を選択
    // 3. 出力を構築
    // 4. 手数料を計算
    // 5. PSBT/Rawトランザクションを生成
}
```

### トランザクション署名 (Keygen/Sign Wallet)

```go
// internal/application/usecase/keygen/btc/sign_transaction.go
type SignTransactionUseCase struct {
    keyStore ports.KeyStore
    // ...
}

func (u *SignTransactionUseCase) Execute(ctx context.Context, req *SignTxRequest) (*SignTxResponse, error) {
    // 1. 秘密鍵を取得
    // 2. トランザクションに署名
    // 3. 署名済みトランザクションを返す
}
```

## Configuration

```toml
# config/wallet/btc_watch.toml
coin_type = "btc"
network_type = "mainnet"  # mainnet, testnet, signet, regtest

[bitcoin]
host = "localhost"
port = 8332
user = "user"
pass = "password"
```

## Testing

```bash
# BTC Use Caseテスト
go test ./internal/application/usecase/watch/btc/...
go test ./internal/application/usecase/keygen/btc/...
go test ./internal/application/usecase/sign/btc/...

# BTC Infrastructureテスト
go test ./internal/infrastructure/api/bitcoin/btc/...

# 統合テスト
make gotest-integration
```

## Related Documentation

- [BTC README](../../../../docs/crypto/btc/README.md) - 詳細なBTCドキュメント
- [Taproot Guide](../../../../docs/crypto/btc/TAPROOT_GUIDE.md) - Taproot実装
- [PSBT Guide](../../../../docs/crypto/btc/psbt_user_guide.md) - PSBT使用方法
- [MuSig2 Guide](../../../../docs/crypto/btc/musig2_guide.md) - MuSig2実装
- [Descriptor Guide](../../../../docs/descriptor_user_guide.md) - Descriptor使用方法

## Common Operations

| 操作 | Wallet | Use Case |
|------|--------|----------|
| アドレス生成 | Keygen | `keygen/btc/export_address.go` |
| トランザクション作成 | Watch | `watch/btc/create_transaction.go` |
| 署名 | Keygen/Sign | `{keygen,sign}/btc/sign_transaction.go` |
| 送信 | Watch | `watch/btc/send_transaction.go` |
| 残高確認 | Watch | Infrastructure API |
| UTXO取得 | Watch | Infrastructure API |

