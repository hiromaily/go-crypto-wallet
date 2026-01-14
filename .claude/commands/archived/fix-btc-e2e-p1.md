# Fix BTC E2E Pattern 1 Errors #{issue_number}

BTC E2Eテスト（パターン1: P2PKH Single-sig）の実行エラーを修正する。

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは `make btc-e2e-p2pkh-reset` 実行時のエラーを診断・修正します。
スクリプトは既に存在するため、エラーの原因特定と修正に集中します。

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p1`
- **コミットタイプ**: `fix(btc)`
- **スコープ**: BTC E2E Pattern 1

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで修正のみ行います。

## 必須ドキュメント（最初に読み込む）

以下のドキュメントを最初に読み込んでください：

- @scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh - **対象スクリプト**
- @scripts/operation/common.sh - 共通ユーティリティ
- @pkg/config/README.md - **環境変数による設定上書き**
- @docs/crypto/btc/e2e_transaction_patterns.md - 全パターンの詳細

## 環境変数による設定上書き（重要）

### 仕組み

Config ファイルの値は環境変数で上書きできます：

```
Priority: Environment Variables > Config File > Default Values
```

### key_type の自動導出

`key_type` は `address_type` から自動的に導出されます（`internal/domain/address/types.go`）：

| address_type | 自動導出される key_type | 用途 |
|--------------|----------------------|------|
| `legacy` | `bip44` | P2PKH (Pattern 1) |
| `p2sh-segwit` | `bip49` | P2SH-P2WPKH/P2SH-P2WSH (Pattern 8) |
| `bech32` | `bip84` | Native SegWit |
| `taproot` | `bip86` | Taproot |

**注意**: `WALLET_KEY_TYPE` 環境変数は不要です。`WALLET_ADDRESS_TYPE` のみ設定してください。

### Pattern 1 (P2PKH Single-sig) に必要な設定

| 設定 | 環境変数名 | Pattern 1 の値 | Config ファイルのデフォルト |
|------|-----------|---------------|---------------------------|
| `address_type` | `WALLET_ADDRESS_TYPE` | `legacy` | `p2sh-segwit` |

### スクリプト内での設定（確認必須）

`e2e-p2pkh-singlesig.sh` で環境変数をエクスポートしています：

```bash
# Pattern 1 (P2PKH Single-sig) requires:
#   - address_type: "legacy" (derives key_type: bip44 automatically)
export WALLET_ADDRESS_TYPE="legacy"
```

**エラー発生時は、この設定が正しく適用されているか確認してください。**

### 設定の適用先

| Wallet | Config ファイル | 必要な address_type |
|--------|----------------|-------------------|
| Watch | `config/wallet/btc_watch.yaml` | `legacy` (環境変数で上書き) |
| Keygen | `config/wallet/btc_keygen.yaml` | `legacy` (環境変数で上書き) |

**注意**: Config ファイル自体は変更不要です。環境変数による上書きで対応します。

## エラー診断手順

### Step 1: エラー再現

```bash
# 完全リセットしてE2Eテストを実行
make btc-e2e-p2pkh-reset
```

エラーメッセージを確認し、以下のカテゴリに分類：

### Step 2: エラーカテゴリの特定

| エラーメッセージ | カテゴリ | 参照セクション |
|----------------|---------|--------------|
| `No utxo` | UTXO関連 | [UTXO関連エラー](#utxo関連エラー) |
| `connection refused` | インフラ | [インフラエラー](#インフラエラー) |
| `wallet not found` | ウォレット | [ウォレット関連エラー](#ウォレット関連エラー) |
| `signing failed` | 署名 | [署名関連エラー](#署名関連エラー) |
| `descriptor` | Descriptor | [Descriptor関連エラー](#descriptor関連エラー) |
| `address_type` mismatch | 設定 | [設定関連エラー](#設定関連エラー) |
| `duplicate key` | DB | [DB関連エラー](#db関連エラー) |

## よくあるエラーと解決策

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

### インフラエラー

#### "connection refused" エラー

**症状**:

```
Error: failed to connect to Bitcoin node
```

**解決策**:

```bash
# コンテナ状態確認
docker compose -f compose.btc.yaml ps

# ログ確認
docker compose -f compose.btc.yaml logs btc-watch
docker compose -f compose.btc.yaml logs btc-keygen

# 再起動
docker compose -f compose.btc.yaml down
docker compose -f compose.btc.yaml up -d btc-watch btc-keygen
```

#### Database 接続エラー

```bash
# DB コンテナ状態確認
docker compose -f compose.yaml ps

# DB 再起動
docker compose -f compose.yaml down
docker compose -f compose.yaml up -d
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
  scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh
```

**修正**: スクリプト内で以下が設定されていることを確認

```bash
export WALLET_ADDRESS_TYPE="legacy"
```

### DB関連エラー

#### "duplicate key" エラー

**原因**: 前回の実行データが残っている

**解決策**:

```bash
# 完全リセット（推奨）
make btc-e2e-p2pkh-reset

# または手動でボリューム削除
docker compose -f compose.yaml down -v
docker volume rm go-crypto-wallet_wallet-db
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
make btc-e2e-p2pkh-verbose

# または
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --verbose --reset
```

## 修正ファイルの特定

| エラー種別 | 修正対象ファイル |
|-----------|-----------------|
| スクリプトロジック | `scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh` |
| 共通関数 | `scripts/operation/common.sh` |
| Descriptor 生成 | `internal/application/usecase/keygen/btc/` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Bitcoin RPC | `internal/infrastructure/wallet/api/btc/` |
| 設定読み込み | `pkg/config/` |

## 修正後の検証

### 1. シェルスクリプトの lint

```bash
make shfmt
```

### 2. Go コードを修正した場合

```bash
make go-lint
make check-build
make gotest
```

### 3. E2E テスト実行

```bash
# 完全リセットからの実行
make btc-e2e-p2pkh-reset
```

## 試行回数の上限とエスカレーション

**修正→テストのサイクルが5回を超えた場合、進捗を整理して報告してください。**

### エスカレーション条件

以下のいずれかに該当する場合、5回未満でも報告を検討：

- 同じエラーが繰り返し発生する
- Bitcoin 仕様の深い理解が必要と判断される
- Go コード側の大規模な修正が必要になった場合
- 環境変数の上書きが機能していない場合

### 進捗報告フォーマット

```markdown
## 進捗報告

### エラー内容
[発生しているエラーメッセージ]

### 試行した修正
1. [修正内容1]
2. [修正内容2]
...

### 現在の状態
[現在の状態を説明]

### 環境変数の設定状況
- WALLET_ADDRESS_TYPE: [値] (key_type は自動導出)

### 次のステップ
[次に必要なアクション]
```

## 必須スキル

1. `shell-scripts` - シェルスクリプトの修正
2. `go-development` - Go コードの修正が必要な場合
3. `git-workflow` - ブランチ管理とコミット（issue番号指定時）

## 関連ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| `docs/crypto/btc/e2e_transaction_patterns.md` | 全11パターンの説明 |
| `docs/crypto/btc/descriptor_examples.md` | Descriptor の例 |
| `docs/crypto/btc/psbt_developer_guide.md` | PSBT 開発ガイド |
| `scripts/operation/btc/README.md` | E2E スクリプトの詳細 |
| `pkg/config/README.md` | 環境変数による設定上書き |

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch wallet ユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC 実装 |
| `internal/domain/wallet/key/` | 鍵ドメインモデル |
| `pkg/config/loader.go` | 設定ローダー |

## 注意事項

### ビルドルール

**Go コード修正後は必ずバイナリを再ビルドしてください。**

```bash
make build-all
```

### Config ファイルは変更しない

Config ファイル (`btc_watch.yaml`, `btc_keygen.yaml`) は変更せず、
**環境変数 (`WALLET_ADDRESS_TYPE`) による上書き**で対応してください。
`key_type` は `address_type` から自動導出されます。

### 既存スクリプトへの影響を避ける

- Pattern 8 (`e2e-p2sh-p2wsh-3of3.sh`) の動作を壊さないこと
- `common.sh` を修正する場合は、他パターンへの影響を確認
- 環境変数の設定は各スクリプト内でローカルに行う

### セキュリティ

- 秘密鍵をログに出力しない
- テスト用のパスフレーズ/RPC クレデンシャルは本番で使用しない
- `docs/standards/security.md` を参照

## クリーンアップ

```bash
# コンテナ停止のみ
make btc-e2e-p2pkh-cleanup

# 完全リセット（データ含む）
make btc-e2e-p2pkh-reset
```
