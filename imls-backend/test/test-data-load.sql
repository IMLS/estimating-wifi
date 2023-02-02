CREATE TABLE IF NOT EXISTS
 public.durations_v2 (
    id INTEGER PRIMARY KEY,
    pi_serial character varying(16),
    fcfs_seq_id character varying(16),
    device_tag character varying(32),
    session_id character varying(255),
    patron_index integer,
    start text,
    "end" text,
    timezone character varying(6)
);

\COPY public.durations_v2 (id, pi_serial, fcfs_seq_id, device_tag, session_id, patron_index, start, "end", timezone) FROM 'durations_v2.csv' csv;

INSERT INTO imlswifi.libraries
    (SELECT DISTINCT(fcfs_seq_id) FROM durations_v2);

INSERT INTO imlswifi.libraries(fscs_id) values ('KY0069-002');

INSERT INTO imlswifi.sensors(fscs_id)
    (SELECT DISTINCT fcfs_seq_id AS fscs_id 
     FROM durations_v2);

INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, timezone)
    (SELECT to_timestamp(start::int) AS start_time, to_timestamp("end"::int) AS end_time, fcfs_seq_id AS fscs_id, 
    timezone AS timezone
     FROM durations_v2 d, imlswifi.sensors s);
