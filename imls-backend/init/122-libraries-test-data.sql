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

-- Note we're giving sensors random ids in the range 1 to 1000.
INSERT INTO imlswifi.sensors(sensor_id, fscs_id)
    (SELECT floor(random() * 1000 + 1)::int AS sensor_id, fcfs_seq_id AS fscs_id
     FROM durations_v2
     GROUP BY pi_serial, fcfs_seq_id);


-- CREATE TABLE imlswifi.presences (
--     presence_id integer NOT NULL,
--     start_time timestamp with time zone NOT NULL,
--     end_time timestamp with time zone NOT NULL,
--     fscs_id character varying(16) NOT NULL,
--     sensor_id integer NOT NULL,
--     manufacturer_index integer
-- );

INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, sensor_id, manufacturer_index)
    (SELECT to_timestamp(start::int) AS start_time, to_timestamp("end"::int) AS end_time, fcfs_seq_id AS fscs_id, s.sensor_id,
        manufacturer_index AS manufacturer_index
     FROM durations_v2 d, imlswifi.sensors s);