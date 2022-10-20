-- migrate:up
GRANT EXECUTE ON FUNCTION api.lib_search_fscs TO web_anon;
GRANT EXECUTE ON FUNCTION api.lib_search_name TO web_anon;
GRANT EXECUTE ON FUNCTION api.lib_search_state TO web_anon;
GRANT USAGE ON SCHEMA data TO web_anon;
GRANT SELECT ON data.imls_data TO web_anon;
-- migrate:down
REVOKE ALL ON FUNCTION api.lib_search_fscs FROM web_anon;
REVOKE ALL ON FUNCTION api.lib_search_name FROM web_anon;
REVOKE ALL ON FUNCTION api.lib_search_state FROM web_anon;
REVOKE ALL ON SCHEMA data FROM web_anon;
REVOKE ALL ON data.imls_data FROM web_anon;

