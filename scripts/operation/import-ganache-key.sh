#!/bin/sh

set -eu

docker compose exec keygen-db mysql -u root -proot -e "$(cat ./sql/ganache_key.sql)"
