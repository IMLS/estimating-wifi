
-- Anonymous endpoints

-- Schemas
GRANT USAGE ON SCHEMA api TO web_anon;
GRANT USAGE ON SCHEMA data TO web_anon;

GRANT sensor TO authenticator;
GRANT web_anon TO authenticator;

-- Tables
GRANT SELECT ON TABLE api.presences TO web_anon;
GRANT SELECT ON data.timezone_lookup TO web_anon;

GRANT EXECUTE ON FUNCTION api.login(text, text) to web_anon;
GRANT EXECUTE ON FUNCTION api.get_library_timezone TO web_anon;


-- Authenticated endpoints
GRANT EXECUTE ON FUNCTION api.beat_the_heart(character varying, character varying, character varying) TO sensor;
GRANT SELECT, INSERT ON imlswifi.heartbeats TO sensor;
GRANT USAGE ON SCHEMA api TO sensor;
GRANT USAGE ON SCHEMA imlswifi TO sensor;
GRANT SELECT, INSERT ON imlswifi.presences TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.heartbeats_heartbeat_id_seq TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.presences_presence_id_seq TO sensor;
