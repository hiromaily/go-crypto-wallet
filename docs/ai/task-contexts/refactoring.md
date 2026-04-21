<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/ai/task-contexts/refactoring.tpl.md · Run `make docs` to regenerate.
-->

---
task_type: refactoring
description: リファクタリングタスク用のコンテキスト
version: 1.0.0
---

# Task: Refactoring (リファクタリング)

## When to Use

- コードの構造を改善する場合（機能変更なし）
- Clean Architectureに準拠させるためのレイヤー移動
- 重複コードの削減・共通化
- 命名規則の統一
- パッケージ構造の整理
- 技術的負債の解消

## Required Context Documents

### Must Read (必須)

以下のドキュメントは必ず読み込んでください：

| Document | Path | Purpose |
|----------|------|---------|
| Architecture | `docs/guidelines/architecture.md` | Clean Architecture原則 |
| Coding Standards | `docs/guidelines/coding-standards.md` | コーディング規約 |
| Refactoring Plan | `docs/issues/REFACTORING_PLAN.md` | リファクタリング計画 |
| Testing | `docs/guidelines/testing.md` | 既存テストの維持 |

### Conditional Read (条件付き)

条件に応じて追加で読み込むドキュメント：

| Condition | Document | Path |
|-----------|----------|------|
| internalパッケージ移動 | Internal Guidelines | `internal/AGENTS.md` |
| pkgパッケージ移動 | Pkg Guidelines | `pkg/AGENTS.md` |
| モック関連の変更 | Code Generation | `docs/guidelines/code-generation.md` |
| DB関連の変更 | Database | `docs/database/db-management.md` |

## Task-Specific Rules

### 1. 機能変更なし原則

```
❌ リファクタリングと同時に機能追加/変更
✅ リファクタリングのみ（外部から見た動作は同一）
```

- 入出力の動作は変更しない
- 既存のテストが全てパスすること

### 2. 段階的な変更

```
❌ 一度に大規模な変更
✅ 小さな変更を積み重ねる
```

- 各コミットは独立して検証可能
- ロールバック可能な単位で進める

### 3. テストの維持

```
❌ テストを壊す/スキップする
✅ テストを維持し、必要に応じて更新
```

- 既存テストが引き続きパスすること
- テスト自体のリファクタリングは別PRで

### 4. Clean Architecture原則の適用

```
依存の方向: Interface Adapters → Application → Domain
                    ↓
             Infrastructure
```

- Domain層は外部依存なし
- Infrastructure依存はPortsで抽象化

### 5. モック設定の更新

コードを移動する場合、`.mockery.yaml` の設定も更新が必要な場合があります：

```yaml
# .mockery.yaml の確認が必要なケース
- インターフェースの移動
- パッケージ名の変更
- 新しいインターフェースの追加
```

## Pre-Task Checklist

- [ ] リファクタリングの目的と範囲を明確にしたか
- [ ] 既存のテストカバレッジを確認したか
- [ ] 依存関係を把握したか（どのパッケージがこのコードを使用しているか）
- [ ] 段階的な変更計画を立てたか
- [ ] REFACTORING_PLAN.md の関連タスクを確認したか
- [ ] モック設定への影響を確認したか

## Verification Commands

```bash
# 必須の検証コマンド（各ステップ後に実行）
make go-lint      # リンターチェック
make tidy         # 依存関係の整理
make check-build  # ビルド確認
make go-test       # テスト実行

# モック再生成（インターフェース変更時）
make mockery
```

## Refactoring Workflow

### Step 1: 現状分析

1. リファクタリング対象のコードを特定
2. 依存関係を調査（このコードを使用している箇所）
3. 既存のテストカバレッジを確認
4. REFACTORING_PLAN.md で関連タスクを確認

### Step 2: 計画

1. 変更を小さな単位に分割
2. 各ステップの検証方法を決定
3. リスクの高い変更を特定

### Step 3: 段階的実装

```
Step A: 準備
  - 必要なテストを追加（カバレッジ向上）
  - 検証: make go-test

Step B: リファクタリング実行
  - 小さな単位で変更
  - 検証: make go-lint && make go-test

Step C: クリーンアップ
  - 不要なコードの削除
  - インポートの整理
  - 検証: make go-lint && make tidy
```

### Step 4: 最終検証

```bash
# 全ての検証コマンドを実行
make go-lint && make tidy && make check-build && make go-test
```

### Step 5: PRドラフト

```bash
git add .
git commit -m "refactor: {what was refactored}

- {change detail 1}
- {change detail 2}

No functional changes."

gh pr create --title "Refactor: {description}"
```

## Common Refactoring Patterns

### Pattern 1: レイヤー間の移動

```
Before: internal/infrastructure/some_logic.go
After:  internal/domain/some_logic.go

Steps:
1. 新しい場所にファイルをコピー
2. 依存関係を更新
3. インポートを更新
4. テスト実行で確認
5. 古いファイルを削除
```

### Pattern 2: インターフェース抽出

```
Before: 具象クラスに直接依存
After:  インターフェース経由で依存

Steps:
1. Portsにインターフェースを定義
2. 依存側をインターフェースに変更
3. DIで具象クラスを注入
4. .mockery.yaml を更新
5. make mockery でモック生成
```

### Pattern 3: 共通処理の抽出

```
Before: 複数箇所で同じロジック
After:  共通関数/パッケージに集約

Steps:
1. 共通ロジックを特定
2. 適切な場所に共通関数を作成
3. 各箇所を共通関数呼び出しに置換
4. テストで動作確認
```

## Testing Requirements

### 必須

- [ ] 既存のテストが全てパスすること
- [ ] リンターエラーがないこと

### 推奨

- [ ] リファクタリング前にテストカバレッジを向上
- [ ] 複雑なロジックには追加テスト

## Examples

### Example 1: Infrastructure → Domain への移動

```
User: "BTCのトランザクションロジックをDomain層に移動して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. 現在の実装を確認（internal/infrastructure/wallet/btc/）
3. ドメインロジックとインフラロジックを分離
4. Domain層にエンティティ/値オブジェクトを作成
5. Infrastructure層からDomain層への依存を作成
6. テスト実行で動作確認
7. 古いコードを削除
8. 検証コマンド実行
```

### Example 2: インターフェース抽出

```
User: "Repository実装をインターフェースで抽象化して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. 現在の具象クラスを確認
3. internal/application/ports/ にインターフェース定義
4. 依存側をインターフェースに変更
5. .mockery.yaml を更新
6. make mockery でモック生成
7. DIコンテナを更新
8. テスト実行
9. 検証コマンド実行
```

### Example 3: パッケージ構造の整理

```
User: "BTCとBCHで重複しているコードを共通化して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. 重複コードを特定
3. 共通パッケージの配置場所を決定
4. 共通ロジックを抽出
5. 各実装から共通ロジックを呼び出し
6. テスト実行
7. 検証コマンド実行
```

### Pattern 4: Splitting a Monolithic Infrastructure Client

When splitting a monolithic struct (e.g. `XRP` with both public + admin methods) into focused
sub-clients, read these files **upfront** to map the full dependency chain before writing any code:

```
Upfront reading order:
1. internal/application/ports/api/<chain>/interface.go  — current interfaces
2. internal/infrastructure/api/<chain>/xrp.go           — monolithic struct
3. internal/di/container.go                             — how the client is injected
4. internal/interface-adapters/wallet/<chain>/<chain>.go — wallet adapter(s) (keygen, watch, sign)
5. internal/interface-adapters/cli/<wallet>/api/<chain>/ — CLI commands using the client
6. internal/application/usecase/<wallet>/<chain>/        — use cases (check constructor args)
7. internal/infrastructure/api/<chain>/testutil/xrp.go  — testutil factory (often missed!)
```

Key decision: for each method on the monolith, assign it to public-only, admin-only, or both.

Steps:
1. Read all 7 files above
2. Define new `XRPPublicClient` and `XRPAdminClient` port interfaces in `ports/api/<chain>/`
3. Create `internal/infrastructure/api/<chain>/public/` and `admin/` subdirectories
4. Update DI container to construct two clients (public singleton, admin on-demand)
5. Update wallet adapters + CLI to use `XRPPublicClient`
6. Update `.mockery.yaml`, run `make mockery`
7. Update `testutil/xrp.go` (often missed — see infrastructure-layer.md)
8. Update tests (combined mock struct if use cases compose interfaces — see mockery.md)

## Anti-Patterns to Avoid

### ❌ Big Bang Refactoring

```
問題: 一度に全てを変更
リスク: デバッグ困難、ロールバック不可

対策: 小さな変更を積み重ねる
```

### ❌ Refactoring Without Tests

```
問題: テストなしでリファクタリング
リスク: 動作変更に気づかない

対策: 先にテストを追加してからリファクタリング
```

### ❌ Mixing Refactoring and Features

```
問題: リファクタリングと機能追加を混在
リスク: レビュー困難、バグ混入

対策: 別々のPRに分ける
```

## Related Documents

- [Architecture Guidelines](../../guidelines/architecture.md) - Clean Architecture詳細
- [Refactoring Plan](https://github.com/hiromaily/go-crypto-wallet/blob/main/docs/issues/REFACTORING_PLAN.md) - プロジェクトのリファクタリング計画
- [Code Generation](../../guidelines/code-generation.md) - モック生成など
- `internal/AGENTS.md` - internalパッケージ構造
