-- For testing, we need some timezone data in the imls_lookup table.
-- We do that by selecting libraries from the libraries table, and pulling it in.
INSERT INTO imlswifi.imls_lookup (fscs_id, timezone) (
    SELECT DISTINCT fscs_id,'00:00:00-04'::TIMETZ FROM imlswifi.libraries
);

