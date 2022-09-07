-- migrate:up

CREATE OR REPLACE FUNCTION api.bin_devices_per_hour(_day DATE, _fscs_id TEXT) RETURNS INT[]
    AS $$
DECLARE 
    _start TIMESTAMPTZ;
    _end TIMESTAMPTZ;
    _count INT;
    _hour INT := 0;
    _day_end INT := 24;
    num_devices_arr INT[];
    _timezone_offset INT;
BEGIN
    SELECT api.get_timezone_from_fscs_id(_fscs_id) INTO _timezone_offset;
    _hour := _hour - _timezone_offset;
    _day_end := _day_end - _timezone_offset;

    -- Hardcoded EDT for now. Will add the look up table next to pass in the time zone
    WHILE _hour < _day_end LOOP

        -- Casting the DATE variable to a TIMESTAMP to add it to the interval
        _start = _day::TIMESTAMP + make_interval(hours=> _hour);
        _end =  _day + make_interval(hours=> _hour, mins => 59, secs => 59);

        -- This select stores the result in the variable _count.
        SELECT count(*) INTO _count
        FROM api.presences
        WHERE  fscs_id = _fscs_id 
        AND (presences.start_time::TIMESTAMPTZ < presences.end_time::TIMESTAMPTZ)
        AND (presences.start_time::TIMESTAMPTZ <= _end::TIMESTAMPTZ)
        AND (presences.end_time > _start::TIMESTAMPTZ);
        num_devices_arr := array_append(num_devices_arr, _count);

        _hour := _hour + 1;
    END LOOP;
    RETURN num_devices_arr;

END
$$ LANGUAGE plpgsql;

-- migrate:down
DROP FUNCTION api.bin_devices_per_hour;