
GRANT USAGE ON SCHEMA api TO web_anon;
GRANT USAGE ON SCHEMA data TO web_anon;
GRANT SELECT ON TABLE data.imls_data TO web_anon;
GRANT SELECT ON TABLE api.presences TO web_anon;

-- Public 
GRANT sensor TO authenticator;

GRANT web_anon TO authenticator;

-- Authentication
GRANT EXECUTE ON FUNCTION api.login(text, text) to web_anon;
GRANT EXECUTE ON FUNCTION api.verify_presence(character varying, timestamptz, timestamptz) to web_anon;

-- GRANT EXECUTE ON FUNCTION api.jwt_gen(text, text) TO web_anon;
-- GRANT EXECUTE ON FUNCTION api.sensor_setup(character varying, character varying, character varying) TO sensor;
-- GRANT EXECUTE ON FUNCTION api.sensor_info(integer, character varying) TO sensor;

-- Private

GRANT EXECUTE ON FUNCTION api.beat_the_heart(character varying, character varying) TO sensor;
GRANT EXECUTE ON FUNCTION api.update_presence(character varying, character varying) TO sensor;
GRANT SELECT, INSERT ON imlswifi.heartbeats TO sensor;
GRANT USAGE ON SCHEMA api TO sensor;
GRANT USAGE ON SCHEMA imlswifi TO sensor;
GRANT USAGE ON SCHEMA basic_auth TO sensor;
GRANT SELECT, INSERT ON imlswifi.sensors to sensor;
GRANT SELECT, INSERT ON basic_auth.users TO sensor;
GRANT SELECT, INSERT ON imlswifi.presences TO sensor;
GRANT SELECT, INSERT ON imlswifi.sensors to sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.sensors_sensor_id_seq TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.heartbeats_heartbeat_id_seq TO sensor;
GRANT usage, SELECT ON SEQUENCE imlswifi.presences_presence_id_seq TO sensor;
GRANT EXECUTE ON FUNCTION api.sensor_setup TO postgres;
GRANT EXECUTE ON FUNCTION api.sensor_setup TO sensor;



