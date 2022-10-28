-- Always drop first. If you change parameters/return types,
-- you cannot necessarily replace the function.
DROP FUNCTION IF EXISTS api.beat_the_heart;
CREATE OR REPLACE FUNCTION api.beat_the_heart(
	_fscs_id character varying,
	_sensor_serial character varying,
	_sensor_version character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
INSERT INTO imlswifi.heartbeats(fscs_id, sensor_serial, sensor_version)
   VALUES(_fscs_id, _sensor_serial, _sensor_version);
   RETURN '{"result":"OK"}'::json;
END;
$BODY$;

DROP FUNCTION IF EXISTS api.get_library_timezone;
CREATE OR REPLACE FUNCTION api.get_library_timezone(
    fscs_id CHARACTER VARYING)
    RETURNS json
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
    _tz TIME WITH TIME ZONE;
BEGIN
    SELECT timezone into _tz FROM data.timezone_lookup
    WHERE timezone_lookup.fscs_id = get_library_timezone.fscs_id;
    RETURN json_build_object('time', to_json(_tz));
END;
$BODY$;

NOTIFY pgrst, 'reload schema';
