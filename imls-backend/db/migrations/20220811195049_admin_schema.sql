-- migrate:up

CREATE SCHEMA admin;

CREATE VIEW admin.sensors AS SELECT * FROM imlswifi.sensors;
CREATE VIEW admin.libraries AS SELECT * FROM imlswifi.libraries;

-- migrate:down

DROP VIEW admin.libraries;
DROP VIEW admin.sensors;
DROP SCHEMA admin;
