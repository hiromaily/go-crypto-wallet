###############################################################################
# Watch Only Wallet
###############################################################################
###############################################################################
# import address
###############################################################################
# import public addresses exported by keygen wallet
#make filepath=./data/pubkey/client_1535423628425011000.csv import-pubkey
.PHONY: import-address
import-address:
	watch import address -file ${filepath}


###############################################################################
# create unsigned transaction
###############################################################################
# detect deposit addresses and create unsigned transaction for client
#make fee=0.5 create-deposit-tx
.PHONY: create-deposit-tx
create-deposit-tx:
	watch create deposit -fee ${fee}

# Note: debug use
# WIP: execute series of flows from creation of a receiving transaction to sending of a transaction
.PHONY: create-deposit-all
create-deposit-all:
	watch create deposit -debug

# create payment request from payment table
#make fee=0.5 create-payment-tx
.PHONY: create-payment-tx
create-payment-tx:
	watch create payment -fee ${fee}

# Note: debug use
# WIP: execute series of flows from creation of payment transaction to sending of a transaction
.PHONY: create-payment-all
create-payment-all:
	watch create payment -debug

# create transfer unsigned transaction among accounts
.PHONY: create-transfer-tx
create-transfer-tx:
	watch create transfer -account1 ${acnt1} -account2 ${acnt2} -amount ${amount}

###############################################################################
# send transaction
###############################################################################
# send signed transaction (deposit/payment/transfer)
#make filepath=./data/tx/deposit/deposit_8_signed_1534832879778945174 send-tx
.PHONY: send-tx
send-tx:
	watch send -file ${filepath}


###############################################################################
# monitor transaction
###############################################################################
# check status of sent tx until 6 confirmations then update status
#make acnt=client monitor-tx
.PHONY: monitor-tx
monitor-tx:
	watch monitor senttx -account ${acnt}

# WIP: monitor account balance
#make acnt=client monitor-balance
.PHONY: monitor-balance
monitor-balance:
	watch monitor balance -account ${acnt}


###############################################################################
# operation for debug / creating test data on database
###############################################################################
# create payment data on database
# Note: available after generated pub keys are imported on wallet
.PHONY: create-testdata
create-testdata:
	watch create db

# reset payment testdata
.PHONY: reset-testdata
reset-testdata:
	watch db reset


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
	watch api balance

# estimate fee
.PHONY: api-estimatefee
api-estimatefee:
	watch api estimatefee

# call getnetworkinfo
.PHONY: api-getnetworkinfo
api-getnetworkinfo:
	watch api getnetworkinfo

# call listunspent
.PHONY: api-listunspent
api-listunspent:
	watch api listunspent

# logging
.PHONY: api-logging
api-logging:
	watch api logging

# unlock locked transaction for unspent transaction
.PHONY: api-unlocktx
api-unlocktx:
	watch api unlocktx

# validateaddress
#make addr=xxxxx api-validateaddress
.PHONY: api-validateaddress
api-validateaddress:
	watch api validateaddress -address ${addr}
