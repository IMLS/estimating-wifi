-- migrate:up

DROP VIEW api.presences;

ALTER TABLE imlswifi.heartbeats
ADD COLUMN sensor_serial VARCHAR(32) NOT NULL,
ADD COLUMN sensor_version VARCHAR(16) NOT NULL;

ALTER TABLE imlswifi.sensors
DROP COLUMN sensor_serial CASCADE,
DROP COLUMN sensor_version CASCADE;

CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
GRANT SELECT ON TABLE api.presences TO web_anon;

ALTER TABLE imlswifi.heartbeats
RENAME COLUMN hourly_ping TO ping_time;

-- migrate:down

ALTER TABLE imlswifi.heartbeats
RENAME COLUMN ping_time TO hourly_ping;

ALTER TABLE imlswifi.sensors
ADD COLUMN sensor_serial VARCHAR(32) NOT NULL,
ADD COLUMN sensor_version VARCHAR(16) NOT NULL;

ALTER TABLE imlswifi.heartbeats
DROP COLUMN sensor_serial CASCADE,
DROP COLUMN sensor_version CASCADE;
