-- migrate:up

INSERT INTO imlswifi.libraries
    (SELECT DISTINCT(fcfs_seq_id) FROM durations_v2);

INSERT INTO imlswifi.sensors(sensor_serial, fscs_id, sensor_version)
    (SELECT DISTINCT(pi_serial) AS sensor_serial, fcfs_seq_id AS fscs_id, '3.0' AS sensor_version
     FROM durations_v2
     GROUP BY pi_serial, fcfs_seq_id);

INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, sensor_id)
    (SELECT to_timestamp(start::int) AS start_time, to_timestamp("end"::int) AS end_time, fcfs_seq_id AS fscs_id, s.sensor_id
     FROM durations_v2 d
     LEFT JOIN imlswifi.sensors s
     ON d.pi_serial = s.sensor_serial AND
        s.sensor_version = '3.0');

-- migrate:down

DELETE FROM imlswifi.presences
    WHERE start_time >= timestamp '2022-05-01' AND start_time < timestamp '2022-06-01';
DELETE FROM imlswifi.sensors
    WHERE sensor_version = '3.0';
DELETE FROM imlswifi.libraries
    WHERE fscs_id like 'AA%';
