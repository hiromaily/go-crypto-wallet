###############################################################################
# keygen wallet
###############################################################################

# development
develop:
	coldwallet1 -d

# seedを生成する
gen-seed:
	coldwallet1 -k -m 1


# Clientのkeyを生成する
gen-client-key:
	coldwallet1 -k -m 10

# Receiptのkeyを生成する
gen-receipt-key:
	coldwallet1 -k -m 11

# Paymentのkeyを生成する
gen-payment:
	coldwallet1 -k -m 12


# Clientのprivate keyをcoldwalletに登録する
add-client-priv-key:
	coldwallet1 -k -m 20

# Receiptのprivate keyをcoldwalletに登録する
add-receipt-priv-key:
	coldwallet1 -k -m 21

# Paymentのprivate keyをcoldwalletに登録する
add-payment-priv-key:
	coldwallet1 -k -m 22


# Clientのpubアドレスをexportする
export-client-pub-key:
	coldwallet1 -k -m 30

# Receiptのpubアドレスをexportする
export-receipt-pub-key:
	coldwallet1 -k -m 31

# Paymentのpubアドレスをexportする
export-payment-pub-key:
	coldwallet1 -k -m 32


# Receiptのmultisigアドレスをimportする
import-receipt-multisig-address:
	coldwallet1 -k -m 40

# Paymentのmultisigアドレスをimportする
import-payment-multisig-address:
	coldwallet1 -k -m 41

