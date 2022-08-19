-- migrate:up

CREATE OR REPLACE FUNCTION api.get_timezone_from_fscs_id(_fscs_id TEXT) RETURNS INT
    AS $$
DECLARE
    _timezone TIMETZ;
    _timezone_offset INT:=0;
BEGIN
    SELECT imls_lookup.timezone::TIMETZ INTO _timezone::TIMETZ
    FROM api.imls_lookup
    WHERE imls_lookup.fscs_id = _fscs_id;
    
    _timezone_offset := extract(timezone_hour FROM _timezone::TIMETZ);
    SELECT extract(timezone_hour FROM _timezone::TIMETZ) INTO _timezone_offset;

    RETURN _timezone_offset;
END

$$ LANGUAGE plpgsql IMMUTABLE;


-- migrate:down
DROP FUNCTION api.get_timezone_from_fscs_id;
