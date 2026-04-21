### Build the Wallets

Build all three wallet binaries:

```bash
make build
```

This produces `watch`, `keygen`, and `sign` binaries. The sign binary embeds the authorizer name at build time:

```bash
# Manual build (equivalent to make build)
go build -v -o ${GOPATH}/bin/watch  ./cmd/watch/
go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
```

Configuration files are in [`./config/wallet/*.toml`](https://github.com/hiromaily/go-crypto-wallet/tree/main/config/wallet).
