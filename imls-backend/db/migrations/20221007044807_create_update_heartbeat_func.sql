-- migrate:up
-- FUNCTION: api.update_hb(character varying, integer, timestamp with time zone, character varying, character varying)

-- DROP FUNCTION IF EXISTS api.update_hb(character varying, integer, timestamp with time zone, character varying, character varying);

CREATE OR REPLACE FUNCTION api.update_hb(
	_fscs character varying,
	_sensor integer,
	_hb timestamp with time zone,
	_serial character varying,
	_version character varying)
    RETURNS character varying
    LANGUAGE 'plpgsql'
AS $BODY$
begin
INSERT INTO imlswifi.heartbeats(fscs_id, sensor_id, ping_time, sensor_serial, sensor_version)
   VALUES(_fscs, _sensor, _hb, _serial, _version);
   RETURN _sensor;
end;
$BODY$;

ALTER FUNCTION api.update_hb(character varying, integer, timestamp with time zone, character varying, character varying)
    OWNER TO postgres;

GRANT EXECUTE ON FUNCTION api.update_hb(character varying, integer, timestamp with time zone, character varying, character varying) TO postgres;

REVOKE ALL ON FUNCTION api.update_hb(character varying, integer, timestamp with time zone, character varying, character varying) FROM PUBLIC;


-- migrate:down

DROP FUNCTION api.update_hb;