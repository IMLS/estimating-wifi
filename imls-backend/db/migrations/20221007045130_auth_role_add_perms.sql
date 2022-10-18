-- migrate:up
CREATE ROLE authenticator noinherit NOLOGIN;
GRANT web_anon TO authenticator;


CREATE ROLE sensor nologin;
GRANT sensor TO authenticator;

GRANT USAGE ON SCHEMA api TO sensor;
GRANT USAGE ON SCHEMA imlswifi TO sensor;
-- GRANT EXECUTE ON FUNCTION api.update_hb TO sensor;
-- GRANT EXECUTE ON FUNCTION api.update_presence TO sensor;
GRANT SELECT, INSERT ON imlswifi.heartbeats TO sensor;
GRANT SELECT, INSERT ON imlswifi.presences TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.heartbeats_heartbeat_id_seq TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.presences_presence_id_seq TO sensor;


-- migrate:down
REVOKE ALL ON SCHEMA api FROM sensor;
REVOKE ALL ON SCHEMA imlswifi FROM sensor;
REVOKE ALL ON imlswifi.heartbeats FROM sensor;
REVOKE ALL ON imlswifi.presences FROM sensor;
-- REVOKE ALL ON FUNCTION api.update_hb FROM sensor;
-- REVOKE ALL ON FUNCTION api.update_presence FROM sensor;
REVOKE ALL ON SEQUENCE imlswifi.heartbeats_heartbeat_id_seq FROM sensor;
REVOKE ALL ON SEQUENCE imlswifi.presences_presence_id_seq FROM sensor;

REVOKE web_anon FROM authenticator;

DROP ROLE sensor;
DROP ROLE authenticator;

