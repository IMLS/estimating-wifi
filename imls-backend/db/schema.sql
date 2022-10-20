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
-- Name: data; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA data;


--
-- Name: imlswifi; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA imlswifi;


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: pgjwt; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgjwt WITH SCHEMA public;


--
-- Name: EXTENSION pgjwt; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgjwt IS 'JSON Web Token API for Postgresql';


--
-- Name: jwt_token; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.jwt_token AS (
	token text
);


--
-- Name: bin_devices_over_time(date, text, boolean, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bin_devices_over_time(_start date, _fscs_id text, _direction boolean, _days integer) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
DECLARE
	_new_start DATE;
	_cnt INTEGER;
	_full INTEGER[][]= '{{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}}';
	--_full INTEGER[][];
	_day_return INTEGER[];
BEGIN
    _cnt := 0;
	_new_start := _start;
	WHILE _cnt < _days LOOP
		IF _cnt != 0 THEN
			IF _direction THEN
				_new_start := _new_start::date + 1;
			ELSE
				_new_start := _new_start::date - 1;
			END IF;
		END IF;

		raise notice 'Value: %', _new_start;

		SELECT api.bin_devices_per_hour(_new_start, _fscs_id) INTO _day_return;

		_full := array_cat(_full, _day_return);

	    _cnt := _cnt + 1;

    END LOOP;
	SELECT (_full)[2:_cnt +1] INTO _full;
    RETURN _full;

END
$$;


--
-- Name: bin_devices_per_hour(date, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.bin_devices_per_hour(_start date, _fscs_id text) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
DECLARE
    _init_start TIMESTAMPTZ;
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
        _init_start = _start::TIMESTAMP + make_interval(hours=> _hour);
        _end =  _start + make_interval(hours=> _hour, mins => 59, secs => 59);

        -- This select stores the result in the variable _count.
        SELECT count(*) INTO _count
        FROM api.presences
        WHERE  fscs_id = _fscs_id
        AND (presences.start_time::TIMESTAMPTZ < presences.end_time::TIMESTAMPTZ)
        AND (presences.start_time::TIMESTAMPTZ <= _end::TIMESTAMPTZ)
        AND (presences.end_time > _init_start::TIMESTAMPTZ);
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


--
-- Name: jwt_gen(text, text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.jwt_gen(s_key text, s_role text) RETURNS public.jwt_token
    LANGUAGE sql
    AS $$
  SELECT public.sign(
    row_to_json(r), s_key
  ) AS token
  FROM (
    SELECT
      s_role as role,
      extract(epoch from now())::integer + 315360000  AS exp
  ) r;
$$;


--
-- Name: lib_search_fscs(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_fscs(_fscs_id text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT row_to_json(X) FROM
(SELECT *, CONCAT(fscskey,'-',TO_CHAR(fscs_seq, 'fm000')) AS fscs_id
FROM data.imls_data) AS X
WHERE X.fscs_id = _fscs_id;
$$;


--
-- Name: lib_search_name(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_name(_name text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE libname LIKE '%'|| UPPER(_name) || '%') AS X;
$$;


--
-- Name: lib_search_state(text); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.lib_search_state(_state_code text) RETURNS json
    LANGUAGE sql
    AS $$
SELECT json_agg(X) FROM
(SELECT *  FROM data.imls_data WHERE stabr LIKE UPPER(_state_code) || '%') AS X;
$$;


--
-- Name: sensor_info(integer, character varying); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.sensor_info(_sensor integer, _install_key character varying) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
declare
_jwt varchar;
begin
SELECT jwt FROM imlswifi.sensors WHERE sensor_id = _sensor AND install_key = _install_key INTO _jwt;
   RETURN _jwt;
end;
$$;


--
-- Name: sensor_setup(character varying, character varying, character varying); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.sensor_setup(_fscs character varying, _label character varying, _install_key character varying) RETURNS integer
    LANGUAGE plpgsql
    AS $$
declare
_jwt varchar;
_sensor integer;
begin
SELECT api.jwt_gen(current_setting('app.jwt_secret'), 'sensor') INTO _jwt;
INSERT INTO imlswifi.sensors(fscs_id, labels, install_key, jwt)
   VALUES(_fscs, _label, _install_key, _jwt);
SELECT currval(pg_get_serial_sequence('imlswifi.sensors','sensor_id')) INTO _sensor;
   RETURN _sensor;
end;
$$;


--
-- Name: update_hb(character varying, integer, timestamp with time zone, character varying, character varying); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.update_hb(_fscs character varying, _sensor integer, _hb timestamp with time zone, _serial character varying, _version character varying) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
begin
INSERT INTO imlswifi.heartbeats(fscs_id, sensor_id, ping_time, sensor_serial, sensor_version)
   VALUES(_fscs, _sensor, _hb, _serial, _version);
   RETURN _sensor;
end;
$$;


--
-- Name: update_presence(timestamp with time zone, timestamp with time zone, character varying, integer, integer); Type: FUNCTION; Schema: api; Owner: -
--

CREATE FUNCTION api.update_presence(_start timestamp with time zone, _end timestamp with time zone, _fscs character varying, _sensor integer, _manufacture integer) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
begin
INSERT INTO imlswifi.presences(start_time, end_time, fscs_id, sensor_id, manufacturer_index)
   VALUES(_start, _end, _fscs, _sensor, _manufacture);
   RETURN _sensor;
end;
$$;


--
-- Name: pgrst_watch(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.pgrst_watch() RETURNS event_trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NOTIFY pgrst, 'reload schema';
END;
$$;


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
-- Name: imls_data; Type: TABLE; Schema: data; Owner: -
--

CREATE TABLE data.imls_data (
    stabr character(2),
    fscskey character(6),
    fscs_seq integer,
    c_fscs character(1),
    libid character varying(64),
    libname character varying(256),
    address character varying(256),
    city character varying(32),
    zip character(5),
    zip4 character(4),
    cnty character varying(64),
    phone character(10),
    c_out_ty character(2),
    sq_feet integer,
    f_sq_ft character(4),
    l_num_bm integer,
    hours integer,
    f_hours character(4),
    wks_open integer,
    f_wksopn character(4),
    yr_sub integer,
    obereg integer,
    statstru integer,
    statname integer,
    stataddr integer,
    longitud double precision,
    latitude double precision,
    incitsst integer,
    incitsco integer,
    gnisplac character varying(6),
    cntypop integer,
    locale character varying(2),
    centract double precision,
    cenblock integer,
    cdcode integer,
    cbsa integer,
    microf character(1),
    geostatus character(1),
    geoscore double precision,
    geomtype character varying(32),
    c19wkscl integer,
    c19wkslo integer
);


--
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
    fscs_id character varying(16) NOT NULL,
    labels character varying,
    install_key character varying,
    jwt character varying
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
-- Name: heartbeats heartbeat_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats ALTER COLUMN heartbeat_id SET DEFAULT nextval('imlswifi.heartbeats_heartbeat_id_seq'::regclass);


--
-- Name: heartbeats sensor_id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats ALTER COLUMN sensor_id SET DEFAULT nextval('imlswifi.heartbeats_sensor_id_seq'::regclass);


--
-- Name: imls_lookup id; Type: DEFAULT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup ALTER COLUMN id SET DEFAULT nextval('imlswifi.imls_lookup_id_seq'::regclass);


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
-- Name: heartbeats heartbeats_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.heartbeats
    ADD CONSTRAINT heartbeats_pkey PRIMARY KEY (heartbeat_id);


--
-- Name: imls_lookup imls_lookup_pkey; Type: CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup
    ADD CONSTRAINT imls_lookup_pkey PRIMARY KEY (id);


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
-- Name: imls_lookup fk_lookup_library; Type: FK CONSTRAINT; Schema: imlswifi; Owner: -
--

ALTER TABLE ONLY imlswifi.imls_lookup
    ADD CONSTRAINT fk_lookup_library FOREIGN KEY (fscs_id) REFERENCES imlswifi.libraries(fscs_id);


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
-- Name: pgrst_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER pgrst_watch ON ddl_command_end
   EXECUTE FUNCTION public.pgrst_watch();


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
    ('20220818154959'),
    ('20220822125647'),
    ('20220831042321'),
    ('20220902170318'),
    ('20220907165121'),
    ('20220912214016'),
    ('20220923222643'),
    ('20220923223659'),
    ('20220923230122'),
    ('20221007044807'),
    ('20221007045130'),
    ('20221007045430'),
    ('20221007045511'),
    ('20221017220651'),
    ('20221017221845'),
    ('20221018044654'),
    ('20221020173024');
