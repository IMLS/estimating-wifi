-- rabbit_raw:
-- store all API calls (before validation)
DROP TABLE IF EXISTS public.rabbit_raw_v2;
CREATE TABLE public.rabbit_raw_v2 (
    id serial PRIMARY KEY,
    date_created timestamp with time zone DEFAULT current_timestamp,
    collection character varying(255),
    data json,
    content_type character varying(255)
);

-- rabbit_review:
-- should an API call fail validation, store the object and validation error.
DROP TABLE IF EXISTS public.rabbit_review_v2;
CREATE TABLE public.rabbit_review_v2 (
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
DROP TABLE IF EXISTS public.validators_v2;
CREATE TABLE public.validators_v2 (
    date_created timestamp with time zone DEFAULT current_timestamp,
    name character varying(255) NOT NULL,
    validator json NOT NULL
);
ALTER TABLE ONLY public.validators_v2
    ADD CONSTRAINT validators_v2_pkey PRIMARY KEY (name);

-- counts:
-- store wifi summaries
DROP TABLE IF EXISTS public.counts_v2;
CREATE TABLE public.counts_v2 (
    id serial PRIMARY KEY,
    pi_serial character varying(16),
    fcfs_seq_id character varying(16),
    device_tag character varying(32),
    session_id character varying(255),
    minimum_minutes integer,
    maximum_minutes integer,
    patron_count integer,
    patron_minutes integer,
    device_count integer,
    device_minutes integer,
    transient_count integer,
    transient_minutes integer
);

-- events:
-- store events from the Raspberry Pi session collector.
DROP TABLE IF EXISTS public.events_v2;
CREATE TABLE public.events_v2 (
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
