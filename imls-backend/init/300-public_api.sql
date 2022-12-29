-- migrate:up

--
-- Name: bin_devices_over_time(date, text, boolean, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bin_devices_over_time(_start date, _fscs_id text, _direction boolean, _days integer) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
DECLARE
	_new_start DATE;
	_cnt INTEGER;
	_full INTEGER[][]= '{{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}';
	_day_return INTEGER[];
BEGIN
    _cnt := 0;
	_new_start := _start;
	WHILE _cnt < _days LOOP
		IF _cnt != 0 THEN
			IF _direction THEN
				_new_start := _new_start::date + 1;
			ELSE
				_new_start := _new_start::date - 1;
			END IF;
		END IF;

		SELECT api.bin_devices_per_hour(_new_start, _fscs_id) INTO _day_return;

		_full := array_cat(_full, _day_return);

	    _cnt := _cnt + 1;

    END LOOP;
	SELECT (_full)[2:_cnt +1] INTO _full;
    RETURN _full;

END
$$;


--
-- Name: bin_devices_per_hour(date, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bin_devices_per_hour(_start date, _fscs_id text) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
DECLARE
    _init_start TIMESTAMPTZ;
    _end TIMESTAMPTZ;
    _count INT;
    _hour INT := 0;
    _day_end INT := 24;
    num_devices_arr INT[];
    _timezone_offset INT;
    -- FIXME: These are hard-coded in.
    -- We could pass them in as parameters, for future flexibility.
    _min_minutes INT := 5;
    _max_minutes INT := 600;
BEGIN
    SELECT api.get_timezone_from_fscs_id(_fscs_id) INTO _timezone_offset;
    _hour := _hour - _timezone_offset;
    _day_end := _day_end - _timezone_offset;

    WHILE _hour < _day_end LOOP

        -- Casting the DATE variable to a TIMESTAMP to add it to the interval
        _init_start = _start::TIMESTAMP + make_interval(hours=> _hour);
        _end =  _start + make_interval(hours=> _hour, mins => 59, secs => 59);

        -- This select stores the result in the variable _count.
        SELECT count(*) INTO _count
        FROM api.presences
        WHERE  fscs_id = _fscs_id
        AND (presences.start_time::TIMESTAMPTZ < presences.end_time::TIMESTAMPTZ)
        AND (presences.start_time::TIMESTAMPTZ <= _end::TIMESTAMPTZ)
        AND (presences.end_time > _init_start::TIMESTAMPTZ)
        AND EXTRACT(EPOCH FROM (presences.end_time::TIMESTAMPTZ - presences.start_time::TIMESTAMPTZ))/60 >= _min_minutes
        AND EXTRACT(EPOCH FROM (presences.end_time::TIMESTAMPTZ - presences.start_time::TIMESTAMPTZ))/60 <= _max_minutes;

        num_devices_arr := array_append(num_devices_arr, _count);

        _hour := _hour + 1;
    END LOOP;
    RETURN num_devices_arr;

END
$$;

-- migrate:down


--
-- Name: lib_search_fscs(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_fscs(_fscs_id text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT row_to_json(X) FROM
(SELECT *, CONCAT(fscskey,'-',TO_CHAR(fscs_seq, 'fm000')) AS fscs_id
FROM data.imls_data) AS X
WHERE X.fscs_id = _fscs_id;
$$;


--
-- Name: lib_search_name(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_name(_name text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE libname LIKE '%'|| UPPER(_name) || '%') AS X;
$$;


--
-- Name: lib_search_state(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_state(_state_code text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE stabr LIKE UPPER(_state_code) || '%') AS X;
$$;

-- migrate:down