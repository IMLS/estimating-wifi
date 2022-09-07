-- migrate:up
CREATE ROLE authenticator noinherit LOGIN PASSWORD '<fill_password>';
GRANT web_anon TO authenticator;


CREATE ROLE sensor nologin;
GRANT sensor TO authenticator;

GRANT usage ON SCHEMA api TO sensor;
GRANT usage ON SCHEMA imlswifi TO sensor;
GRANT SELECT, INSERT ON imlswifi.heartbeats TO sensor;
GRANT SELECT, INSERT ON imlswifi.presences TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.heartbeats_heartbeat_id_seq TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.presences_presence_id_seq TO sensor;

-- migrate:down
REVOKE usage on schema api to sensor;
REVOKE usage on schema imlswifi to sensor;
REVOKE SELECT, INSERT on imlswifi.heartbeats to sensor;
REVOKE SELECT, INSERT on imlswifi.presences to sensor;
REVOKE usage, select on sequence imlswifi.heartbeats_heartbeat_id_seq to sensor;
REVOKE usage, select on sequence imlswifi.presences_presence_id_seq to sensor;

REVOKE web_anon to authenticator;

DROP ROLE sensor;
DROP ROLE authenticator;
