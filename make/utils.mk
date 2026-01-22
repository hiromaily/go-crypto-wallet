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
