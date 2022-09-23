-- migrate:up
CREATE OR REPLACE FUNCTION api.lib_search_name(
	_name text)
    RETURNS JSON
    LANGUAGE 'sql'
    COST 100
AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE libname LIKE '%'|| UPPER(_name) || '%') AS X;
$$;

ALTER FUNCTION api.lib_search_name
    OWNER TO postgres;

-- migrate:down
DROP FUNCTION api.lib_search_name;

