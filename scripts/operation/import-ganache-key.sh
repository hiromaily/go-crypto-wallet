#!/bin/sh

set -eu

docker compose exec btc-keygen-db mysql -u root -proot  -e "$(cat ./sql/ganache_key.sql)"