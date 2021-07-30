-- rabbit_raw:
-- store all API calls (before validation)
DROP TABLE IF EXISTS public.rabbit_raw_v1;
CREATE TABLE public.rabbit_raw_v1 (
    id serial PRIMARY KEY,
    date_created timestamp with time zone DEFAULT current_timestamp,
    collection character varying(255),
    data json,
    content_type character varying(255)
);

-- rabbit_review:
-- should an API call fail validation, store the object and validation error.
DROP TABLE IF EXISTS public.rabbit_review_v1;
CREATE TABLE public.rabbit_review_v1 (
    id serial PRIMARY KEY,
    date_created timestamp with time zone DEFAULT current_timestamp,
    headers json,
    whole_table_errors json,
    rows json,
    valid_row_count integer,
    invalid_row_count integer
);

-- validators:
-- store GoodTables validators for API calls.
DROP TABLE IF EXISTS public.validators_v1;
CREATE TABLE public.validators_v1 (
    date_created timestamp with time zone DEFAULT current_timestamp,
    name character varying(255) NOT NULL,
    validator json NOT NULL
);
ALTER TABLE ONLY public.validators_v1
    ADD CONSTRAINT validators_v1_pkey PRIMARY KEY (name);

-- wifi:
-- store incoming wifi data from the Raspberry Pi session collector.
DROP TABLE IF EXISTS public.wifi_v1;
CREATE TABLE public.wifi_v1 (
    id serial PRIMARY KEY,
    event_id integer,
    pi_serial character varying(16),
    fcfs_seq_id character varying(16),
    device_tag character varying(32),
    "localtime" timestamp with time zone,
    servertime timestamp with time zone DEFAULT current_timestamp,
    -- session_id: the unique identifier for session, randomly generated when session-counter starts
    session_id character varying(255),
    manufacturer_index integer,
    patron_index integer
);

-- events:
-- store events from the Raspberry Pi session collector.
DROP TABLE IF EXISTS public.events_v1;
CREATE TABLE public.events_v1 (
    id serial PRIMARY KEY,
    pi_serial character varying(16),
    fcfs_seq_id character varying(16),
    device_tag character varying(32),
    session_id character varying(255),
    "localtime" timestamp with time zone,
    servertime timestamp with time zone DEFAULT current_timestamp,
    tag character varying(255),
    info text
);
