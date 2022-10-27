#!/bin/bash

set -e

# The image this is embedded in runs after the DB comes to life.
# Then, we can inject things into the DB at runtime.


echo --------------------------
echo LIVE CONFIG SQL SCRIPTS...
echo --------------------------

for FILE in `ls /runtime`; do   
    
    echo CHECKING $FILE

    if [[ $FILE == *.sql ]]
    then
        echo --------------------------
        echo SQL :: $FILE
        echo --------------------------
        psql ${POSTGRES_DB_URI} -a -f /runtime/${FILE}
    fi
done

echo --------------------------
echo LIVE CONFIG ONE-OFFS...
echo --------------------------

psql ${POSTGRES_DB_URI} -v ON_ERROR_STOP=1 <<-EOSQL
    ALTER DATABASE imls SET "app.jwt_secret" TO "${POSTGRES_JWT_SECRET}";
    INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor');
    NOTIFY pgrst, 'reload schema'
EOSQL
