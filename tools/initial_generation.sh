#!/bin/sh

#After `docker-compose up`

#coldwallet1
coldwallet1 -k -m 1

coldwallet1 -k -m 10
coldwallet1 -k -m 11
coldwallet1 -k -m 12

coldwallet1 -k -m 20
coldwallet1 -k -m 21
coldwallet1 -k -m 22

coldwallet1 -k -m 30
coldwallet1 -k -m 31
coldwallet1 -k -m 32

#coldwallet2
coldwallet2 -k -m 1
coldwallet2 -k -m 13
coldwallet2 -k -m 23

coldwallet2 -k -m 33 -i ./data/pubkey/xxx.csv
coldwallet2 -k -m 34 -i ./data/pubkey/xxx.csv

coldwallet2 -k -m 50
coldwallet2 -k -m 51

coldwallet2 -k -m 60
coldwallet2 -k -m 61

#coldwallet1
coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv
coldwallet1 -k -m 42 -i ./data/pubkey/xxx.csv

coldwallet1 -k -m 50
coldwallet1 -k -m 51

#watch only wallet
wallet -k -m 1 -x -i ./data/pubkey/xxx.csv
wallet -k -m 2 -x -i ./data/pubkey/xxx.csv
wallet -k -m 3 -x -i ./data/pubkey/xxx.csv

wallet -d -m 1
