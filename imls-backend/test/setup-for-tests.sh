#!/bin/bash

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor') ON CONFLICT DO NOTHING;
EOSQL

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 -f test-data.sql

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO imlswifi.timezone_lookup (fscs_id, timezone) (
        SELECT DISTINCT fscs_id,'00:00:00-04'::TIMETZ FROM imlswifi.libraries
    );
EOSQL
