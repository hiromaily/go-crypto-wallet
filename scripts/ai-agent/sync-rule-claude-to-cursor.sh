#!/bin/bash
#
# convert-claude-to-cursor.sh
#
# Claude Code のルール (.claude/rules/*.md) を
# Cursor のルール (.cursor/rules/*.mdc) に変換するスクリプト
#
# 使用方法:
#   ./convert-claude-to-cursor.sh [SOURCE_DIR] [DEST_DIR]
#
# デフォルト:
#   SOURCE_DIR: .claude/rules
#   DEST_DIR: .cursor/rules
#
# 変換ルール:
#   - paths あり → globs: + alwaysApply: false
#   - paths なし → alwaysApply: true (グローバルルール)
#   - 最初の # 見出し → description: (frontmatter に追加)
#   - .md → .mdc (拡張子変換)
#   - サブディレクトリ構造を維持
#

set -euo pipefail

# カラー出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 使用方法
usage() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS] [SOURCE_DIR] [DEST_DIR]

Claude Code のルール (.claude/rules/*.md) を
Cursor のルール (.cursor/rules/*.mdc) に変換します。

Arguments:
  SOURCE_DIR    変換元ディレクトリ (デフォルト: .claude/rules)
  DEST_DIR      変換先ディレクトリ (デフォルト: .cursor/rules)

Options:
  -h, --help    このヘルプを表示
  -d, --dry-run 実際には変換せず、処理内容を表示
  -v, --verbose 詳細な出力を表示
  -f, --force   既存ファイルを上書き

Examples:
  $(basename "$0")
  $(basename "$0") .claude/rules .cursor/rules
  $(basename "$0") --dry-run
  $(basename "$0") -v -f /path/to/claude/rules /path/to/cursor/rules

EOF
}

# Markdown ファイルを変換
convert_file() {
    local src_file="$1"
    local dest_file="$2"
    local dry_run="$3"
    local verbose="$4"

    if [[ "$verbose" == "true" ]]; then
        log_info "Converting: $src_file → $dest_file"
    fi

    if [[ "$dry_run" == "true" ]]; then
        echo "  [DRY-RUN] Would convert: $src_file → $dest_file"
        return 0
    fi

    local content
    content=$(cat "$src_file")

    local has_frontmatter=false
    local frontmatter=""
    local body=""
    local paths_value=""
    local description=""

    # frontmatter の解析
    # ファイルの先頭が --- で始まるかチェック
    local first_line
    first_line=$(echo "$content" | head -n 1)

    if [[ "$first_line" == "---" ]]; then
        has_frontmatter=true

        # frontmatter の終了行を探す (2行目以降で最初の --- を探す)
        local end_line
        end_line=$(echo "$content" | tail -n +2 | grep -n '^---$' | head -n 1 | cut -d: -f1)

        if [[ -n "$end_line" ]]; then
            # frontmatter を抽出 (1行目と終了行を除く)
            frontmatter=$(echo "$content" | sed -n "2,${end_line}p")
            # body を抽出 (終了行の次の行から)
            local body_start=$((end_line + 2))
            body=$(echo "$content" | tail -n +"$body_start")
        else
            # frontmatter が閉じられていない場合は全体を body として扱う
            body="$content"
            has_frontmatter=false
        fi

        # paths: を抽出
        if [[ "$has_frontmatter" == "true" ]] && echo "$frontmatter" | grep -q '^paths:'; then
            paths_value=$(echo "$frontmatter" | grep '^paths:' | sed 's/^paths:[[:space:]]*//')
        fi
    else
        body="$content"
    fi

    # 最初の # 見出しを description として抽出
    local first_heading
    first_heading=$(echo "$body" | grep -m1 '^#[[:space:]]' | sed 's/^#[[:space:]]*//' || true)

    if [[ -n "$first_heading" ]]; then
        description="$first_heading"
    fi

    # 新しい frontmatter を構築
    local new_frontmatter="---"

    # description を追加
    if [[ -n "$description" ]]; then
        new_frontmatter="${new_frontmatter}
description: ${description}"
    fi

    # globs と alwaysApply を追加 (paths の有無で分岐)
    if [[ -n "$paths_value" ]]; then
        # paths あり → globs + alwaysApply: false
        new_frontmatter="${new_frontmatter}
globs: ${paths_value}
alwaysApply: false"
    else
        # paths なし → alwaysApply: true (グローバルルール)
        new_frontmatter="${new_frontmatter}
alwaysApply: true"
    fi

    # 元の frontmatter から paths 以外の項目を保持
    if [[ "$has_frontmatter" == "true" ]]; then
        local other_fields
        other_fields=$(echo "$frontmatter" | grep -v '^paths:' || true)
        if [[ -n "$other_fields" ]]; then
            new_frontmatter="${new_frontmatter}
${other_fields}"
        fi
    fi

    new_frontmatter="${new_frontmatter}
---"

    # 最終的なコンテンツを構築
    local final_content
    if [[ "$new_frontmatter" == $'---\n---' ]]; then
        # frontmatter が空の場合は body のみ
        final_content="$body"
    else
        final_content="${new_frontmatter}

${body}"
    fi

    # 出力ディレクトリを作成
    local dest_dir
    dest_dir=$(dirname "$dest_file")
    mkdir -p "$dest_dir"

    # ファイルを書き込み
    echo "$final_content" > "$dest_file"

    if [[ "$verbose" == "true" ]]; then
        log_success "Converted: $dest_file"
        if [[ -n "$paths_value" ]]; then
            echo "    paths: $paths_value → globs: $paths_value, alwaysApply: false"
        else
            echo "    paths: (none) → alwaysApply: true"
        fi
        if [[ -n "$description" ]]; then
            echo "    description: $description"
        fi
    fi
}

# メイン処理
main() {
    local dry_run=false
    local verbose=false
    local force=false
    local source_dir=".claude/rules"
    local dest_dir=".cursor/rules"

    # 引数解析
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            -d|--dry-run)
                dry_run=true
                shift
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            -f|--force)
                force=true
                shift
                ;;
            -*)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
            *)
                if [[ -z "${positional_args:-}" ]]; then
                    source_dir="$1"
                    positional_args=1
                else
                    dest_dir="$1"
                fi
                shift
                ;;
        esac
    done

    log_info "Claude → Cursor ルール変換"
    log_info "Source: $source_dir"
    log_info "Destination: $dest_dir"

    # ソースディレクトリの確認
    if [[ ! -d "$source_dir" ]]; then
        log_error "Source directory not found: $source_dir"
        exit 1
    fi

    # 変換先ディレクトリの作成
    if [[ "$dry_run" == "false" ]]; then
        mkdir -p "$dest_dir"
    fi

    # .md ファイルを検索して変換
    local count=0
    local skipped=0

    while IFS= read -r -d '' src_file; do
        # 相対パスを計算
        local rel_path="${src_file#$source_dir/}"
        # 拡張子を .mdc に変更
        local dest_file="${dest_dir}/${rel_path%.md}.mdc"

        # 既存ファイルのチェック
        if [[ -f "$dest_file" && "$force" == "false" && "$dry_run" == "false" ]]; then
            log_warn "Skipping (already exists): $dest_file"
            log_warn "  Use -f/--force to overwrite"
            ((skipped++))
            continue
        fi

        convert_file "$src_file" "$dest_file" "$dry_run" "$verbose"
        ((count++))
    done < <(find "$source_dir" -type f -name "*.md" -print0)

    echo ""
    log_success "変換完了!"
    log_info "Converted: $count files"
    if [[ $skipped -gt 0 ]]; then
        log_warn "Skipped: $skipped files"
    fi

    if [[ "$dry_run" == "true" ]]; then
        log_warn "DRY-RUN モードでした。実際の変換は行われていません。"
    fi
}

main "$@"
