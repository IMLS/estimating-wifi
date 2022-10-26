
CREATE OR REPLACE FUNCTION api.beat_the_heart(
	_fscs_id character varying,
	_sensor_id integer,
	_sensor_serial character varying,
	_sensor_version character varying)
    RETURNS character varying
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
INSERT INTO imlswifi.heartbeats(fscs_id, sensor_id, sensor_serial, sensor_version)
   VALUES(_fscs_id, _sensor_id,  _sensor_serial, _sensor_version);
   RETURN _sensor_id;
END;
$BODY$;
