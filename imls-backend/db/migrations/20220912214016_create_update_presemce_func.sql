-- migrate:up
-- FUNCTION: api.update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer)

CREATE OR REPLACE FUNCTION api.update_presence(
	_start timestamptz,
	_end timestamptz,
	_fscs character varying (16),
	_sensor integer,
	_manufacture integer)
    RETURNS character varying
    LANGUAGE 'plpgsql'
AS $BODY$
begin
INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, sensor_id, manufacturer_index)
   VALUES(_start, _end, _fscs, _sensor, _manufacture);
   RETURN _sensor;
end;
$BODY$;

ALTER FUNCTION api.update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer)
    OWNER TO postgres;

GRANT EXECUTE ON FUNCTION api.update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer) TO postgres;

GRANT EXECUTE ON FUNCTION api.update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer) TO sensor;

REVOKE ALL ON FUNCTION api.update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer) FROM PUBLIC;

-- migrate:down
DROP FUNCTION api.update_presence;
