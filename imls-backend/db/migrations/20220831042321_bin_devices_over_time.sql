-- migrate:up
-- FUNCTION: api.bin_devices_over_time(date, text, boolean, integer)

-- DROP FUNCTION IF EXISTS api.bin_devices_over_time(date, text, boolean, integer);

CREATE OR REPLACE FUNCTION api.bin_devices_over_time(
	_start date,
	_fscs_id text,
	_direction boolean,
	_days integer)
    RETURNS integer[]
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
DECLARE 
	_new_start DATE;
	_cnt INTEGER;
	_full INTEGER[][]= '{{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}';
	_day_return INTEGER[];
BEGIN
    _cnt := 0;
	_new_start := _start;
	WHILE _cnt < _days LOOP
		IF _direction THEN 
			_new_start := _new_start::date + _cnt;
		ELSE 
			_new_start := _new_start::date - _cnt;
		END IF;
		
		raise notice 'Value: %', _new_start;
		
		SELECT api.bin_devices_per_hour(_new_start, _fscs_id) INTO _day_return;
	
		_full := array_cat(_full, _day_return);
		
	    _cnt := _cnt + 1;

    END LOOP;
	SELECT (_full)[2:_cnt +1] INTO _full;
    RETURN _full;

END
$BODY$;

ALTER FUNCTION api.bin_devices_over_time(date, text, boolean, integer)
    OWNER TO postgres;

-- Permissions used for testing

GRANT EXECUTE ON FUNCTION api.bin_devices_over_time(date, text, boolean, integer) TO PUBLIC;

GRANT EXECUTE ON FUNCTION api.bin_devices_over_time(date, text, boolean, integer) TO postgres;

GRANT EXECUTE ON FUNCTION api.bin_devices_over_time(date, text, boolean, integer) TO web_anon;



-- migrate:down

DROP FUNCTION api.bin_devices_over_time;