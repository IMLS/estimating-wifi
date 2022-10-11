-- migrate:up
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pgjwt;

-- Inferred from
-- https://github.com/PostgREST/postgrest-docs/issues/280
CREATE TYPE jwt_token AS (
  token TEXT
);

CREATE OR REPLACE FUNCTION api.jwt_gen(
	s_key text,
	s_role text)
    RETURNS jwt_token
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
  SELECT public.sign(
    row_to_json(r), s_key
  ) AS token
  FROM (
    SELECT
      s_role as role,
      extract(epoch from now())::integer + 315360000  AS exp
  ) r;
$BODY$;

ALTER FUNCTION api.jwt_gen(text, text)
    OWNER TO postgres;


-- migrate:down
DROP EXTENSION IF EXISTS pgjwt;
DROP EXTENSION IF EXISTS pgcrypto;


DROP FUNCTION api.jwt_gen;
