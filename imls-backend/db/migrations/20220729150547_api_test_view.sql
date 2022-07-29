-- migrate:up
CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
GRANT SELECT ON TABLE api.presences TO web_anon;

-- migrate:down
DROP VIEW api.presences;

