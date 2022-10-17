-- migrate:up
-- FUNCTION: api.sensor_setup(character varying, character varying, character varying)

-- DROP FUNCTION IF EXISTS api.sensor_setup(character varying, character varying, character varying);

CREATE OR REPLACE FUNCTION api.sensor_setup(
	_fscs character varying,
	_label character varying,
	_install_key character varying)
    RETURNS integer
    LANGUAGE 'plpgsql'

AS $BODY$
declare 
_jwt varchar;
_sensor integer;
begin
SELECT api.jwt_gen(current_setting('app.jwt_secret'), 'sensor') INTO _jwt;
INSERT INTO imlswifi.sensors(fscs_id, labels, install_key, jwt)
   VALUES(_fscs, _label, _install_key, _jwt);
SELECT currval(pg_get_serial_sequence('imlswifi.sensors','sensor_id')) INTO _sensor;
   RETURN _sensor;
end;
$BODY$;

ALTER FUNCTION api.sensor_setup(character varying, character varying, character varying)
    OWNER TO postgres;

GRANT EXECUTE ON FUNCTION api.sensor_setup(character varying, character varying, character varying) TO postgres;

GRANT EXECUTE ON FUNCTION api.sensor_setup(character varying, character varying, character varying) TO sensor;

-- migrate:down

DROP FUNCTION api.sensor_setup;