-- migrate:up
-- FUNCTION: api.sensor_info(integer, character varying)

-- DROP FUNCTION IF EXISTS api.sensor_info(integer, character varying);

CREATE OR REPLACE FUNCTION api.sensor_info(
	_sensor integer,
	_install_key character varying)
    RETURNS character varying
    LANGUAGE 'plpgsql'
    
AS $BODY$
declare 
_jwt varchar;
begin
SELECT jwt FROM imlswifi.sensors WHERE sensor_id = _sensor AND install_key = _install_key INTO _jwt;
   RETURN _jwt;
end;
$BODY$;

ALTER FUNCTION api.sensor_info(integer, character varying)
    OWNER TO postgres;

GRANT EXECUTE ON FUNCTION api.sensor_info(integer, character varying) TO postgres;

GRANT EXECUTE ON FUNCTION api.sensor_info(integer, character varying) TO sensor;

-- migrate:down

DROP FUNCTION api.sensor_info;