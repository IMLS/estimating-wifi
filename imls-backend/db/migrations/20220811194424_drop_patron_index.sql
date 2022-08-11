-- migrate:up

DROP VIEW api.presences;
ALTER TABLE imlswifi.presences DROP COLUMN patron_index;
CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
GRANT SELECT ON TABLE api.presences TO web_anon;

-- migrate:down

ALTER TABLE imlswifi.presences ADD COLUMN patron_index INTEGER;
