###############################################################################
# Build Targets
###############################################################################

# Build on local
# - authName on sign works as account name
#.PHONY: tidy
tidy:
	# go mod verify
	go mod tidy

.PHONY: check-build
check-build: tidy
	go build -v -o /dev/null ./cmd/watch/
	go build -v -o /dev/null ./cmd/keygen/
	go build -ldflags "-X main.authName=auth1" -v -o /dev/null ./cmd/sign/

.PHONY: build-all
build-all: tidy
	go build -v -o ${GOPATH}/bin/watch ./cmd/watch/
	go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/

.PHONY: build-watch
build-watch:
	go build -v -o ${GOPATH}/bin/watch ./cmd/watch/

.PHONY: build-keygen
build-keygen:
	go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/

.PHONY: build-sign
build-sign:
	go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign ./cmd/sign/
	go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/

# Build from inside docker container
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/watch ./cmd/watch/main.go
	CGO_ENABLED=0 GOOS=linux go build -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.authName=auth1" -o /go/bin/sign ./cmd/sign/main.go

.PHONY: run-watch
run-watch:
	go run ./cmd/watch/ -conf ./config/wallet/btc/watch.yaml
