-- migrate:up
-- FUNCTION: api.jwt_gen(text, text)

-- DROP FUNCTION IF EXISTS api.jwt_gen(text, text);

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

DROP FUNCTION api.jwt_gen;
