-- Cold database: MuSig2 nonce commitments table
-- Used by keygen and sign wallets for secure nonce storage

CREATE TABLE musig2_nonces (
  id             BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  signer_id      VARCHAR(255) NOT NULL COMMENT 'Signer identifier',
  transaction_id VARCHAR(255) NOT NULL COMMENT 'Transaction identifier',
  public_nonce   BINARY(66) NOT NULL COMMENT 'Public nonce (66 bytes: two 33-byte compressed EC points R1||R2)',
  is_used        BOOL NOT NULL DEFAULT false COMMENT 'true: nonce has been used in signing',
  created_at     DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'creation date',
  used_at        DATETIME DEFAULT NULL COMMENT 'date when nonce was marked as used',
  PRIMARY KEY (id),
  UNIQUE KEY idx_signer_transaction (signer_id, transaction_id) COMMENT 'Prevent duplicate nonces per signer-tx pair (CRITICAL for security)',
  INDEX idx_transaction_id (transaction_id),
  INDEX idx_is_used (is_used),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='MuSig2 nonce commitments for secure storage';
