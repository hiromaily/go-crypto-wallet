
###############################################################################
# Cold wallet2 (Authorizationのキー/ Receipt/PaymentのMultisigアドレス管理)
###############################################################################

# seedを生成する
gen-seed2:
	coldwallet2 -k -m 1


# Authorizationのkeyを生成する
gen-authorization-key:
	coldwallet2 -k -m 13


# Authorizationのprivate keyをcoldwalletに登録する
add-authorization-priv-key:
	coldwallet2 -k -m 23


# ReceiptのPublicアドレス(full public key)をimportする
import-receipt-pub-key:
	coldwallet2 -k -m 33 -i ./data/pubkey/receipt_1535613888391656000.csv

# PaymentのPublicアドレス(full public key)をimportする
import-payment-pub-key:
	coldwallet2 -k -m 34 -i ./data/pubkey/payment_1535613934762230000.csv


# Receiptのmultisigアドレスを生成し、登録する
add-multisig-receipt:
	coldwallet2 -k -m 50

# Paymentのmultisigアドレスを生成し、登録する
add-multisig-payment:
	coldwallet2 -k -m 51


# Receiptのmultisigアドレスをexportする
export-multisig-receipt:
	coldwallet2 -k -m 60

# Paymentのmultisigアドレスをexportする
export-multisig-payment:
	coldwallet2 -k -m 61
