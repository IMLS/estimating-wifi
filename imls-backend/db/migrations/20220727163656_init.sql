-- migrate:up

CREATE SCHEMA imlswifi;

CREATE TABLE imlswifi.libraries (
    fscs_id VARCHAR(16) PRIMARY KEY
);

CREATE TABLE imlswifi.sensors (
    sensor_id SERIAL PRIMARY KEY,
    sensor_serial VARCHAR(32) NOT NULL,
    sensor_version VARCHAR(16) NOT NULL,
    fscs_id VARCHAR(16) NOT NULL,
    CONSTRAINT fk_sensor_library
        FOREIGN KEY(fscs_id)
            REFERENCES imlswifi.libraries(fscs_id)
);

CREATE TABLE imlswifi.presences (
    presence_id SERIAL PRIMARY KEY,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    fscs_id VARCHAR(16) NOT NULL,
    sensor_id SERIAL,
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

-- migrate:down

DROP TABLE imlswifi.presences;
DROP TABLE imlswifi.sensors;
DROP TABLE imlswifi.libraries;
DROP SCHEMA imlswifi;
