-- migrate:up

CREATE OR REPLACE FUNCTION api.get_timezone_from_fscs_id(_fscs_id TEXT) RETURNS INT
    AS $$
DECLARE
    _timezone TIMETZ;
    _time_offset INT:=0;
    _fscs_id TEXT;
BEGIN
    SELECT imls_lookup.timezone::TIMETZ INTO _timezone::TIMETZ
    FROM api.imls_lookup
    WHERE imls_lookup.fscs_id = _fscs_id;

    RAISE NOTICE 'VAlUE: %', _timezone;
    
    -- _time_offset := extract(timezone_hour FROM _timezone::TIMETZ);
    SELECT extract(timezone_hour FROM _timezone::TIMETZ) INTO _time_offset;

    RAISE NOTICE 'VALUE: %', _time_offset;

    RETURN _time_offset;
END

$$ LANGUAGE plpgsql IMMUTABLE;


-- migrate:down
DROP FUNCTION api.get_timezone_from_fscs_id;
