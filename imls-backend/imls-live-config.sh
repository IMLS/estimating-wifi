#!/bin/bash

set -e

# The image this is embedded in runs after the DB comes to life.
# Then, we can inject things into the DB at runtime.

echo --- POSTGRES URI ---
echo ${POSTGRES_DB_URI}
echo --------------------

for FILE in `ls /runtime`; do   
    
    echo CHECKING $FILE

    if [[ $FILE == *.sql ]]
    then
        echo --------------------------
        echo SQL :: $FILE
        echo --------------------------
        psql ${POSTGRES_DB_URI} -a -f /runtime/${FILE}
    fi
    if [[ $FILE == *.sh ]]
    then
        echo --------------------------
        echo SHELL :: $FILE
        echo --------------------------
        source /runtime/${FILE}
    fi
done
