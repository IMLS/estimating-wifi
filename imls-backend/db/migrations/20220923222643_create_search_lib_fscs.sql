-- migrate:up
CREATE OR REPLACE FUNCTION api.lib_search_fscs(
	_fscs_id text)
    RETURNS JSON
    LANGUAGE 'sql'
    COST 100
AS $$
SELECT row_to_json(X) FROM
(SELECT *, CONCAT(fscskey,'-',TO_CHAR(fscs_seq, 'fm000')) AS fscs_id 
FROM data.imls_data) AS X 
WHERE X.fscs_id = _fscs_id;
$$;

ALTER FUNCTION api.lib_search_fscs
    OWNER TO postgres;

-- migrate:down
DROP FUNCTION api.lib_search_fscs;
