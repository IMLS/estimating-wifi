CREATE INDEX fk_sensor_library_index ON imlswifi.sensors(fscs_id);

CREATE INDEX fk_presence_library_index ON imlswifi.presences(fscs_id);

CREATE INDEX fk_presence_sensor_index ON imlswifi.presences(presence_id);

CREATE INDEX fk_heartbeat_library_index ON imlswifi.heartbeats(fscs_id);

CREATE INDEX fk_heartbeat_sensor_index ON imlswifi.heartbeats(heartbeat_id);
