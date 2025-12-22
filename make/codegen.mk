###############################################################################
# Code Generator Targets
###############################################################################

###############################################################################
# sqlc
#------------------------------------------------------------------------------
# Generate Go code from SQL queries using sqlc
# Schemas: tools/sqlc/schemas/*.sql
# Queries: tools/sqlc/queries/*.sql
# Output: pkg/db/rdb/sqlcgen/
#------------------------------------------------------------------------------
.PHONY: sqlc
sqlc:
	cd tools/sqlc && sqlc generate

###############################################################################
# sqlboiler (DEPRECATED - Migrating to sqlc)
#------------------------------------------------------------------------------
# To generate all schema, modify `docker/mysql/watch/init.d/watch.sql` according to comments
# Then recreate database
# ```
# $ docker compose rm -f -s watch-db
# $ docker volume rm -f go-crypto-wallet_watch-db
# $ docker compose up watch-db
# ```
# Make sure `watch-db` includes tables of keygen-db/sign-db
# Then, run `make sqlboiler`
# Make sure `make build` works
# Revert `docker/mysql/watch/init.d/watch.sql`
#------------------------------------------------------------------------------
.PHONY: sqlboiler
sqlboiler:
	sqlboiler --wipe mysql

.PHONY: sqlboiler-with-template
sqlboiler-with-template:
	sqlboiler --wipe mysql \
	--templates \
	${GOPATH}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.8.6/templates/main, \
	${GOPATH}/pkg/mod/github.com/volatiletech/sqlboiler/v4@v4.8.6/templates/test, \
	templates

.PHONY: generate-abi
generate-abi:
	abigen --abi ./data/contract/token.abi --pkg contract --type Token --out ./pkg/contract/token-abi.go

###############################################################################
# Protocol Buffer
#------------------------------------------------------------------------------
# run `make install-proto-plugin` in advance
###############################################################################
.PHONY: get-third-proto
get-third-proto:
	./scripts/get_third_proto.sh

.PHONY: lint-proto
lint-proto:
	prototool lint

.PHONY:protoc-go
protoc-go: clean-pb
	$(PROTOC_BIN) \
	--go_out=./pkg/wallet/api/xrpgrp/xrp/ --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/wallet/api/xrpgrp/xrp/ --go-grpc_opt=paths=source_relative  \
	--proto_path=./data/proto/rippleapi \
	--proto_path=./data/proto/third_party \
	data/proto/**/*.proto

.PHONY: clean-pb
clean-pb:
	rm -rf pkg/wallet/api/xrpgrp/xrp/*.pb.go
