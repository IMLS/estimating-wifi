-- migrate:up

CREATE TABLE imlswifi.heartbeats (
    heartbeat_id SERIAL PRIMARY KEY,
    fscs_id VARCHAR(16) NOT NULL,
    sensor_id SERIAL NOT NULL,
    hourly_ping TIMESTAMP,
    CONSTRAINT fk_heartbeat_library
        FOREIGN KEY (fscs_id)
            REFERENCES imlswifi.libraries(fscs_id),
    CONSTRAINT fk_heartbeat_sensor
        FOREIGN KEY (sensor_id)
            REFERENCES imlswifi.sensors(sensor_id)

);

CREATE INDEX fk_heartbeat_library_index ON imlswifi.heartbeats(fscs_id);

CREATE INDEX fk_heartbeat_sensor_index ON imlswifi.heartbeats(sensor_id);

INSERT INTO imlswifi.heartbeats(sensor_id, fscs_id)
    SELECT presences.sensor_id, presences.fscs_id
     FROM imlswifi.presences;

-- migrate:down

DROP TABLE imlswifi.heartbeats;
