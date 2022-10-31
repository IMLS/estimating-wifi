
CREATE OR REPLACE FUNCTION api.beat_the_heart(
	_sensor_serial character varying,
	_sensor_version character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
   _validated_id character varying := current_setting('request.jwt.claims', true)::json->>'fscs_id';
BEGIN
INSERT INTO imlswifi.heartbeats(fscs_id, sensor_serial, sensor_version)
   VALUES(_validated_id, _sensor_serial, _sensor_version);
   RETURN '{"result":"OK"}'::json;
END;
$BODY$;

-- TODO: validate fscs based on the api key, instead of relying on senders to send their fscs id
CREATE OR REPLACE FUNCTION api.update_presence(
       _start TIMESTAMPTZ,
       _end TIMESTAMPTZ)
    RETURNS CHARACTER varying
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
   _validated_id CHARACTER varying := current_setting('request.jwt.claims', true)::json->>'fscs_id';
BEGIN
INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, manufacturer_index)
   VALUES(_start, _end, _validated_id, 0);
   RETURN _validated_id;
END;
$BODY$;

CREATE OR REPLACE FUNCTION api.verify_presence(
   _fscs_id CHARACTER varying(16),
   _start TIMESTAMPTZ,
   _end TIMESTAMPTZ)
   RETURNS INTEGER
   LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
   _result INTEGER;
BEGIN 
   SELECT presence_id INTO _result
   FROM api.presences
   WHERE 
      (presences.fscs_id = _fscs_id) AND
      (presences.start_time = _start) AND
      (presences.end_time = _end);
   RETURN _result;
END;
$BODY$;