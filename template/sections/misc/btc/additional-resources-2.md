## Additional Resources

### Documentation

- [BIP 174: PSBT Specification](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [btcd PSBT Package](https://pkg.go.dev/github.com/btcsuite/btcd/btcutil/psbt)
- [PSBT Implementation Details](./implementation)
- [PSBT User Guide](./user-guide)

### Code References

- Bitcoin API: `internal/infrastructure/api/btc/btc/psbt.go`
- File Storage: `internal/infrastructure/storage/file/transaction.go`
- Watch Use Cases: `internal/application/usecase/watch/btc/`
- Keygen Use Cases: `internal/application/usecase/keygen/btc/`
- Sign Use Cases: `internal/application/usecase/sign/btc/`

### Tools

- [Bitcoin Core](https://bitcoincore.org/) - PSBT decoding/analysis
- [btcdeb](https://github.com/bitcoin-core/btcdeb) - Bitcoin script debugger
- [PSBT Toolkit](https://github.com/bitcoin/bitcoin/blob/master/doc/psbt.md) - Bitcoin Core PSBT tools

---

**Last Updated**: 2025-01-27
**Version**: 1.0 (PSBT Phase 2 Complete)
