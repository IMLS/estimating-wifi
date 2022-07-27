CREATE SCHEMA imlswifi;

CREATE TABLE imlswifi.libraries (
    fscs_id VARCHAR PRIMARY KEY
    -- Unsure of other fields the libraries table needs
);

CREATE TABLE imlswifi.sensors (
    sensor_id SERIAL PRIMARY KEY,
    sensor_serial VARCHAR,
    sensor_version VARCHAR,
    heartbeat BOOL,
    CONSTRAINT fk_sensor_library
        FOREIGN KEY(fscs_id)
            REFERENCES imlswifi.libraries(fscs_id)
);

CREATE TABLE imlswifi.presences (
    presence_id SERIAL PRIMARY KEY,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    CONSTRAINT fk_presence_library
        FOREIGN KEY (fscs_id)
            REFERENCES imlswifi.libraries(fscs_id),
    CONSTRAINT fk_presence_sensor
        FOREIGN KEY (sensor_id)
            REFERENCES imlswifi.sensors(sensor_id)
);

CREATE INDEX fk_sensor_library_index ON imlswifi.sensors(fscs_id);

CREATE INDEX fk_presence_library_index ON imlswifi.presences(fscs_id);

CREATE INDEX fk_presence_sensor_index ON imlswifi.presences(sensor_id);