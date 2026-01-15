# Fix BTC E2E Pattern 2 Test #{issue_number}

BTC E2Eテスト（パターン2: P2PKH 2-of-3 Multisig）を **regtest環境** で実行し、エラーを修正する。

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは`scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh`を **Bitcoin Core regtest環境** で実行し、発生したエラーを分析・修正します。

> **Note**: このE2Eテストはローカルのregtest（Regression Test）環境で実行されます。
> 実際のBitcoinネットワーク（mainnet/testnet）には接続しません。

### Pattern 2 の技術仕様

| 項目 | 値 |
|------|-----|
| **パターン番号** | 2 |
| **ネットワーク** | **regtest** (ローカル環境) |
| **鍵タイプ** | P2PKH (BIP44 Legacy) |
| **スクリプトタイプ** | 2-of-3 Multisig (P2SH wrapped) |
| **アドレス形式** | `2...` (regtest/testnet P2SH) |
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
| account設定 | `account_singlesig.yaml` | `account_2of3.yaml` |

### Pattern 8 (3-of-3) との違い

| 項目 | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| 鍵タイプ | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| 署名要件 | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |
| 署名フロー | 2回署名で完了 | 3回署名必要 |

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p2`
- **コミットタイプ**: `fix(btc)`
- **スコープ**: BTC E2E Pattern 2

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで修正のみ行います。

## 必須ドキュメント（最初に読み込む）

以下のドキュメントを最初に読み込んでください：

- @docs/crypto/btc/e2e_transaction_patterns.md - 全パターンの詳細（**必須**）
- @scripts/operation/btc/README.md - E2Eスクリプトの使い方
- @scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh - 実行対象のスクリプト

### 参照スクリプト

- @scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh - Pattern 1 スクリプト（Single-sig部分の参考）
- @scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh - Pattern 8 スクリプト（Multisig部分の参考）
- @scripts/operation/common.sh - 共通ユーティリティ

### 設定ファイル

- @config/wallet/account_2of3.yaml - 2-of-3 マルチシグ設定
- @config/wallet/account.yaml - 現在のマルチシグ設定（参考）

## 事前確認: Config設定

**パターン2を実行するには、すべてのウォレット設定で `address_type: "legacy"` が必要です。**

以下のファイルを確認してください：

- @config/wallet/btc_watch.yaml - `address_type: "legacy"` であること
- @config/wallet/btc_keygen.yaml - `address_type: "legacy"` であること
- @config/wallet/btc_sign1.yaml - `address_type: "legacy"` であること
- @config/wallet/btc_sign2.yaml - `address_type: "legacy"` であること

**重要**: 4つのウォレット設定すべてで同じ `address_type` を使用する必要があります。不一致があると、アドレス生成や署名で問題が発生します。

```yaml
# 正しい設定例（パターン2用）
address_type: "legacy"
```

他のパターンを使用する場合は、`address_type` を適切に変更してください：

| パターン | address_type |
|----------|--------------|
| **1, 2 (P2PKH)** | **`"legacy"`** |
| 3, 4 (P2SH-P2WPKH) | `"p2sh-segwit"` |
| 5, 6, 7 (P2WPKH/P2WSH) | `"bech32"` |
| 8 (P2SH-P2WSH) | `"p2sh-segwit"` |
| 9, 10, 11 (P2TR) | `"bech32m"` |

## 試行回数の上限とエスカレーション

**修正→確認のサイクルが5回を超えた場合、GitHub issueを作成して処理を終了してください。**

これはコンテキストのトークン消費を抑え、複雑な問題を適切にエスカレーションするためです。

### エスカレーション条件

以下のいずれかに該当する場合、5回未満でもissue作成を検討：

- 同じエラーが繰り返し発生する
- エラーの原因が特定できない
- Bitcoin仕様の深い理解が必要と判断される
- Goコード側の修正が必要になった場合

### GitHub Issue作成手順

```bash
gh issue create --title "BTC E2E Pattern 2: [エラー概要]" --body "$(cat <<'EOF'
## 概要

BTC E2Eテスト（パターン2: P2PKH 2-of-3）で発生したエラー

## 試行履歴

| 回 | 修正内容 | 結果 |
|----|----------|------|
| 1 | [修正内容] | [エラー内容] |
| 2 | [修正内容] | [エラー内容] |
| ... | ... | ... |

## 最終エラー

```

[エラーメッセージ全文]

```

## 分析

- **エラー発生フェーズ**: [Prerequisites/Infrastructure/Key Generation/Multisig Setup/UTXO Generation/Transaction Flow]
- **推定原因**: [推定される原因]
- **試した解決策**: [試した内容のリスト]

## 関連コード

- [修正したファイルのパス]

## 再現手順

```bash
make btc-e2e-p2pkh-2of3-reset
```

## 参考資料

- docs/crypto/btc/e2e_transaction_patterns.md
- scripts/operation/btc/README.md

EOF
)" --label "bug,chain:btc,lang:go"

```

### issue作成後

1. issue番号を報告
2. 作成したブランチがあればコミット（未完成でも）
3. 処理を終了

## 必須スキル

1. `git-workflow` - ブランチ管理とコミット
2. `go-development` - Goコードの修正時
3. `shell-scripts` - シェルスクリプトの修正時

## 実行手順

### Step 0: バイナリをビルド（必須）

**E2Eテストを実行する前に、必ずバイナリをビルドしてください。**

E2Eスクリプトは `watch`, `keygen`, `sign1`, `sign2` バイナリを `${GOPATH}/bin/` から呼び出します。

```bash
# 全ウォレットバイナリをビルド（推奨）
make build-all

# 出力先: ${GOPATH}/bin/watch, ${GOPATH}/bin/keygen, ${GOPATH}/bin/sign1, ${GOPATH}/bin/sign2
```

**重要: `go build` を直接実行しないでください。** 詳細は `go-development` スキルの「Build Rules」を参照。

### Step 1: E2Eテストを実行

```bash
# フルリセットして実行（推奨）
make btc-e2e-p2pkh-2of3-reset

# または既存状態から実行
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh

# デバッグ出力付き
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --verbose
```

### Step 2: エラー分析

エラーが発生したフェーズを特定し、対応するコードを調査：

| Phase | 関連コード | 説明 |
|-------|-----------|------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Multisig Setup | Descriptor export/import | `internal/application/usecase/watch/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

### Step 3: コードを修正

エラー種別に応じて適切なスキルをロード：

- **Goコード修正**: `go-development` SKILLを読み込む
- **シェルスクリプト修正**: `shell-scripts` SKILLを読み込む
- **設定ファイル修正**: 関連ドキュメントを参照

### Step 4: 検証

**Goコードを修正した場合、再ビルドしてからE2Eテストを実行してください。**

```bash
# 1. Goコード修正時の検証
make go-lint && make check-build && make gotest

# 2. 修正後、バイナリを再ビルド（Goコード修正時は必須）
make build-all

# 3. E2Eテストを再実行
make btc-e2e-p2pkh-2of3-reset
```

**ワンライナー（Goコード修正後の完全検証）:**

```bash
make go-lint && make check-build && make gotest && make build-all && make btc-e2e-p2pkh-2of3-reset
```

## 技術仕様: パターン2 (P2PKH 2-of-3)

### ネットワーク

このE2Eテストは **regtest（Regression Test）環境** で実行されます。

- **regtest**: ローカルでブロックを即座に生成可能なテスト環境
- mainnet/testnet には接続しない
- Docker コンテナ内で Bitcoin Core が regtest モードで起動

### アドレス形式

```
Network: regtest (Bitcoin Core local test environment)
Address Type: P2SH-wrapped P2PKH Multisig (BIP44 Legacy)
Address Format: 2... (Regtest/Testnet P2SH prefix)
Descriptor: sh(multi(2, xpub1, xpub2, xpub3))
```

> **Note**: Mainnet では `3...` で始まりますが、regtest/testnet では `2...` で始まります。

### 署名フロー（2-of-3）

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature) ← ここで完了
    ↓
Watch Wallet (broadcast)

※ Sign2 は不要 - 2つの署名で2-of-3が満たされる
```

### 重要な技術ポイント

1. **2-of-3 vs 3-of-3**
   - Pattern 2: 任意の2つの署名で完了
   - Pattern 8: 全3つの署名が必要
   - 署名フローの制御が異なる

2. **Descriptor形式**
   - `sh(multi(2, ...))` - P2SH wrapper + 2-of-3 multisig
   - `sortedmulti` ではなく `multi` を使用可能
   - 鍵の順序に注意

3. **HD Key派生パス**
   - BIP44: `m/44'/0'/account'/change/index`
   - Pattern 8 (BIP49) とは異なる派生パス

## よくあるエラーと解決策

### "No utxo" エラー

**原因**: UTXOが見つからない、またはアドレスがsolvableでない

**解決策**:

1. Descriptorが正しくインポートされているか確認
2. `getaddressinfo`でsolvable/ismine状態を確認
3. ブロック生成(101ブロック)が完了しているか確認

```bash
# デバッグ: アドレス情報確認
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
```

### Descriptor形式エラー

**症状**: Descriptor export/import 時にエラー

**原因**: BIP49 形式になっている（P2SH-SegWit）

**解決策**: `address_type` が `legacy` であることを確認

```bash
echo $WALLET_ADDRESS_TYPE  # "legacy" であること
```

### アドレス形式が異なる

**症状**: `m...` や `n...` で始まるアドレスが生成される（Single-sig P2PKHアドレス）

**原因**: Single-sig アドレスが生成されている

**解決策**:

- `account_2of3.yaml` の `multisig` セクションを確認
- fullpubkey のインポートが成功しているか確認
- `required: 2` が設定されているか確認

### 署名が足りない / 署名が多すぎる

**症状**: トランザクション送信時に「署名が不完全」または「署名が多すぎる」エラー

**原因**:

- 2回目の署名が正しく適用されていない
- 3回目の署名を行っている（不要）

**確認**:

```bash
# PSBTの署名状態を確認
btc_cli "btc-watch" analyzepsbt "${psbt_hex}"
```

### fullpubkey インポートエラー

**症状**: Multisig セットアップ時にエラー

**原因**: fullpubkey の形式不一致または順序の問題

**解決策**:

1. `sign1`, `sign2` から正しくエクスポートされているか確認
2. `keygen` へのインポート順序を確認
3. 関連コード: `internal/infrastructure/wallet/key/fullpubkey/`

### RPC接続エラー

**原因**: Dockerコンテナが起動していない

**解決策**:

```bash
# コンテナ状態確認
docker compose -f compose.btc.yaml ps

# ログ確認
docker compose -f compose.btc.yaml logs btc-watch
```

## 関連ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| `docs/crypto/btc/e2e_transaction_patterns.md` | 全パターンの説明 |
| `docs/crypto/btc/descriptor_examples.md` | Descriptorの例 |
| `docs/crypto/btc/psbt_developer_guide.md` | PSBT開発ガイド |
| `docs/crypto/btc/README.md` | BTC技術リファレンス |
| `scripts/operation/btc/README.md` | E2Eスクリプトの詳細 |

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch walletユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC実装 |
| `internal/infrastructure/wallet/key/fullpubkey/` | fullpubkey処理 |
| `internal/domain/wallet/key/` | 鍵ドメインモデル |

## 注意事項

### ビルドルール

**Goコード修正後は必ずバイナリを再ビルドしてください。**

- E2Eテストは `${GOPATH}/bin/` のバイナリを使用
- `go build` を直接実行しない（`make build-all` を使用）
- 詳細: `go-development` スキルの「Build Rules」セクション

### 他パターンへの影響を避ける

- パターン2固有の修正は`P2PKH 2-of-3`関連コードに限定
- 共通コードを修正する場合は、他パターン（特に1, 8）への影響を確認
- 共通関数を修正する場合は単体テストで回帰を確認

### セキュリティ

- 秘密鍵をログに出力しない
- テスト用のパスフレーズ/RPCクレデンシャルは本番で使用しない
- `docs/standards/security.md`を参照

## クリーンアップ

```bash
# コンテナ停止のみ
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --cleanup

# 完全リセット（データ含む）
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --reset
```
