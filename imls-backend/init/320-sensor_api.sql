
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

-- TODO: validate fscs based on the api key, instead of relying on senders to send their fscs id
CREATE OR REPLACE FUNCTION api.update_presence(
       _start timestamptz,
       _end timestamptz,
       _fscs character varying (16))
    RETURNS character varying
    LANGUAGE 'plpgsql'
AS $BODY$
begin
INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, manufacturer_index)
   VALUES(_start, _end, _fscs, 0);
   RETURN _fscs;
end;
$BODY$;
