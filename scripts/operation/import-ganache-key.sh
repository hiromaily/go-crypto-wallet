#!/bin/sh

set -eu

docker compose exec wallet-db mysql -u root -proot keygen -e "$(cat ./sql/ganache_key.sql)"
