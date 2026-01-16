#!/bin/bash
#
# sync-rule-claude-to-cursor.sh
#
# Converts Claude Code rules (.claude/rules/*.md) to
# Cursor rules (.cursor/rules/*.mdc)
#
# Usage:
#   ./sync-rule-claude-to-cursor.sh [SOURCE_DIR] [DEST_DIR]
#
# Defaults:
#   SOURCE_DIR: .claude/rules
#   DEST_DIR: .cursor/rules
#
# Conversion Rules:
#   - paths present → globs: + alwaysApply: false
#   - paths absent → alwaysApply: true (global rule)
#   - First # heading → description: (added to frontmatter)
#   - .md → .mdc (extension conversion)
#   - Preserves subdirectory structure
#

set -euo pipefail

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log functions
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

# Usage
usage() {
	cat <<EOF
Usage: $(basename "$0") [OPTIONS] [SOURCE_DIR] [DEST_DIR]

Converts Claude Code rules (.claude/rules/*.md) to
Cursor rules (.cursor/rules/*.mdc).

Arguments:
  SOURCE_DIR    Source directory (default: .claude/rules)
  DEST_DIR      Destination directory (default: .cursor/rules)

Options:
  -h, --help    Show this help message
  -d, --dry-run Show what would be done without actually converting
  -v, --verbose Show detailed output
  -f, --force   Overwrite existing files

Examples:
  $(basename "$0")
  $(basename "$0") .claude/rules .cursor/rules
  $(basename "$0") --dry-run
  $(basename "$0") -v -f /path/to/claude/rules /path/to/cursor/rules

EOF
}

# Convert markdown file
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

	# Parse frontmatter
	# Check if file starts with ---
	local first_line
	first_line=$(echo "$content" | head -n 1)

	if [[ "$first_line" == "---" ]]; then
		has_frontmatter=true

		# Find the closing --- (first --- after line 2)
		local end_line
		end_line=$(echo "$content" | tail -n +2 | grep -n '^---$' | head -n 1 | cut -d: -f1)

		if [[ -n "$end_line" ]]; then
			# Extract frontmatter (excluding first and last ---)
			frontmatter=$(echo "$content" | sed -n "2,${end_line}p")
			# Extract body (from after the closing ---)
			local body_start=$((end_line + 2))
			body=$(echo "$content" | tail -n +"$body_start")
		else
			# Frontmatter not properly closed, treat entire content as body
			body="$content"
			has_frontmatter=false
		fi

		# Extract paths: (supports both single-line and multi-line YAML array format)
		if [[ "$has_frontmatter" == "true" ]] && echo "$frontmatter" | grep -q '^paths:'; then
			local paths_line
			paths_line=$(echo "$frontmatter" | grep '^paths:' | sed 's/^paths:[[:space:]]*//')

			if [[ -n "$paths_line" ]]; then
				# Single-line format: paths: ["value1", "value2"] or paths: "value"
				paths_value="$paths_line"
			else
				# Multi-line format: paths: followed by indented array elements
				# Get the line number of paths:
				local paths_line_num
				paths_line_num=$(echo "$frontmatter" | grep -n '^paths:' | cut -d: -f1)

				# Collect indented lines starting with - after paths:
				local array_items=""
				local in_paths_array=false
				local line_num=0

				while IFS= read -r line; do
					((line_num++))
					if [[ $line_num -eq $paths_line_num ]]; then
						in_paths_array=true
						continue
					fi
					if [[ "$in_paths_array" == "true" ]]; then
						# Collect indented lines starting with -
						if [[ "$line" =~ ^[[:space:]]+-[[:space:]] ]]; then
							if [[ -n "$array_items" ]]; then
								array_items="${array_items}
${line}"
							else
								array_items="$line"
							fi
						else
							# End of array
							break
						fi
					fi
				done <<<"$frontmatter"

				if [[ -n "$array_items" ]]; then
					# Preserve multi-line format
					paths_value="
${array_items}"
				fi
			fi
		fi
	else
		body="$content"
	fi

	# Extract first # heading as description
	local first_heading
	first_heading=$(echo "$body" | grep -m1 '^#[[:space:]]' | sed 's/^#[[:space:]]*//' || true)

	if [[ -n "$first_heading" ]]; then
		description="$first_heading"
	fi

	# Build new frontmatter
	local new_frontmatter="---"

	# Add description
	if [[ -n "$description" ]]; then
		new_frontmatter="${new_frontmatter}
description: ${description}"
	fi

	# Add globs and alwaysApply (based on presence of paths)
	if [[ -n "$paths_value" ]]; then
		# paths present → globs + alwaysApply: false
		new_frontmatter="${new_frontmatter}
globs: ${paths_value}
alwaysApply: false"
	else
		# paths absent → alwaysApply: true (global rule)
		new_frontmatter="${new_frontmatter}
alwaysApply: true"
	fi

	# Preserve other fields from original frontmatter
	# (excluding paths: line and its indented array elements)
	if [[ "$has_frontmatter" == "true" ]]; then
		local other_fields=""
		local skip_array=false
		local first_field=true

		while IFS= read -r line; do
			if [[ "$line" =~ ^paths: ]]; then
				# Skip paths: line
				# If value is empty, also skip following array elements
				local paths_inline_value
				paths_inline_value=$(echo "$line" | sed 's/^paths:[[:space:]]*//')
				if [[ -z "$paths_inline_value" ]]; then
					skip_array=true
				fi
				continue
			fi

			if [[ "$skip_array" == "true" ]]; then
				# Skip indented lines starting with - (array elements)
				if [[ "$line" =~ ^[[:space:]]+-[[:space:]] ]]; then
					continue
				else
					# End of array
					skip_array=false
				fi
			fi

			# Preserve other fields (including blank lines)
			if $first_field; then
				other_fields="$line"
				first_field=false
			else
				other_fields="${other_fields}
${line}"
			fi
		done <<<"$frontmatter"

		if [[ -n "$other_fields" ]]; then
			new_frontmatter="${new_frontmatter}
${other_fields}"
		fi
	fi

	new_frontmatter="${new_frontmatter}
---"

	# Build final content
	local final_content
	if [[ "$new_frontmatter" == $'---\n---' ]]; then
		# If frontmatter is empty, use body only
		final_content="$body"
	else
		final_content="${new_frontmatter}

${body}"
	fi

	# Create output directory
	local dest_dir
	dest_dir=$(dirname "$dest_file")
	mkdir -p "$dest_dir"

	# Write file
	echo "$final_content" >"$dest_file"

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

# Main process
main() {
	local dry_run=false
	local verbose=false
	local force=false
	local source_dir=".claude/rules"
	local dest_dir=".cursor/rules"

	# Parse arguments
	while [[ $# -gt 0 ]]; do
		case $1 in
		-h | --help)
			usage
			exit 0
			;;
		-d | --dry-run)
			dry_run=true
			shift
			;;
		-v | --verbose)
			verbose=true
			shift
			;;
		-f | --force)
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

	log_info "Claude → Cursor rule conversion"
	log_info "Source: $source_dir"
	log_info "Destination: $dest_dir"

	# Verify source directory exists
	if [[ ! -d "$source_dir" ]]; then
		log_error "Source directory not found: $source_dir"
		exit 1
	fi

	# Create destination directory
	if [[ "$dry_run" == "false" ]]; then
		mkdir -p "$dest_dir"
	fi

	# Find and convert .md files
	local count=0
	local skipped=0

	while IFS= read -r -d '' src_file; do
		# Calculate relative path
		local rel_path="${src_file#$source_dir/}"
		# Change extension to .mdc
		local dest_file="${dest_dir}/${rel_path%.md}.mdc"

		# Check for existing file
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
	log_success "Conversion complete!"
	log_info "Converted: $count files"
	if [[ $skipped -gt 0 ]]; then
		log_warn "Skipped: $skipped files"
	fi

	if [[ "$dry_run" == "true" ]]; then
		log_warn "DRY-RUN mode was enabled. No actual conversion was performed."
	fi
}

main "$@"
