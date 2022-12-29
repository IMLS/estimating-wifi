-- migrate:up

--
-- Name: libraries; Type: VIEW; Schema: admin; Owner: -
--

CREATE VIEW admin.libraries AS
 SELECT libraries.fscs_id
   FROM imlswifi.libraries;

--
-- Name: presences; Type: VIEW; Schema: api; Owner: -
--

CREATE VIEW api.presences AS SELECT * FROM imlswifi.presences;
