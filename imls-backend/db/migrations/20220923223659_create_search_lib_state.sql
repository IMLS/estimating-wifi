-- migrate:up
CREATE OR REPLACE FUNCTION api.lib_search_state(
	_state_code text)
    RETURNS JSON
    LANGUAGE 'sql'
    COST 100
AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE stabr LIKE UPPER(_state_code) || '%') AS X;
$$;

ALTER FUNCTION api.lib_search_state
    OWNER TO postgres;

-- migrate:down
DROP FUNCTION api.lib_search_state;
