# Ganache
- [Ganache](https://www.trufflesuite.com/ganache)

## Run by docker
```
docker compose -f docker-compose.eth.yml up ganache
```

## Setup ethereum environment with Ganache
1. run ganache
2. create sql to register private key and address displayed on console by running ganache
    - run db by `docker compose up watch-db keygen-db sign-db`
3. run sql to insert data into database. [./docker/mysql/insert/ganache.example.sql](https://github.com/hiromaily/go-crypto-wallet/blob/master/docker/mysql/insert/ganache.example.sql)
4. import private key from database
    - run `direnv allow` to set environment variable to run cli. please install [direnv](https://direnv.net/) if not.
    - run command like the below.
    ```
    keygen -coin eth import privkey -account client
    keygen -coin eth import privkey -account deposit
    keygen -coin eth import privkey -account payment
    keygen -coin eth import privkey -account stored
    ```
5. export addresses
```
keygen -coin eth export address -account client
keygen -coin eth export address -account deposit
keygen -coin eth export address -account payment
keygen -coin eth export address -account stored
```

6. import addresses into watch wallet
```
watch -coin eth import address -file xxxxx
```

## Make transaction
just run scripts
```
# create depost
./scripts/operation/create-eth-tx-deposit.sh

# create transfer
./scripts/operation/create-eth-tx-transfer.sh
```
