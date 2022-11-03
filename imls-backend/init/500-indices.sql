CREATE INDEX fk_sensor_library_index ON imlswifi.sensors(fscs_id);

CREATE INDEX fk_presence_library_index ON imlswifi.presences(fscs_id);

CREATE INDEX fk_presence_sensor_index ON imlswifi.presences(presence_id);

CREATE INDEX fk_heartbeat_library_index ON imlswifi.heartbeats(fscs_id);

CREATE INDEX fk_heartbeat_sensor_index ON imlswifi.heartbeats(heartbeat_id);

-- For performance on binning queries. 
-- We think we can do better. The queries currently have loops, which cannot
-- be optimized well. However, they work. These indices provide reasonable front-end
-- performance, until we can revisit the DB structure. Most of what the queries are doing
-- could be pre-computed at the moment that data is submitted from the sensor. (FWIW,
-- we could pre-compute this information *on the sensor* before sending.)
create index presences_start_end_index on imlswifi.presences (start_time, end_time);
create index presences_id_start_end_index on imlswifi.presences (fscs_id, start_time, end_time);