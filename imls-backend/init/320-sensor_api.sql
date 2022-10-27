
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
