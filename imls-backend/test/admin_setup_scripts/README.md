# Admin Tool

## Setup 

### ENV file
This script assumes an .env file in the same folder with 

SERVER_ADDR=http://0.0.0.0:Port
ADMIN_USERNAME=
ADMIN_PASSWORD=

### SQL Function and Permission Requirements 
- Must be run before the tool can be used

```
//— Library Setup

CREATE OR REPLACE FUNCTION api.library_setup(
	_fscs_id character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
IF EXISTS(SELECT 1 FROM imlswifi.libraries WHERE fscs_id = _fscs_id) THEN
	INSERT INTO imlswifi.libraries(fscs_id) 
		VALUES (_fscs_id) ON CONFLICT (fscs_id) DO NOTHING;
	   RETURN '{"result":"Library Inserted"}'::json;
ELSE
	RETURN '{"result":"Library Inserted"}'::json;
END IF;
END;
$BODY$;

//— Library Delete

CREATE OR REPLACE FUNCTION api.library_remove(
	_fscs_id character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
DELETE FROM imlswifi.libraries WHERE fscs_id = _fscs_id;
   RETURN '{"result":"Library Removed"}'::json;
END;
$BODY$;

//—Sensor Setup

CREATE OR REPLACE FUNCTION api.sensor_setup(
	_sensor character varying,
	_key character varying,
	_label character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
BEGIN
IF EXISTS(SELECT 1 FROM imlswifi.sensors WHERE fscs_id = _sensor) THEN
	   RETURN '{"result":"Sensor Exists"}'::json;
ELSE
	INSERT INTO basic_auth.users 
		VALUES (_sensor, _key, 'sensor') ON CONFLICT DO NOTHING;
	INSERT INTO imlswifi.sensors(fscs_id, labels)
	   VALUES(_sensor, _label) ON CONFLICT DO NOTHING;
	   RETURN '{"result":"Sensor Inserted"}'::json;
END IF;
END;
$BODY$;

//—Sensor Remove

CREATE OR REPLACE FUNCTION api.sensor_remove(
	_sensor character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
BEGIN
IF EXISTS(SELECT 1 FROM imlswifi.sensors WHERE fscs_id = _sensor) THEN
	DELETE FROM basic_auth.users 
		WHERE fscs_id = _sensor;
	DELETE FROM imlswifi.sensors
	   WHERE fscs_id = _sensor;
	   RETURN '{"result":"Sensor Deleted"}'::json;
ELSE
	RETURN '{"result":"Sensor Does Not Exist"}'::json;
END IF;
END;
$BODY$;

//—Update Password

CREATE OR REPLACE FUNCTION api.update_password(
	_sensor character varying,
	_key character varying)
    RETURNS json
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
BEGIN
UPDATE basic_auth.users SET api_key = _key WHERE fscs_id = _sensor; 
   RETURN '{"result":"Password Updated"}'::json;
END;
$BODY$;


CREATE ROLE dbadmin;
GRANT ALL ON basic_auth.users TO dbadmin;
GRANT ALL ON imlswifi.sensors TO dbadmin;
GRANT ALL ON imlswifi.libraries TO dbadmin;
GRANT USAGE ON SCHEMA imlswifi TO dbadmin;
GRANT USAGE ON SCHEMA basic_auth TO dbadmin;
GRANT USAGE ON SCHEMA api TO dbadmin;
GRANT EXECUTE ON FUNCTION api.sensor_setup TO dbadmin;
GRANT EXECUTE ON FUNCTION api.sensor_remove TO dbadmin;
GRANT EXECUTE ON FUNCTION api.update_password TO dbadmin;
GRANT EXECUTE ON FUNCTION api.library_remove TO dbadmin;
GRANT EXECUTE ON FUNCTION api.library_setup TO dbadmin;
GRANT usage, SELECT ON SEQUENCE imlswifi.sensors_sensor_id_seq TO dbadmin;

```




## Components 

This tool is composed of multiple files and libraries used to execute the admin functions

### admin_tool.py
This is the main tool and is executed with Python. You can use --help to find what commands are available.

### admin_tool_kit.py
This is a library that contains general administrative functions 

### admin_tool_sensos.py
This is a library that contains functions related to the management of sensors 

### admin_tool_libraries.py
This is a library that contains functions related to the management of sensors

### 5000-more-common.text
A dictionary used by the tool to generate XKCD passwords

### pdfs folder
This folder contains the instructions that are printed into the pdf and the folder where PDFs are stored when they are generated during sensor creation. 

## Basic Usage
- python admin_tool.py --help
    - Display command options
- python admin_tool.py addSingleLibrary --library fscs_id
	- Used to add a libary to the database (Note: Must be run before adding the sensor for the library)
- python admin_tool.py addSingleSensor --sensor fscs_id --label 'some description'
    - Used to add a sensor and user. This action will create a PDF in the pdfs folder
    - This function also calls testHB which does a heartbeat test
- python admin_tool.py removeSensor --sensor fscs_id 
    - Removes a sensor from the database.
    - NOTE: If there is a record in the heartbeats or presences you cannot remove it.
- python admin_tool.py removeLibrary --library fscs_id 
    - Removes a library from the database.
    - NOTE: Must be run after the sensor associated to the fscs_id is removed.
