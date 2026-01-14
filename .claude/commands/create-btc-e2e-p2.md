# Create BTC E2E Pattern 2 Script #{issue_number}

BTC E2Eテスト **Pattern 2: P2PKH 2-of-3 Multisig** の新規スクリプトを作成する。

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは、以下のE2Eスクリプトを新規作成します：

```
scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh
```

### Pattern 2 の技術仕様

| 項目 | 値 |
|------|-----|
| **パターン番号** | 2 |
| **鍵タイプ** | P2PKH (BIP44 Legacy) |
| **スクリプトタイプ** | 2-of-3 Multisig (P2SH wrapped) |
| **アドレス形式** | `3...` (mainnet), `2...` (testnet/regtest) |
| **署名要件** | 2-of-3 (任意の2つの署名で完了) |
| **Descriptor** | `sh(multi(2, [fp/44'/0'/0']xpub1/0/*, xpub2/0/*, xpub3/0/*))` |
| **必要なウォレット** | watch, keygen, sign1, sign2 |
| **環境変数** | `WALLET_ADDRESS_TYPE="legacy"` |

### Pattern 1 (Single-sig) との違い

| 項目 | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| 署名要件 | Single-sig (1つ) | 2-of-3 Multisig |
| 必要ウォレット | keygen のみ | keygen + sign1 + sign2 |
| アドレス形式 | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey交換 | 不要 | 必要 |
| account設定 | `account_singlesig.yaml` | `account_2of3.yaml` (新規作成) |

### Pattern 8 (3-of-3) との違い

| 項目 | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| 鍵タイプ | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| 署名要件 | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |
| 署名フロー | 2回署名で完了 | 3回署名必要 |

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `feat/issue-{issue_number}-btc-e2e-p2`
- **コミットタイプ**: `feat(btc)`
- **スコープ**: BTC E2E Pattern 2

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで作成のみ行います。

## 必須ドキュメント（最初に読み込む）

以下のドキュメントを最初に読み込んでください：

### 参照スクリプト（テンプレート）

- @scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh - **Pattern 1 スクリプト（Single-sig部分の参考）**
- @scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh - **Pattern 8 スクリプト（Multisig部分の参考）**
- @scripts/operation/common.sh - 共通ユーティリティ

### 設計ドキュメント

- @docs/crypto/btc/e2e_transaction_patterns.md - 全パターンの詳細
- @pkg/config/README.md - 環境変数による設定上書き

### 設定ファイル

- @config/wallet/account.yaml - 現在のマルチシグ設定（参考）
- @config/wallet/account_singlesig.yaml - シングルシグ設定（参考）

## 作成するファイル一覧

### 1. E2Eスクリプト（必須）

```
scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh
```

### 2. アカウント設定ファイル（必須）

```
config/wallet/account_2of3.yaml
```

2-of-3 マルチシグ用の設定：

```yaml
# Account Configuration for Pattern 2 (P2PKH 2-of-3 Multisig E2E Test)
types:
  - client
  - deposit
  - payment
  - stored

deposit_receiver: deposit
payment_sender: payment

# 2-of-3 multisig configuration
multisig:
  - type: deposit
    required: 2  # 2-of-3
    auth_users:
      - auth1
      - auth2

  - type: payment
    required: 2  # 2-of-3
    auth_users:
      - auth1
      - auth2

  - type: stored
    required: 2  # 2-of-3
    auth_users:
      - auth1
      - auth2
```

### 3. Makefileターゲット追加（必須）

`make/btc_e2e.mk` に追加：

```makefile
###############################################################################
# E2E Testing - Pattern 2: P2PKH 2-of-3 Multisig
###############################################################################
.PHONY: btc-e2e-p2pkh-2of3-reset
btc-e2e-p2pkh-2of3-reset:
 ./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --reset

.PHONY: btc-e2e-p2pkh-2of3
btc-e2e-p2pkh-2of3:
 ./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh

.PHONY: btc-e2e-p2pkh-2of3-verbose
btc-e2e-p2pkh-2of3-verbose:
 ./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --verbose

.PHONY: btc-e2e-p2pkh-2of3-ci
btc-e2e-p2pkh-2of3-ci:
 ./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --non-interactive

.PHONY: btc-e2e-p2pkh-2of3-cleanup
btc-e2e-p2pkh-2of3-cleanup:
 ./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --cleanup
```

## スクリプト実装手順

### Step 1: スクリプトの基本構造

Pattern 8 (`e2e-p2sh-p2wsh-3of3.sh`) をベースに、以下を変更：

1. **ヘッダーコメント更新**

```bash
# Bitcoin E2E Workflow Script - Pattern 2: P2PKH 2-of-3 Multisig
# Transaction Pattern:
#   Pattern 2: BTC P2PKH 2-of-3 Multisig
#   - Address Type: P2PKH (BIP44 Legacy) wrapped in P2SH
#   - Address Format: `3...` (Testnet: `2...`)
#   - Signature Requirement: 2-of-3 (任意の2つの署名)
#   - Descriptor: sh(multi(2, xpub1, xpub2, xpub3))
```

1. **環境変数設定**

```bash
# Pattern 2 (P2PKH 2-of-3 Multisig) requires:
#   - address_type: "legacy" (derives key_type: bip44 automatically)
export WALLET_ADDRESS_TYPE="legacy"
```

1. **アカウント設定ファイル**

```bash
# Use 2-of-3 multisig account configuration for Pattern 2
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account_2of3.yaml"
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"
```

### Step 2: キー生成フェーズ

Pattern 8 と同様に以下を実行：

1. Keygen で Seed 作成
2. Keygen で HD Key 作成（BIP44パスで生成）
3. Sign1/Sign2 で Seed 作成
4. Sign1/Sign2 で HD Key 作成
5. Sign1/Sign2 から fullpubkey エクスポート
6. Keygen に fullpubkey インポート

### Step 3: Descriptor エクスポート

**重要**: Pattern 2 のDescriptor形式

```
sh(multi(2, [fp/44'/0'/0']xpub1/0/*, [fp/44'/0'/0']xpub2/0/*, [fp/44'/0'/0']xpub3/0/*))
```

- `sh()` - P2SH wrapper
- `multi(2, ...)` - 2-of-3 multisig (sortedmulti ではなく multi を使用可能)
- 鍵派生パス: BIP44 (`m/44'/0'/0'`)

### Step 4: 署名フロー（2-of-3）

**Pattern 8 との最大の違い**: 2回の署名で完了

```bash
# Sign with keygen wallet (1st signature)
tx_signed1=$(keygen ... sign signature --file "${tx_unsigned}")

# Sign with sign1 wallet (2nd signature) - これで完了
tx_signed2=$(sign1 ... sign signature --file "${tx_signed1}")

# sign2 は不要 - 2つの署名で2-of-3が満たされる
# Send transaction
watch ... send --file "${tx_signed2}"
```

### Step 5: アドレス検証

Payment アドレスが `2...` (regtest P2SH) で始まることを確認：

```bash
# P2SH address for 2-of-3 multisig (starts with '2' in regtest)
sender_address=$(docker compose exec -T wallet-db mysql ... \
    "SELECT wallet_address FROM address WHERE coin='btc' AND account='payment' AND wallet_address LIKE '2%' LIMIT 1")
```

## 実装チェックリスト

### スクリプト作成

- [ ] `e2e-p2pkh-2of3.sh` を作成
- [ ] ヘッダーコメントを Pattern 2 用に更新
- [ ] `WALLET_ADDRESS_TYPE="legacy"` を設定
- [ ] `CONFIG_ACCOUNT` を `account_2of3.yaml` に変更
- [ ] 署名フローを2回に変更（3回目を削除）

### 設定ファイル

- [ ] `account_2of3.yaml` を作成（`required: 2`）
- [ ] `make/btc_e2e.mk` にターゲット追加

### 検証

- [ ] シェルスクリプトの lint: `make shfmt`
- [ ] 実行テスト: `make btc-e2e-p2pkh-2of3-reset`
- [ ] 生成されるアドレスが `2...` で始まることを確認
- [ ] 2回の署名でトランザクションが完了することを確認

## ドキュメント更新

E2Eスクリプト作成後、以下を更新：

### 1. e2e_transaction_patterns.md

Pattern 2 の対応状況を更新：

```markdown
| 2 | P2PKH (BIP44) | 2-of-3 Multisig | `3...` (P2SH wrapped) | ✅ e2e/e2e-p2pkh-2of3.sh |
```

### 2. E2Eスクリプトの対応表更新

```markdown
| `scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh` | BTC | P2PKH 2-of-3 Multisig | 2-of-3 |
```

## エラー対応

### よくあるエラー

#### 1. Descriptor形式エラー

**症状**: Descriptor export/import 時にエラー

**原因**: BIP49 形式になっている

**解決**: `address_type` が `legacy` であることを確認

```bash
echo $WALLET_ADDRESS_TYPE  # "legacy" であること
```

#### 2. アドレス形式が異なる

**症状**: `m...` や `n...` で始まるアドレスが生成される

**原因**: Single-sig アドレスが生成されている

**解決**:

- `account_2of3.yaml` の `multisig` セクションを確認
- fullpubkey のインポートが成功しているか確認

#### 3. 署名が足りない

**症状**: トランザクション送信時に「署名が不完全」エラー

**原因**: 2回目の署名が正しく適用されていない

**確認**:

```bash
# PSBTの署名状態を確認
btc_cli "btc-watch" analyzepsbt "${psbt_hex}"
```

## 試行回数の上限

**修正→テストのサイクルが5回を超えた場合、進捗を整理して報告してください。**

### エスカレーション条件

- 同じエラーが繰り返し発生
- Descriptor 形式の仕様理解が必要
- Go コード側の修正が必要になった場合

## 必須スキル

1. `shell-scripts` - スクリプト作成
2. `go-development` - Go コードの修正が必要な場合
3. `git-workflow` - ブランチ管理（issue番号指定時）
4. `makefile-update` - Makefile更新

## 関連ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| `docs/crypto/btc/e2e_transaction_patterns.md` | 全11パターンの説明 |
| `docs/crypto/btc/descriptor_examples.md` | Descriptor の例 |
| `scripts/operation/btc/README.md` | E2E スクリプトの詳細 |

## 完了条件

以下がすべて完了したら作業終了：

1. ✅ `e2e-p2pkh-2of3.sh` が作成されている
2. ✅ `account_2of3.yaml` が作成されている
3. ✅ `make/btc_e2e.mk` にターゲットが追加されている
4. ✅ `make btc-e2e-p2pkh-2of3-reset` が成功する
5. ✅ 生成されるアドレスが `2...` (P2SH) で始まる
6. ✅ 2回の署名でトランザクションがブロードキャストされる
7. ✅ `docs/crypto/btc/e2e_transaction_patterns.md` が更新されている
