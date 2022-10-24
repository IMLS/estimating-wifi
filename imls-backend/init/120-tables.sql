-- migrate:up

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: libraries; Type: TABLE; Schema: imlswifi; Owner: -
--

CREATE TABLE imlswifi.libraries (
    fscs_id character varying(16) NOT NULL
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
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


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

-- migrate:down
