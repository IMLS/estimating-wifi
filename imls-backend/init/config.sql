CREATE SCHEMA api;

CREATE ROLE web_anon NOLOGIN;

GRANT USAGE ON SCHEMA api TO web_anon;

CREATE TABLE api.helo (
    uid INT PRIMARY KEY,
    message VARCHAR(42)
)