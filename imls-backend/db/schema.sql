SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: admin; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA admin;


--
-- Name: api; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA api;


--
-- Name: imlswifi; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA imlswifi;


--
-- Name: bin_devices_per_hour(date, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bin_devices_per_hour(_day date, _fscs_id text) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
DECLARE
    _start TIMESTAMPTZ;
    _end TIMESTAMPTZ;
    _count INT;
    _hour INT := 0;
    _day_end INT := 24;
    num_devices_arr INT[];
    _timezone_offset INT;
BEGIN
    SELECT api.get_timezone_from_fscs_id(_fscs_id) INTO _timezone_offset;
    _hour := _hour - _timezone_offset;
    _day_end := _day_end - _timezone_offset;

    -- Hardcoded EDT for now. Will add the look up table next to pass in the time zone
    WHILE _hour < _day_end LOOP

        -- Casting the DATE variable to a TIMESTAMP to add it to the interval
        _start = _day::TIMESTAMP + make_interval(hours=> _hour);
        _end =  _day + make_interval(hours=> _hour, mins => 59, secs => 59);

        -- This select stores the result in the variable _count.
        SELECT count(*) INTO _count
        FROM api.presences
        WHERE  fscs_id = _fscs_id
        AND (presences.start_time::TIMESTAMPTZ < presences.end_time::TIMESTAMPTZ)
        AND (presences.start_time::TIMESTAMPTZ <= _end::TIMESTAMPTZ)
        AND (presences.end_time > _start::TIMESTAMPTZ);
        num_devices_arr := array_append(num_devices_arr, _count);

        _hour := _hour + 1;
    END LOOP;
    RETURN num_devices_arr;

END
$$;


--
-- Name: get_timezone_from_fscs_id(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.get_timezone_from_fscs_id(_fscs_id text) RETURNS integer
    LANGUAGE plpgsql IMMUTABLE
    AS $$
DECLARE
    _timezone TIMETZ;
    _timezone_offset INT:=0;
BEGIN
    SELECT imls_lookup.timezone::TIMETZ INTO _timezone::TIMETZ
    FROM api.imls_lookup
    WHERE imls_lookup.fscs_id = _fscs_id;

    _timezone_offset := extract(timezone_hour FROM _timezone::TIMETZ);
    SELECT extract(timezone_hour FROM _timezone::TIMETZ) INTO _timezone_offset;

    RETURN _timezone_offset;
END

$$;
-- Name: test(); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.test() RETURNS TABLE(start_time timestamp with time zone, end_time timestamp with time zone)
    LANGUAGE plpgsql
    AS $$
BEGIN
RETURN QUERY
SELECT presences.start_time, presences.end_time
FROM api.presences;
END; $$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: libraries; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.libraries (
    fscs_id character varying(16) NOT NULL
);


--
-- Name: libraries; Type: VIEW; Schema: admin; Owner: -
--

CREATE VIEW admin.libraries AS
 SELECT libraries.fscs_id
   FROM imlswifi.libraries;


--
-- Name: helo; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.helo (
    uid integer NOT NULL,
    message character varying(42)
);


--
-- Name: imls_lookup; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.imls_lookup (
    id integer NOT NULL,
    fscs_id character varying(16) NOT NULL,
    timezone time with time zone NOT NULL
);


--
-- Name: imls_lookup; Type: VIEW; Schema: api; Owner: -
--

CREATE VIEW api.imls_lookup AS
 SELECT imls_lookup.id,
    imls_lookup.fscs_id,
    imls_lookup.timezone
   FROM imlswifi.imls_lookup;


--
-- Name: presences; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.presences (
    presence_id integer NOT NULL,
    start_time timestamp with time zone NOT NULL,
    end_time timestamp with time zone NOT NULL,
    fscs_id character varying(16) NOT NULL,
    sensor_id integer NOT NULL,
    manufacturer_index integer
);


--
-- Name: presences; Type: VIEW; Schema: api; Owner: -
--

CREATE VIEW api.presences AS
 SELECT presences.presence_id,
    presences.start_time,
    presences.end_time,
    presences.fscs_id,
    presences.sensor_id,
    presences.manufacturer_index
   FROM imlswifi.presences;


--
-- Name: imls_lookup_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.imls_lookup_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: imls_lookup_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.imls_lookup_id_seq OWNED BY imlswifi.imls_lookup.id;


--
-- Name: libraries; Type: TABLE; Schema: imlswifi; Owner: -
-- Name: heartbeats; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.heartbeats (
    heartbeat_id integer NOT NULL,
    fscs_id character varying(16) NOT NULL,
    sensor_id integer NOT NULL,
    ping_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    sensor_serial character varying(32) NOT NULL,
    sensor_version character varying(16) NOT NULL
);


--
-- Name: heartbeats_heartbeat_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.heartbeats_heartbeat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: heartbeats_heartbeat_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.heartbeats_heartbeat_id_seq OWNED BY imlswifi.heartbeats.heartbeat_id;


--
-- Name: heartbeats_sensor_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.heartbeats_sensor_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: heartbeats_sensor_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.heartbeats_sensor_id_seq OWNED BY imlswifi.heartbeats.sensor_id;


--
-- Name: presences_presence_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.presences_presence_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: presences_presence_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.presences_presence_id_seq OWNED BY imlswifi.presences.presence_id;


--
-- Name: presences_sensor_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.presences_sensor_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: presences_sensor_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.presences_sensor_id_seq OWNED BY imlswifi.presences.sensor_id;


--
-- Name: sensors; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.sensors (
    sensor_id integer NOT NULL,
    fscs_id character varying(16) NOT NULL
);


--
-- Name: sensors_sensor_id_seq; Type: SEQUENCE; Schema: imlswifi; Owner: -
--

CREATE SEQUENCE imlswifi.sensors_sensor_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sensors_sensor_id_seq; Type: SEQUENCE OWNED BY; Schema: imlswifi; Owner: -
--

ALTER SEQUENCE imlswifi.sensors_sensor_id_seq OWNED BY imlswifi.sensors.sensor_id;


--
-- Name: durations_v2; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.durations_v2 (
    id integer NOT NULL,
    pi_serial character varying(16),
    fcfs_seq_id character varying(16),
    device_tag character varying(32),
    session_id character varying(255),
    patron_index integer,
    manufacturer_index integer,
    start text,
    "end" text
);


--
-- Name: durations_v2_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.durations_v2_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: durations_v2_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.durations_v2_id_seq OWNED BY public.durations_v2.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: imls_lookup id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup ALTER COLUMN id SET DEFAULT nextval('imlswifi.imls_lookup_id_seq'::regclass);
-- Name: heartbeats heartbeat_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats ALTER COLUMN heartbeat_id SET DEFAULT nextval('imlswifi.heartbeats_heartbeat_id_seq'::regclass);


--
-- Name: heartbeats sensor_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats ALTER COLUMN sensor_id SET DEFAULT nextval('imlswifi.heartbeats_sensor_id_seq'::regclass);


--
-- Name: presences presence_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.presences ALTER COLUMN presence_id SET DEFAULT nextval('imlswifi.presences_presence_id_seq'::regclass);


--
-- Name: presences sensor_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.presences ALTER COLUMN sensor_id SET DEFAULT nextval('imlswifi.presences_sensor_id_seq'::regclass);


--
-- Name: sensors sensor_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.sensors ALTER COLUMN sensor_id SET DEFAULT nextval('imlswifi.sensors_sensor_id_seq'::regclass);


--
-- Name: durations_v2 id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.durations_v2 ALTER COLUMN id SET DEFAULT nextval('public.durations_v2_id_seq'::regclass);


--
-- Name: helo helo_pkey; Type: CONSTRAINT; Schema: api; Owner: -
--

ALTER TABLE ONLY api.helo
    ADD CONSTRAINT helo_pkey PRIMARY KEY (uid);


--
-- Name: imls_lookup imls_lookup_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup
    ADD CONSTRAINT imls_lookup_pkey PRIMARY KEY (id);
-- Name: heartbeats heartbeats_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats
    ADD CONSTRAINT heartbeats_pkey PRIMARY KEY (heartbeat_id);


--
-- Name: libraries libraries_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.libraries
    ADD CONSTRAINT libraries_pkey PRIMARY KEY (fscs_id);


--
-- Name: presences presences_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.presences
    ADD CONSTRAINT presences_pkey PRIMARY KEY (presence_id);


--
-- Name: sensors sensors_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.sensors
    ADD CONSTRAINT sensors_pkey PRIMARY KEY (sensor_id);


--
-- Name: durations_v2 durations_v2_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.durations_v2
    ADD CONSTRAINT durations_v2_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: fk_heartbeat_library_index; Type: INDEX; Schema: imlswifi; Owner: -
--

CREATE INDEX fk_heartbeat_library_index ON imlswifi.heartbeats USING btree (fscs_id);


--
-- Name: fk_heartbeat_sensor_index; Type: INDEX; Schema: imlswifi; Owner: -
--

CREATE INDEX fk_heartbeat_sensor_index ON imlswifi.heartbeats USING btree (sensor_id);


--
-- Name: fk_presence_library_index; Type: INDEX; Schema: imlswifi; Owner: -
--

CREATE INDEX fk_presence_library_index ON imlswifi.presences USING btree (fscs_id);


--
-- Name: fk_presence_sensor_index; Type: INDEX; Schema: imlswifi; Owner: -
--

CREATE INDEX fk_presence_sensor_index ON imlswifi.presences USING btree (sensor_id);


--
-- Name: fk_sensor_library_index; Type: INDEX; Schema: imlswifi; Owner: -
--

CREATE INDEX fk_sensor_library_index ON imlswifi.sensors USING btree (fscs_id);


--
-- Name: imls_lookup fk_lookup_library; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup
    ADD CONSTRAINT fk_lookup_library FOREIGN KEY (fscs_id) REFERENCES imlswifi.libraries(fscs_id);
-- Name: heartbeats fk_heartbeat_library; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats
    ADD CONSTRAINT fk_heartbeat_library FOREIGN KEY (fscs_id) REFERENCES imlswifi.libraries(fscs_id);


--
-- Name: heartbeats fk_heartbeat_sensor; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats
    ADD CONSTRAINT fk_heartbeat_sensor FOREIGN KEY (sensor_id) REFERENCES imlswifi.sensors(sensor_id);


--
-- Name: presences fk_presence_library; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.presences
    ADD CONSTRAINT fk_presence_library FOREIGN KEY (fscs_id) REFERENCES imlswifi.libraries(fscs_id);


--
-- Name: presences fk_presence_sensor; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.presences
    ADD CONSTRAINT fk_presence_sensor FOREIGN KEY (sensor_id) REFERENCES imlswifi.sensors(sensor_id);


--
-- Name: sensors fk_sensor_library; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.sensors
    ADD CONSTRAINT fk_sensor_library FOREIGN KEY (fscs_id) REFERENCES imlswifi.libraries(fscs_id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20220727163656'),
    ('20220729132839'),
    ('20220729150547'),
    ('20220811192329'),
    ('20220811194424'),
    ('20220811195049'),
    ('20220817173135'),
    ('20220818150144'),
    ('20220818154959');
    ('20220822125647');
