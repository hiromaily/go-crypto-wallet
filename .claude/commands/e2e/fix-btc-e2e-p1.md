# Fix BTC E2E Pattern 1 Errors #{issue_number}

BTC E2Eテスト（パターン1: P2PKH Single-sig）の実行エラーを修正する。

## 前提条件

**以下の共通ルールを最初に読み込むこと：**

- @.claude/rules/btc/e2e-script.md - BTC E2E共通ルール（ビルド、検証、エスカレーション、セキュリティ）

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは `make btc-e2e-p1-reset` 実行時のエラーを診断・修正します。
スクリプトは既に存在するため、エラーの原因特定と修正に集中します。

### Pattern 1 の技術仕様

| 項目 | 値 |
|------|-----|
| **パターン番号** | 1 |
| **ネットワーク** | **regtest** (ローカル環境) |
| **鍵タイプ** | P2PKH (BIP44 Legacy) |
| **スクリプトタイプ** | Single-sig |
| **アドレス形式** | `m.../n...` (regtest/testnet P2PKH) |
| **署名要件** | Single-sig (1つの署名) |
| **Descriptor** | `pkh([fingerprint/44'/0'/0']xpub.../0/*)` |
| **必要なウォレット** | watch, keygen |
| **環境変数** | `WALLET_ADDRESS_TYPE="legacy"` |

### Pattern 2 (2-of-3 Multisig) との違い

| 項目 | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| 署名要件 | Single-sig (1つ) | 2-of-3 Multisig |
| 必要ウォレット | keygen のみ | keygen + sign1 + sign2 |
| アドレス形式 | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey交換 | 不要 | 必要 |
| account設定 | `account_singlesig.yaml` | `account_2of3.yaml` |

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p1`
- **コミットタイプ**: `fix(btc)`
- **スコープ**: BTC E2E Pattern 1

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで修正のみ行います。

## Pattern 1 固有のドキュメント

共通ルールの Required Documentation に加えて、以下を参照：

- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - **対象スクリプト**
- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Pattern 2 スクリプト（Multisig部分の参考）

## 事前確認: 環境変数

**パターン1では `WALLET_ADDRESS_TYPE="legacy"` が必要です。**

スクリプト内で自動設定されますが、確認用：

```bash
echo $WALLET_ADDRESS_TYPE  # "legacy" であること
```

> **Note**: 設定ファイルを直接編集しないでください。環境変数で上書きします。
> 詳細は共通ルールの「Configuration File Policy」を参照。

## エラー診断手順

### Step 1: エラー再現

```bash
# 完全リセットしてE2Eテストを実行
make btc-e2e-p1-reset
```

エラーメッセージを確認し、以下のカテゴリに分類。

### Step 2: エラーカテゴリの特定

| エラーメッセージ | カテゴリ | 参照セクション |
|----------------|---------|--------------|
| `No utxo` | UTXO関連 | [UTXO関連エラー](#utxo関連エラー) |
| `connection refused` | インフラ | 共通ルール参照 |
| `wallet not found` | ウォレット | [ウォレット関連エラー](#ウォレット関連エラー) |
| `signing failed` | 署名 | [署名関連エラー](#署名関連エラー) |
| `descriptor` | Descriptor | [Descriptor関連エラー](#descriptor関連エラー) |
| `address_type` mismatch | 設定 | [設定関連エラー](#設定関連エラー) |
| `duplicate key` | DB | 共通ルール参照 |

## Pattern 1 固有のエラーと解決策

共通エラー（connection refused, duplicate key等）は共通ルールを参照。

### UTXO関連エラー

#### "No utxo" エラー

**症状**:

```
Transaction creation failed
This could indicate:
  - No payment requests in database
  - No UTXOs available for payment account
  - UTXOs not mature enough (need 100+ confirmations)
```

**原因と解決策**:

1. **Descriptor が正しくインポートされていない**

   ```bash
   # デバッグ: アドレス情報確認
   docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
     getaddressinfo "<payment_address>"
   ```

   確認ポイント:
   - `solvable: true` であること
   - `ismine: true` (watch-only wallet では false でも OK)

2. **ブロック生成が不足**

   ```bash
   # ブロック数確認
   docker exec btc-watch bitcoin-cli -regtest getblockcount
   # 101 以上であること
   ```

3. **address_type の不一致**

   ```bash
   # 環境変数確認
   echo $WALLET_ADDRESS_TYPE  # "legacy" であること
   ```

### ウォレット関連エラー

#### "wallet not found" エラー

**解決策**:

```bash
# ウォレット一覧確認
docker exec btc-watch bitcoin-cli -regtest listwallets

# ウォレット作成
docker exec btc-watch bitcoin-cli -regtest createwallet "watch" true true
docker exec btc-keygen bitcoin-cli -regtest createwallet "keygen" false true
```

### 署名関連エラー

#### 署名失敗エラー

**原因**: 秘密鍵がインポートされていない、または address_type の不一致

**確認手順**:

```bash
# 1. 環境変数確認
echo "WALLET_ADDRESS_TYPE: $WALLET_ADDRESS_TYPE"  # "legacy" であること

# 2. Keygen ウォレットの秘密鍵確認
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true
```

### Descriptor関連エラー

#### Descriptor インポート失敗

**確認手順**:

```bash
# Descriptor ファイル確認
cat data/descriptor/btc/payment_descriptors.json

# フォーマット確認 (P2PKH なら pkh(...) 形式)
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
```

**期待される形式** (Pattern 1):

```
pkh([fingerprint/44'/0'/0']xpub.../0/*)
```

### 設定関連エラー

#### address_type 不一致エラー

**症状**: 生成されるアドレスが期待と異なる

**確認手順**:

```bash
# スクリプト内の環境変数エクスポートを確認
grep -A2 "Environment Variable Overrides" \
  scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

**修正**: スクリプト内で以下が設定されていることを確認

```bash
export WALLET_ADDRESS_TYPE="legacy"
```

## デバッグ用コマンド

### 状態確認

```bash
# Bitcoin ノード状態
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# ウォレット残高
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO 一覧
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# DB のアドレス確認
docker compose exec -T wallet-db mysql -u root -proot watch -e \
  "SELECT wallet_address, account FROM address WHERE coin='btc' LIMIT 10"

# Payment request 確認
docker compose exec -T wallet-db mysql -u root -proot watch -e \
  "SELECT * FROM payment_request WHERE coin='btc'"
```

### ログ確認

```bash
# 詳細モードで実行
make btc-e2e-p1-verbose

# または
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh --verbose --reset
```

## 修正ファイルの特定

| エラー種別 | 修正対象ファイル |
|-----------|-----------------|
| スクリプトロジック | `scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh` |
| 共通関数 | `scripts/operation/common.sh` |
| Descriptor 生成 | `internal/application/usecase/keygen/btc/` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Bitcoin RPC | `internal/infrastructure/wallet/api/btc/` |
| 設定読み込み | `pkg/config/` |

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch wallet ユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC 実装 |
| `internal/domain/wallet/key/` | 鍵ドメインモデル |
| `pkg/config/loader.go` | 設定ローダー |

## 注意事項

### 既存スクリプトへの影響を避ける

- Pattern 2 (`e2e-p2-p2pkh-2of3.sh`) の動作を壊さないこと
- Pattern 8 (`e2e-p8-p2sh-p2wsh-3of3.sh`) の動作を壊さないこと
- `common.sh` を修正する場合は、他パターンへの影響を確認
- 環境変数の設定は各スクリプト内でローカルに行う

> **Note**: ビルドルール、検証コマンド、セキュリティは共通ルールを参照。

## クリーンアップ

```bash
# コンテナ停止のみ
make btc-e2e-p1-cleanup

# 完全リセット（データ含む）
make btc-e2e-p1-reset
```
