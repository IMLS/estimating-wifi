-- migrate:up

CREATE TABLE imlswifi.imls_lookup (
    id SERIAL PRIMARY KEY,
    fscs_id VARCHAR(16) NOT NULL,
    timezone TIMETZ NOT NULL,
    CONSTRAINT fk_lookup_library
        FOREIGN KEY(fscs_id)
            REFERENCES imlswifi.libraries(fscs_id)
);

INSERT INTO imlswifi.imls_lookup (fscs_id, timezone) (
    SELECT DISTINCT fscs_id,'00:00:00-04'::TIMETZ FROM imlswifi.libraries
);

CREATE VIEW api.imls_lookup AS SELECT * FROM imlswifi.imls_lookup;
GRANT SELECT ON TABLE api.imls_lookup TO web_anon;

-- migrate:down

DROP VIEW api.imls_lookup;
DROP TABLE imlswifi.imls_lookup;