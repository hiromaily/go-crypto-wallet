###############################################################################
# AI Targets
###############################################################################
PHONY: create-claudemd-symlink
create-claudemd-symlink:
	@find . -name "AGENTS.md" -type f | while read agents_file; do \
		dir=$$(dirname "$$agents_file"); \
		claude_file="$$dir/CLAUDE.md"; \
		ln -sf "AGENTS.md" "$$claude_file"; \
		echo "Created symlink: $$claude_file -> $$dir/AGENTS.md"; \
	done

