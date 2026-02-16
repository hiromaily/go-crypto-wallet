#!/bin/sh

set -eu

docker compose exec wallet-mysql mysql -u root -proot keygen -e "$(cat ./sql/ganache_key.sql)"
