###############################################################################
# Watch Only Wallet
###############################################################################
###############################################################################
# import address
###############################################################################
# import public addresses exported by keygen wallet
#make filepath=./data/pubkey/client_1535423628425011000.csv import-pubkey
.PHONY: import-pubkey
import-pubkey:
	wallet key key importing functionality ${filepath}
	#wallet -k -m 1 -i ./data/pubkey/client_1535423628425011000.csv


###############################################################################
# receipt transaction
###############################################################################
# detect receipt addresses and create unsigned transaction for client
.PHONY: create-receipt-tx
create-receipt-tx:
	wallet receipt create -fee 0.5
	#wallet -r -m 1

# WIP: only check client address
.PHONY: create-client-tx
check-client-address:
	wallet receipt create -check

# Note: debug use
# WIP: execute series of flows from creation of a receiving transaction to sending of a transaction
.PHONY: create-receipt-all
create-receipt-all:
	wallet receipt debug
	#wallet -r -m 10


# wallet receipt create
# sign xxxx

# TODO
# - wallet create [payment/receipt]
# - wallet debug  [payment/receipt]

###############################################################################
# send transaction
###############################################################################
# send signed transaction (receipt/payment/transfer)
#make filepath=./data/tx/receipt/receipt_8_signed_1534832879778945174 send-tx
.PHONY: send-tx
send-tx:
	wallet sending -file ${filepath}


###############################################################################
# monitor transaction
###############################################################################
# check status of sent tx until 6 confirmations then update status
#make acnt=client monitor-tx
.PHONY: monitor-tx
monitor-tx:
	wallet monitoring senttx -account ${acnt}

# WIP: monitor account balance
#make acnt=client monitor-balance
.PHONY: monitor-balance
monitor-balance:
	wallet monitoring balance -account ${acnt}


###############################################################################
# payment transaction
###############################################################################
# create payment request from payment table
.PHONY: create-payment-tx
create-payment-tx:
	wallet payment create -fee 0.5

# Note: debug use
# WIP: execute series of flows from creation of payment transaction to sending of a transaction
.PHONY: create-payment-all
create-payment-all:
	wallet payment debug


###############################################################################
# operation for debug / creating test data on database
###############################################################################
# create payment data on database
# Note: available after generated pub keys are imported on wallet
.PHONY: create-testdata
create-testdata:
	wallet db create

# reset payment testdata
.PHONY: reset-testdata
reset-testdata:
	wallet db reset


###############################################################################
# Bitcoin API
###############################################################################
#balance            get balance for account
#estimatefee        estimate fee
#getnetworkinfo     call getnetworkinfo
#listunspent        call listunspent
#logging            logging
#unlocktx           unlock locked transaction for unspent transaction
#validateaddress    validate address

# get balance for account
.PHONY: api-balance
api-balance:
	wallet api balance

# estimate fee
.PHONY: api-estimatefee
api-estimatefee:
	wallet api estimatefee

# call getnetworkinfo
.PHONY: api-getnetworkinfo
api-getnetworkinfo:
	wallet api getnetworkinfo

# call listunspent
.PHONY: api-listunspent
api-listunspent:
	wallet api listunspent

# logging
.PHONY: api-logging
api-logging:
	wallet api logging

# unlock locked transaction for unspent transaction
.PHONY: api-unlocktx
api-unlocktx:
	wallet api unlocktx

# validateaddress
#make addr=xxxxx api-validateaddress
.PHONY: api-validateaddress
api-validateaddress:
	wallet api validateaddress -address ${addr}
