-- migrate:up

--
-- Name: libraries; Type: VIEW; Schema: admin; Owner: -
--

CREATE VIEW admin.libraries AS
 SELECT libraries.fscs_id
   FROM imlswifi.libraries;


--
-- Name: timezone_lookup; Type: VIEW; Schema: api; Owner: -
--
-- NOTE: What is the difference, really, between these two?
CREATE VIEW api.timezone_lookup AS SELECT * FROM imlswifi.timezone_lookup;
-- CREATE VIEW api.timezone_lookup AS
--  SELECT timezone_lookup.id,
--     timezone_lookup.fscs_id,
--     timezone_lookup.timezone
--    FROM imlswifi.timezone_lookup;


--
-- Name: presences; Type: VIEW; Schema: api; Owner: -
--

-- CREATE VIEW api.presences AS
--  SELECT presences.presence_id,
--     presences.start_time,
--     presences.end_time,
--     presences.fscs_id,
--     presences.sensor_id,
--     presences.manufacturer_index
--    FROM imlswifi.presences;
CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
