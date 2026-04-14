## Wallet Type

This is explained for BTC/BCH for now.
There are mainly 3 wallets separately and these wallets are expected to be installed in each different devices.

### 1.Watch only wallet

- Only this wallet run online to access to BTC/BCH Nodes.
- Only pubkey address is stored. Private key is NOT stored for security reason. That's why this is called `watch only wallet`.
- Major functionalities are
  - creating unsigned transaction
  - sending signed transaction
  - monitoring transaction status.

### 2.Keygen wallet as cold wallet

- Key management functionalities for accounts.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts
  - generating keys based on `HD Wallet`
  - generating multisig addressed according to account setting
  - exporting pubkey addresses as csv file which is imported from `Watch only wallet`
  - signing on unsigned transaction as first sign. However, multisig addresses could not be completed by only this wallet.

### 3.Sign wallet as cold wallet (Auth wallet)

- The internal authorization operators would use this wallet to sign on unsigned transaction for multisig addresses.
- Each of operators would be given own authorization account and Sing wallet apps.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts for own auth account
  - generating keys based on `HD Wallet` for own auth account
  - exporting full-pubkey addresses as csv file which is imported from `Keygen wallet` to generate multisig address
  - signing on unsigned transaction as second or more signs for multisig addresses.
