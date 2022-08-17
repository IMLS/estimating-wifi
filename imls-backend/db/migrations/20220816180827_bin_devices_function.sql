-- migrate:up
-- migrate:up
-- the 'day' variable needs to be passed in as such: '2022-05-01 00:00:00+00'

CREATE OR REPLACE FUNCTION api.bin_devices_per_hour(_day TIMESTAMPTZ, _fscs_id TEXT) RETURNS TABLE (num_devices INT, hour TIMESTAMPTZ)
    AS $$
-- Variable declarations happen in a single block.
-- Tables cannot be declared as variables.
DECLARE 
    _start TIMESTAMPTZ;
    _end TIMESTAMPTZ;
    _count INT;
    -- _hour INT := 0;
    -- _hour TIMESTAMPTZ:=_day;
    _period TIMESTAMPTZ;
BEGIN
    -- Temporary tables are not:
    --   * Scoped to the function.
    --   * Automatically cleaned up after a function execution.
    CREATE TEMP TABLE _results (hour TIMESTAMPTZ, count INT);
    _period := _day::TIMESTAMPTZ + '1 day'::INTERVAL;    
    WHILE _day < _period LOOP
        -- Note that type casting is often necessary to make things work.
        -- Here, I had to pass in a DATE object, but to add it to an interval,
        -- it had to become a TIMESTAMP, which sets the HH:MM:SS to 00:00:00.
        _start = _day;
        _end =  _day + interval '1 hour';
        -- This select stores the result in the variable _count.
        SELECT count(*) INTO _count
        FROM imlswifi.presences
        WHERE  fscs_id = _fscs_id 
        AND (presences.start_time::TIMESTAMPTZ < presences.end_time::TIMESTAMPTZ)
        AND (presences.start_time::TIMESTAMPTZ <= _end::TIMESTAMPTZ)
        AND (presences.end_time > _start::TIMESTAMPTZ);
        INSERT INTO _results VALUES (_day, _count);
        -- This was a handy way to "printf" the loop counter.
        -- RAISE NOTICE 'HOUR: %', _hour;
        -- If you forget this, the loop will not terminate.
        -- _hour := _hour + 1;
        _day := _day + interval '1 hour';
    END LOOP;
    -- Specify the return query.
    RETURN QUERY SELECT * FROM _results;
    -- Drop the temporary table.
    DROP TABLE _results;
    -- Return the two-column table.
    RETURN;
END
$$ LANGUAGE plpgsql;

-- migrate:down
DROP FUNCTION api.bin_devices;
