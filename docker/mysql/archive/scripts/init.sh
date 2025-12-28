#!/bin/bash

# copy proper sql files to docker-entrypoint-initdb.d from /sqls/*.sql
# it may be too late to sync
cp /sqls/user.sql /docker-entrypoint-initdb.d/
if [ "${ENV}" = 'watch' ]; then
  cp /sqls/watch.sql /docker-entrypoint-initdb.d/
elif [ "${ENV}" = 'keygen' ]; then
  cp /sqls/keygen.sql /docker-entrypoint-initdb.d/
else
  cp /sqls/sign.sql /docker-entrypoint-initdb.d/
fi
