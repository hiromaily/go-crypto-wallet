---
task_type: verification
description: ファイルタイプ別の検証コマンドマトリックス
version: 1.0.0
---

# Verification Matrix

タスク完了時に実行すべき検証コマンドは、編集したファイルの種類によって異なります。
このドキュメントは、ファイルタイプ別の検証コマンドを定義します。

## Quick Reference

| File Type | Required Commands | Optional Commands |
|-----------|-------------------|-------------------|
| `*.go` | `make go-lint`, `make check-build` | `make gotest` |
| `*.md` | (none) | `markdownlint` |
| `*.sql`, `*.hcl` | `make atlas-fmt`, `make atlas-lint` | |
| `*.yaml`, `*.toml` | (none) | YAML/TOML validation |
| `*.sh` | (none) | `shellcheck` |
| `*.proto` | `make proto` | |

## Detailed Verification by File Type

### Go Files (`*.go`)

**必須コマンド:**

```bash
make go-lint      # リンターチェック（自動修正含む）
make check-build  # ビルド確認
```

**推奨コマンド:**

```bash
make gotest       # テスト実行（機能変更時）
make tidy         # 依存関係整理（import変更時）
```

**追加コマンド（条件付き）:**

| Condition | Command | Description |
|-----------|---------|-------------|
| インターフェース変更 | `make mockery` | モック再生成 |
| SQLC関連変更 | `make sqlc` | SQLC再生成 |

### Markdown Files (`*.md`, `*.mdc`)

**必須コマンド:**

```bash
# なし - ドキュメントのみの変更では Go 関連コマンドは不要
```

**オプションコマンド:**

```bash
# markdownlint がインストールされている場合
markdownlint docs/**/*.md

# リンク確認
markdown-link-check docs/**/*.md
```

**注意**: Go ファイル内のコメント編集は Go ファイルの編集扱いです。

### Database Files (`*.sql`, `*.hcl`)

**必須コマンド:**

```bash
make atlas-fmt    # スキーマフォーマット
make atlas-lint   # スキーマリント
```

**追加ワークフロー:**

```bash
# スキーマ変更後
make atlas-dev-reset  # マイグレーション再生成
make sqlc             # SQLC再生成
make check-build      # ビルド確認
```

詳細は [db-change.md](./db-change.md) を参照。

### Configuration Files (`*.yaml`, `*.toml`)

**必須コマンド:**

```bash
# なし - 設定ファイルのみの変更では検証コマンド不要
```

**推奨確認:**

- YAML/TOML 構文の有効性
- 設定値の妥当性（手動確認）

### Shell Scripts (`*.sh`)

**必須コマンド:**

```bash
# なし
```

**オプションコマンド:**

```bash
# shellcheck がインストールされている場合
shellcheck scripts/**/*.sh
```

### Protocol Buffers (`*.proto`)

**必須コマンド:**

```bash
make proto        # protoc 再生成
make check-build  # ビルド確認
```

## Task Type × File Type Matrix

タスクタイプとファイルタイプの組み合わせによる検証コマンド:

| Task Type | Go Files | MD Files | SQL/HCL | Config |
|-----------|----------|----------|---------|--------|
| bug-fix | lint, build, test | (none) | atlas-fmt | (none) |
| feature-add | lint, build, test | (none) | atlas-* | (none) |
| refactoring | lint, build, test | (none) | atlas-fmt | (none) |
| db-change | lint, build | (none) | atlas-*, sqlc | (none) |
| documentation | (none) | (optional) | (none) | (none) |

## Auto-Detection Logic

AI Agent は以下のロジックでファイルタイプを検出し、適切な検証コマンドを決定します：

```
1. 編集されたファイルの拡張子を確認
2. 各拡張子に対応する検証コマンドを収集
3. 重複を除去して最小限のコマンドセットを構築
4. コマンドを実行
```

### 例: 複数ファイルタイプの編集

```
編集ファイル:
- internal/domain/wallet/wallet.go  → Go file
- docs/overview.md                   → MD file
- config/wallet/btc_watch.yaml       → Config file

検証コマンド:
- make go-lint     (Go file のため)
- make check-build (Go file のため)
- (MD/Config は検証不要)
```

## Skip Verification Scenarios

以下のシナリオでは検証をスキップまたは最小化できます：

| Scenario | Skippable Commands | Reason |
|----------|-------------------|--------|
| ドキュメントのみ変更 | 全ての Go 関連 | コード変更なし |
| Config のみ変更 | 全ての Go 関連 | コード変更なし |
| コメントのみ追加 | `gotest` | 機能変更なし |
| typo 修正 | `gotest` | 機能変更なし |

## Command Reference

### Make Targets

| Command | Description | When to Use |
|---------|-------------|-------------|
| `make go-lint` | golangci-lint 実行（自動修正） | Go ファイル変更時 |
| `make check-build` | 全バイナリのビルド確認 | Go ファイル変更時 |
| `make gotest` | Go テスト実行 | 機能変更時 |
| `make tidy` | go mod tidy | import 変更時 |
| `make mockery` | モック再生成 | インターフェース変更時 |
| `make sqlc` | SQLC コード再生成 | SQL/スキーマ変更時 |
| `make atlas-fmt` | Atlas スキーマフォーマット | HCL 変更時 |
| `make atlas-lint` | Atlas スキーマリント | HCL 変更時 |
| `make proto` | protobuf 再生成 | proto ファイル変更時 |

### Full Verification (All Checks)

全ての検証を実行する場合：

```bash
make go-lint && make tidy && make check-build && make gotest
```

**注意**: これは Go ファイルを編集した場合のみ必要です。

## Integration with Task Contexts

各タスクコンテキストファイルは、この Verification Matrix を参照しています：

- [bug-fix.md](./bug-fix.md) → Go 検証 + テスト
- [feature-add.md](./feature-add.md) → Go 検証 + テスト
- [refactoring.md](./refactoring.md) → Go 検証 + テスト
- [db-change.md](./db-change.md) → Go 検証 + Atlas 検証
- [documentation.md](./documentation.md) → 検証なし

## Related Documents

- [Task Contexts README](./README.md) - タスクコンテキスト一覧
- [Workflow Guidelines](../guidelines/workflow.md) - ワークフローガイドライン
- [Coding Standards](../guidelines/coding-standards.md) - コーディング規約
