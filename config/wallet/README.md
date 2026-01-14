# Wallet Configuration Files

このディレクトリには、各ウォレットタイプ・暗号通貨ごとの設定ファイルが格納されています。

## ディレクトリ構成

```
config/wallet/
├── btc_*.yaml         # Bitcoin wallet configurations
├── bch_*.yaml         # Bitcoin Cash wallet configurations
├── eth_*.yaml         # Ethereum wallet configurations
├── xrp_*.yaml         # XRP (Ripple) wallet configurations
├── account*.yaml      # Account type configurations
└── archive/           # Legacy TOML configurations (deprecated)
```

## 設定ファイルの種類

### Wallet Config (`{chain}_{wallet_type}.yaml`)

| ファイル | 用途 |
|----------|------|
| `btc_watch.yaml` | BTC Watch-only wallet (オンライン) |
| `btc_keygen.yaml` | BTC Keygen wallet (オフライン環境推奨) |
| `btc_sign.yaml` | BTC Sign wallet (オフライン環境推奨) |
| `btc_sign1.yaml` / `btc_sign2.yaml` | マルチシグ用署名者1/2 |

### Account Config (`account*.yaml`)

| ファイル | 用途 |
|----------|------|
| `account.yaml` | デフォルトのアカウント構成 |
| `account_singlesig.yaml` | シングルシグ用アカウント構成 |
| `account_2of3.yaml` | 2-of-3 マルチシグ用アカウント構成 |

## ⚠️ 重要: 設定ファイルの運用方針

**これらの設定ファイルはスクリプト実行ごとに編集・上書きすることを想定していません。**

### 環境変数による設定の上書き

異なる設定でCLIを実行する場合は、**環境変数を使用**してください。

```bash
# 設定ファイルを編集せずに、環境変数で上書き
export WALLET_ADDRESS_TYPE=legacy
export WALLET_BITCOIN_HOST=127.0.0.1:18443
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh
```

#### 環境変数の命名規則

```
WALLET_{KEY}              # トップレベル設定
WALLET_{SECTION}_{KEY}    # ネストされた設定
```

#### 主要な環境変数

| 設定ファイルのキー | 環境変数 |
|-------------------|----------|
| `address_type` | `WALLET_ADDRESS_TYPE` |
| `bitcoin.host` | `WALLET_BITCOIN_HOST` |
| `bitcoin.user` | `WALLET_BITCOIN_USER` |
| `bitcoin.pass` | `WALLET_BITCOIN_PASS` |
| `mysql.host` | `WALLET_MYSQL_HOST` |
| `mysql.dbname` | `WALLET_MYSQL_DBNAME` |
| `logger.level` | `WALLET_LOGGER_LEVEL` |

#### 優先順位

1. **環境変数** (最優先)
2. **設定ファイル**
3. **デフォルト値** (最低優先)

詳細は [pkg/config/README.md](../../pkg/config/README.md) を参照。

## Bitcoin RPC Host の設定

`bitcoin.host` の設定には、**ドメインのみ**と**パス付き**の2つの形式があります。

### 形式1: ドメインのみ

```yaml
bitcoin:
  host: "127.0.0.1:20332"
```

### 形式2: パス付き (`/wallet/<name>`)

```yaml
bitcoin:
  host: "127.0.0.1:18332/wallet/watch"
```

### どちらを使うべきか？

| ウォレットタイプ | 推奨形式 | 理由 |
|-----------------|----------|------|
| **Watch** | パス付き | Descriptorをインポートしてbitcoindでアドレス/UTXO管理を行う |
| **Keygen** | パス付き | Descriptorをインポートしてbitcoindでアドレス/UTXO管理を行う |
| **Sign** | ドメインのみ | Descriptorを使用しない。秘密鍵をアプリ内で管理 |

### パス付き形式が必要な場合

Bitcoindの [Descriptor Wallet](https://bitcoincore.org/en/doc/24.0.0/rpc/wallet/importdescriptors/) 機能を使用する場合、事前にbitcoind側でwalletを作成する必要があります。

```bash
# bitcoindでwalletを作成
bitcoin-cli createwallet "watch"
bitcoin-cli createwallet "keygen"

# CLIで特定のwalletを指定してRPCを実行
bitcoin-cli -rpcwallet=watch getwalletinfo
```

パス付き形式 (`host: "127.0.0.1:18332/wallet/watch"`) は、`bitcoin-cli -rpcwallet=watch` と同等の動作をします。

#### 作成スクリプト

Dockerコンテナ内でwalletを作成するスクリプト:

```bash
./scripts/operation/btc/create-bitcoind-wallet.sh
```

### 補足: Sign walletがパスを使わない理由

Sign walletは以下の理由でbitcoindのwallet機能を使用しません：

1. **オフライン環境**: Sign walletはセキュリティ上オフライン環境での運用を想定
2. **秘密鍵管理**: アプリケーション内で秘密鍵を管理し、bitcoindには渡さない
3. **トランザクション署名**: 未署名のPSBTを受け取り、アプリ内で署名処理を実行

## 関連ドキュメント

- [pkg/config/README.md](../../pkg/config/README.md) - 設定パッケージの詳細
- [docs/crypto/btc/descriptor_user_guide.md](../../docs/crypto/btc/descriptor_user_guide.md) - Descriptor Walletの解説
- [scripts/operation/btc/](../../scripts/operation/btc/) - BTC操作スクリプト
