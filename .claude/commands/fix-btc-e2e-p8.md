# Fix BTC E2E Pattern 8 Test #{issue_number}

BTC E2Eテスト（パターン8: P2SH-P2WSH 3-of-3 Multisig）を実行し、エラーを修正する。

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは`scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh`を実行し、発生したエラーを分析・修正します。
Bitcoinの深い仕様理解が必要なため、以下の技術情報を参照してください。

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p8`
- **コミットタイプ**: `fix(btc)`
- **スコープ**: BTC E2E Pattern 8

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで修正のみ行います。

## 必須ドキュメント（最初に読み込む）

以下のドキュメントを最初に読み込んでください：

- @docs/crypto/btc/e2e_transaction_patterns.md - 全8パターンの詳細（**必須**）
- @scripts/operation/btc/README.md - E2Eスクリプトの使い方
- @scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh - 実行対象のスクリプト

## 事前確認: Config設定

**パターン8を実行するには、すべてのウォレット設定で `address_type: "p2sh-segwit"` が必要です。**

以下のファイルを確認してください：

- @config/wallet/btc_watch.yaml - `address_type: "p2sh-segwit"` であること
- @config/wallet/btc_keygen.yaml - `address_type: "p2sh-segwit"` であること
- @config/wallet/btc_sign1.yaml - `address_type: "p2sh-segwit"` であること
- @config/wallet/btc_sign2.yaml - `address_type: "p2sh-segwit"` であること

**重要**: 4つのウォレット設定すべてで同じ `address_type` を使用する必要があります。不一致があると、アドレス生成や署名で問題が発生します。

```yaml
# 正しい設定例（パターン8用）
address_type: "p2sh-segwit"
```

他のパターンを使用する場合は、`address_type` を適切に変更してください：

| パターン | address_type |
|----------|--------------|
| 1, 2 (P2PKH) | `"legacy"` |
| 3, 4 (P2SH-P2WPKH) | `"p2sh-segwit"` |
| 5, 6, 7 (P2WPKH/P2WSH) | `"bech32"` |
| **8 (P2SH-P2WSH)** | **`"p2sh-segwit"`** |
| 9, 10, 11 (P2TR) | `"bech32m"` |

## 試行回数の上限とエスカレーション

**修正→確認のサイクルが5回を超えた場合、GitHub issueを作成して処理を終了してください。**

これはコンテキストのトークン消費を抑え、複雑な問題を適切にエスカレーションするためです。

### エスカレーション条件

以下のいずれかに該当する場合、5回未満でもissue作成を検討：

- 同じエラーが繰り返し発生する
- エラーの原因が特定できない
- Bitcoin仕様の深い理解が必要と判断される

### GitHub Issue作成手順

```bash
gh issue create --title "BTC E2E Pattern 8: [エラー概要]" --body "$(cat <<'EOF'
## 概要

BTC E2Eテスト（パターン8: P2SH-P2WSH 3-of-3）で発生したエラー

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
make btc-e2e-p2sh-3of3-reset
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
make btc-e2e-p2sh-3of3-reset

# または既存状態から実行
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh

# デバッグ出力付き
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --verbose
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
make btc-e2e-p2sh-3of3-reset
```

**ワンライナー（Goコード修正後の完全検証）:**

```bash
make go-lint && make check-build && make gotest && make build-all && make btc-e2e-p2sh-3of3-reset
```

## 技術仕様: パターン8 (P2SH-P2WSH 3-of-3)

### アドレス形式

```
Address Type: P2SH-P2WSH (BIP49 wrapped SegWit)
Address Format: 3... (Mainnet), 2... (Testnet/Regtest)
Descriptor: sh(wsh(sortedmulti(3, xpub1, xpub2, xpub3)))
```

### 署名フロー

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature)
    ↓
Sign2 Wallet (3rd signature)
    ↓
Watch Wallet (broadcast)
```

### 重要な技術ポイント

1. **Descriptor Solvability**
   - P2SH-P2WSHアドレスはDescriptorインポートが必要
   - 従来のアドレスインポートではUTXOが「unsolvable」になる
   - `importdescriptors` RPCで解決

2. **3-of-3 Multisig**
   - 3つのxpubから派生したアドレスを使用
   - 全ての署名者の秘密鍵が必要
   - PSBTフォーマットで部分署名を伝達

3. **HD Key派生パス**
   - BIP49: `m/49'/0'/account'/change/index`
   - 各ウォレットで同じインデックスの鍵を使用

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

### "Unsolvable" UTXOエラー

**原因**: P2SH-P2WSHアドレスがDescriptorなしでインポートされた

**解決策**:

1. `descriptor export`で正しい形式をエクスポート
2. `descriptor import`で`--account`を指定してインポート
3. 関連コード: `internal/infrastructure/wallet/api/btc/descriptor.go`

### 署名エラー

**原因**: 秘密鍵がインポートされていない、または鍵の不一致

**解決策**:

1. 各ウォレットで`import privkey`が実行されているか確認
2. fullpubkeyのインポート順序を確認
3. 関連コード: `internal/infrastructure/wallet/api/btc/sign.go`

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
| `docs/crypto/btc/e2e_transaction_patterns.md` | 全8パターンの説明 |
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
| `internal/domain/wallet/key/` | 鍵ドメインモデル |

## 注意事項

### ビルドルール

**Goコード修正後は必ずバイナリを再ビルドしてください。**

- E2Eテストは `${GOPATH}/bin/` のバイナリを使用
- `go build` を直接実行しない（`make build-all` を使用）
- 詳細: `go-development` スキルの「Build Rules」セクション

### 他パターンへの影響を避ける

- パターン8固有の修正は`P2SH-P2WSH`関連コードに限定
- 共通コードを修正する場合は、他パターン（特に1-7, 9-11）への影響を確認
- 共通関数を修正する場合は単体テストで回帰を確認

### セキュリティ

- 秘密鍵をログに出力しない
- テスト用のパスフレーズ/RPCクレデンシャルは本番で使用しない
- `docs/standards/security.md`を参照

## クリーンアップ

```bash
# コンテナ停止のみ
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --cleanup

# 完全リセット（データ含む）
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --reset
```
