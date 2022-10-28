#!/bin/bash

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor') ON CONFLICT DO NOTHING;
EOSQL

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 -f test-data.sql

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO data.libraries
        (SELECT DISTINCT(fcfs_seq_id) FROM durations_v2);

    INSERT INTO data.libraries(fscs_id) values ('KY0069-002');

    INSERT INTO imlswifi.sensors(fscs_id)
        (SELECT DISTINCT fcfs_seq_id AS fscs_id 
        FROM durations_v2);

    INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, manufacturer_index)
        (SELECT to_timestamp(start::int) AS start_time, to_timestamp("end"::int) AS end_time, fcfs_seq_id AS fscs_id,
            manufacturer_index AS manufacturer_index
        FROM durations_v2 d, imlswifi.sensors s);
EOSQL

psql ${DATABASE_URL} -v ON_ERROR_STOP=0 <<-EOSQL
    INSERT INTO data.timezone_lookup (fscs_id, timezone) (
        SELECT DISTINCT fscs_id,'00:00:00-04'::TIMETZ FROM data.libraries
    );
EOSQL
