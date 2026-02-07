###############################################################################
# Utility Targets
###############################################################################

# clean
.PHONY: clean
clean: rm-db-volumes rm-local-wallet-dat

# get timestamp
.PHONY: timestamp
timestamp:
	@echo $(timestamp)

# remove local binaries
.PHONY: rm-local-binaries
rm-local-binaries:
	rm -rf watch keygen sign

# remove local files
.PHONY: rm-files
rm-files:
	rm -rf ./data/btc/address/*.csv
	rm -rf ./data/btc/pubkey/*.csv
	rm -rf ./data/btc/tx/deposit/*
	rm -rf ./data/btc/tx/payment/*
	rm -rf ./data/btc/tx/transfer/*
	touch ./data/btc/tx/deposit/.gitkeep
	touch ./data/btc/tx/payment/.gitkeep
	touch ./data/btc/tx/transfer/.gitkeep

# find empty directories
.PHONY: find-empty-dirs
find-empty-dirs:
	find . -type d -empty

# clear gha cache within 1 day not accessed
.PHONY: clear-gha-cache-1d
clear-gha-cache-1d:
	gh cache list --limit 500 --json id,createdAt,lastAccessedAt | \
	jq -r '
		.[]
		| select(.createdAt == .lastAccessedAt)
		| select((now - (.createdAt | fromdateiso8601)) > 1*24*3600)
		| .id
	' | \
	xargs -r -I {} gh cache delete {}

# clear gha cache
.PHONY: clear-gha-cache
clear-gha-cache:
	gh cache list --limit 500 --json id,createdAt,lastAccessedAt | \
	jq -r '.[] | select(.createdAt == .lastAccessedAt) | .id' | \
	xargs -I {} gh cache delete {}
