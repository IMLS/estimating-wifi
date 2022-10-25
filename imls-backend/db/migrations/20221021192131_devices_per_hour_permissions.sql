-- migrate:up
GRANT EXECUTE ON FUNCTION api.bin_devices_per_hour(DATE, TEXT) TO web_anon;


-- migrate:down

REVOKE ALL ON api.bin_devices_per_hour FROM web_anon