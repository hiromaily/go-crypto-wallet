###############################################################################
# Build Targets
###############################################################################

# Build on local
# - appVersion is injected via LDFLAGS (defined in vars.mk)
# - authName on sign works as account name
#.PHONY: tidy
tidy:
	# go mod verify
	go mod tidy

.PHONY: check-build
check-build: tidy
	go build -ldflags "$(LDFLAGS)" -v -o /dev/null ./cmd/watch/
	go build -ldflags "$(LDFLAGS)" -v -o /dev/null ./cmd/keygen/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -v -o /dev/null ./cmd/sign/

.PHONY: build-all
build-all: tidy
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/watch ./cmd/watch/
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/

.PHONY: build-watch
build-watch:
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/watch ./cmd/watch/

.PHONY: build-keygen
build-keygen:
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/keygen ./cmd/keygen/

.PHONY: build-sign
build-sign:
	go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/

# Build from inside docker container
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o /go/bin/watch ./cmd/watch/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -o /go/bin/sign ./cmd/sign/main.go

# Show current version from git tag
.PHONY: show-version
show-version:
	@echo "VERSION: $(VERSION)"

.PHONY: run-watch
run-watch:
	go run ./cmd/watch/ -conf ./config/wallet/btc/watch.yaml
