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
	wallet import key -file ${filepath}
	#wallet -k -m 1 -i ./data/pubkey/client_1535423628425011000.csv


###############################################################################
# create unsigned transaction
###############################################################################
# detect receipt addresses and create unsigned transaction for client
#make fee=0.5 create-receipt-tx
.PHONY: create-receipt-tx
create-receipt-tx:
	wallet create receipt -fee ${fee}
	#wallet -r -m 1

# Note: debug use
# WIP: execute series of flows from creation of a receiving transaction to sending of a transaction
.PHONY: create-receipt-all
create-receipt-all:
	wallet create receipt -debug
	#wallet -r -m 10

# create payment request from payment table
#make fee=0.5 create-payment-tx
.PHONY: create-payment-tx
create-payment-tx:
	wallet create payment -fee ${fee}

# Note: debug use
# WIP: execute series of flows from creation of payment transaction to sending of a transaction
.PHONY: create-payment-all
create-payment-all:
	wallet create payment -debug

# create transfer unsigned transaction among accounts
.PHONY: create-transfer-tx
create-transfer-tx:
	wallet create transfer -account1 ${acnt1} -account2 ${acnt2}

###############################################################################
# send transaction
###############################################################################
# send signed transaction (receipt/payment/transfer)
#make filepath=./data/tx/receipt/receipt_8_signed_1534832879778945174 send-tx
.PHONY: send-tx
send-tx:
	wallet send -file ${filepath}


###############################################################################
# monitor transaction
###############################################################################
# check status of sent tx until 6 confirmations then update status
#make acnt=client monitor-tx
.PHONY: monitor-tx
monitor-tx:
	wallet monitor senttx -account ${acnt}

# WIP: monitor account balance
#make acnt=client monitor-balance
.PHONY: monitor-balance
monitor-balance:
	wallet monitor balance -account ${acnt}


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
