package sign

// Signer is Signer service
type Signer interface {
	SignTx(filePath string) (string, bool, string, error)
}

// FullPubkeyExporter is FullPubkeyExporter service
type FullPubkeyExporter interface {
	ExportFullPubkey() (string, error)
}
