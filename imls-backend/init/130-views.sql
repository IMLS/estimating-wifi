-- migrate:up

--
-- Name: libraries; Type: VIEW; Schema: admin; Owner: -
--

CREATE VIEW admin.libraries AS
 SELECT libraries.fscs_id
   FROM imlswifi.libraries;


--
-- Name: imls_lookup; Type: VIEW; Schema: api; Owner: -
--
-- NOTE: What is the difference, really, between these two?
CREATE VIEW api.imls_lookup AS SELECT * FROM imlswifi.imls_lookup;
-- CREATE VIEW api.imls_lookup AS
--  SELECT imls_lookup.id,
--     imls_lookup.fscs_id,
--     imls_lookup.timezone
--    FROM imlswifi.imls_lookup;


--
-- Name: presences; Type: VIEW; Schema: api; Owner: -
--

CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
