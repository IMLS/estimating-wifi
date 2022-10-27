#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" <<-EOSQL
    ALTER DATABASE imls SET "app.jwt_secret" TO "${POSTGRES_JWT_SECRET}";
    INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor');
    NOTIFY pgrst, 'reload schema'
EOSQL