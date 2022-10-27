-- (id, pi_serial, fcfs_seq_id, device_tag, session_id, patron_index, manufacturer_index, start, "end")


-- CREATE TABLE imlswifi.sensors (
--     sensor_id integer NOT NULL,
--     fscs_id character varying(16) NOT NULL,
--     labels character varying,
--     install_key character varying,
--     jwt character varying
-- );

INSERT INTO imlswifi.libraries
    (SELECT DISTINCT(fcfs_seq_id) FROM durations_v2);

-- FIXME -- This is part of the *test* data.
-- Note we're giving sensors random ids in the range 1 to 1M.
INSERT INTO imlswifi.sensors(fscs_id)
    (SELECT DISTINCT fcfs_seq_id AS fscs_id 
     FROM durations_v2);


-- CREATE TABLE imlswifi.presences (
--     presence_id integer NOT NULL,
--     start_time timestamp with time zone NOT NULL,
--     end_time timestamp with time zone NOT NULL,
--     fscs_id character varying(16) NOT NULL,
--     sensor_id integer NOT NULL,
--     manufacturer_index integer
-- );

INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, manufacturer_index)
    (SELECT to_timestamp(start::int) AS start_time, to_timestamp("end"::int) AS end_time, fcfs_seq_id AS fscs_id,
        manufacturer_index AS manufacturer_index
     FROM durations_v2 d, imlswifi.sensors s);
