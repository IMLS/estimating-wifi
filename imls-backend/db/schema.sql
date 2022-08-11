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
-- Name: api; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA api;


--
-- Name: imlswifi; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA imlswifi;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: helo; Type: TABLE; Schema: api; Owner: -
--

CREATE TABLE api.helo (
    uid integer NOT NULL,
    message character varying(42)
);


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
-- Name: libraries; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.libraries (
    fscs_id character varying(16) NOT NULL
);


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
    sensor_serial character varying(32) NOT NULL,
    sensor_version character varying(16) NOT NULL,
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
    ('20220811194424');
