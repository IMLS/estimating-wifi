SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE IF NOT EXISTS
imlswifi.libraries (
    fscs_id character varying(16) PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS
imlswifi.imls_lookup (
    id SERIAL PRIMARY KEY,
    fscs_id character varying(16) NOT NULL REFERENCES imlswifi.libraries(fscs_id),
    timezone time with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS
imlswifi.presences (
    presence_id SERIAL PRIMARY KEY,
    start_time timestamp with time zone NOT NULL,
    end_time timestamp with time zone NOT NULL,
    fscs_id character varying(16) NOT NULL REFERENCES imlswifi.libraries(fscs_id),
    -- FIXME: what about our tag provided by the library IT director?
    -- We're missing something, but it should not be a UID of some sort.
    manufacturer_index integer
);

CREATE TABLE IF NOT EXISTS
imlswifi.heartbeats (
    heartbeat_id SERIAL PRIMARY KEY,
    fscs_id character varying(16) NOT NULL REFERENCES imlswifi.libraries(fscs_id),
    ping_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    sensor_serial character varying(32) NOT NULL,
    sensor_version character varying(16) NOT NULL
);

CREATE TABLE IF NOT EXISTS
imlswifi.sensors (
    sensor_id SERIAL PRIMARY KEY,
    fscs_id character varying(16) NOT NULL REFERENCES imlswifi.libraries(fscs_id),
    labels character varying,
    api_key character varying
    -- jwt character varying
);

CREATE TABLE IF NOT EXISTS
 public.schema_migrations (
    version character varying(255) PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS
basic_auth.users (
    fscs_id   text primary key check ( fscs_id ~* '^[A-Z][A-Z][0-9]{4}-[0-9]{3}$' ),
    api_key   text not null check (length(api_key) < 512),
    role      name not null check (length(role) < 512)
);
