#!/bin/bash

source "${BASH_SOURCE%/*}/../.env"

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor') ON CONFLICT DO NOTHING;
EOSQL

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 -f test-data-load.sql

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    NOTIFY pgrst, 'reload schema';
EOSQL
