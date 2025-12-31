
.PHONY: setup
setup-mac:
	brew install protobuf
	brew tap ktr0731/evans
	brew install evans

.PHONY: protoc-ts
protoc-ts:
	./scripts/protoc-ts.sh

.PHONY: update-dependencies
update-dependencies:
	#yarn yarn-upgrade-all
	yarn upgrade-interactive --latest
	yarn upgrade --latest

.PHONY: run
run:
	yarn run dev
