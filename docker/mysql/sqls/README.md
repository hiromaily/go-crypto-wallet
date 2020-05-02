# Database schema

## Caution when generating models by SqlBoiler

##### 1. remove comment out in watch.sql 
```
-- source /sqls/definition_keygen.sql
```

##### 2. modify whitelist in sqlboiler.toml
```
whitelist = [
    "tx",
    "tx_input",
    "tx_output",
    "payment_request",
    "address",
    "seed",
    "account_key",
    "multisig_history"
]
```

##### 3. tweak templates/*.go.tpl  
- sqlboiler doesn't generate new file because of bug in SqlBoiler.