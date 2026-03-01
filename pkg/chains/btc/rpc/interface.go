package rpc

// BTCRPC defines all Bitcoin RPC operations exposed by Client.
type BTCRPC interface {
	// address.go
	GetAddressInfo(addr string) (*GetAddressInfoResult, error)
	ValidateAddress(addr string) (*ValidateAddressResult, error)
	GetAddressesByLabel(labelName string) (map[string]Purpose, error)

	// balance.go
	GetBalance(confirmationBlock int) (float64, error)

	// descriptor.go
	ImportDescriptors(requests []ImportDescriptorsRequest) ([]ImportDescriptorsResponse, error)
	ImportMulti(requests []ImportMultiRequest, options *ImportMultiOptions) ([]ImportMultiResponse, error)
	GetDescriptorInfo(descriptor string) (*DescriptorInfo, error)
	ListDescriptors(privateDescriptors bool) (*ListDescriptorsResult, error)

	// fee.go
	EstimateSmartFee(confirmationBlock int) (float64, error)

	// import.go
	ImportAddress(address, label string, rescan bool) error
	ImportPrivKey(wifKey, label string, rescan bool) error

	// label.go
	SetLabel(addr, label string) error

	// logging.go
	Logging() (*LoggingResult, error)

	// multisig.go
	AddMultisigAddress(
		requiredSigs int, addresses []string, accountName, addressType string,
	) (*AddMultisigAddressResult, error)

	// networkinfo.go
	GetNetworkInfo() (*GetNetworkInfoResult, error)
	GetBlockchainInfo() (*GetBlockchainInfoResult, error)
	GetBlockCount() (int64, error)

	// transaction.go
	GetTransaction(txID string) (*GetTransactionResult, error)
	DecodeRawTransaction(hexTx string) (*TxRawResult, error)
	FundRawTransaction(hexTx string, opts *FundRawTransactionOptions) (*FundRawTransactionResult, error)
	SignRawTransactionWithWallet(hexTx string, prevTxs []PrevTx) (*SignRawTransactionResult, error)
	SignRawTransactionWithKey(hexTx string, privKeys []string, prevTxs []PrevTx) (*SignRawTransactionResult, error)
	CreateRawTransaction(inputs []TxInput, outputs map[string]float64) (string, error)
	SendRawTransaction(hexTx string) (string, error)
	GetTxOut(txID string, index uint32, mempool bool) (*GetTxOutResult, error)

	// wallet.go
	CreateWallet(fileName string, opts *CreateWalletOptions) error
	LoadWallet(fileName string) error
	UnloadWallet(fileName string) error
}
