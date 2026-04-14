### Common Setup

1. install Golang, Docker
2. build `watch`, `keygen`, `auth` wallets

- only each sign wallet includes corresponding account name as `authName` into binary

```sh
make build
 or
go build -v -o ${GOPATH}/bin/watch ./cmd/watch/
go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
```

1. configure config files in [./config/wallet/*.toml](https://github.com/hiromaily/go-crypto-wallet/tree/main/config/wallet)
2. run Database containers

```
docker compose up wallet-mysql
```
