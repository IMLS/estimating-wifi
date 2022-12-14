# Notes on Working with the Database

### Run Migrations

- install [dbmate](https://github.com/amacneil/dbmate)
- configure your `.env` for dbmate:
  - `DATABASE_URL="postgres://postgres:imlsimls@localhost:5432/imls?sslmode=disable"`
- Open a separate terminal and run `dbmate up`

## Connect to the DB via CLI

- `psql -h localhost -U postgres -W`
- Enter the password from the docker compose file
- Connect to the database: `\c imls`
- View schemas `\dn`
- View users `\du`
- Use a schema `SET schema 'imlswifi';`
- List tables `\dt`
- List views `\dv`


### Using a GUI DB browser

There are many options for graphical DB browsers. DBeaver is a good option, and the community edition will work just fine. pgAdmin is another

#### DBeaver

[Instructions for installing DBeaver](https://dbeaver.io/download/).

You will need to create a new database connection. The connection parameters will be the same as in the `.env` file in the imls-backend/ folder. If you did not modify the environment variables, you should be able to use the username `postgres`, database `imls`, and password `imlsimlsimls`.

#### pgAdmin

- Ensure you've started the Docker containers as outlined above
- Open a browser and go to `localhost:8080`
- Enter the pgAdmin creds from docker-compose
  - imls@gsa.gov
  - imlsimls
- Click 'Add New Server'
- In the 'General' tab
  - Choose a name (can be anything)
- In the 'Connection' tab
  - The 'Host name/address' is `db`
  - The 'Port' is `5432`
  - The 'Username' is `postgres`
  - The 'Password' is `imlsimls`


## Query the DB

- `curl -s http://127.0.0.1:3000/presences`
- Fields from the `presences` table
  - presence_id SERIAL PRIMARY KEY,
  - start_time TIMESTAMPTZ NOT NULL,
  - end_time TIMESTAMPTZ NOT NULL,
  - fscs_id VARCHAR(16) NOT NULL [FOREIGN KEY],
  - sensor_id SERIAL [FOREIGN KEY],
  - manufacturer_index INTEGER,

## Query Stored Procedures and Functions from the DB

- `curl "http://localhost:3000/rpc/{function_or_sp_name}`
- Current stored procedures and functions in the `api` schema
  - bin_devices_per_hour
    - `curl "http://localhost:3000/rpc/{function_or_sp_name}?_start={DATE variable}&_fscs_id={TEXT variable}"`
    - EXAMPLE:
        `curl "http://localhost:3000/rpc/bin_devices_per_hour?_start=2022-05-10&_fscs_id=AA0003-001"`
    - Returns an array of INTs (device counts per hour) starting at 12 AM EDT, length 24
    - EXAMPLE:
      [22,23,23,27,26,21,23,37,44,50,66,75,75,75,70,88,88,86,70,30,25,25,25,25]

  - bin_devices_over_time
    - `curl “http://localhost:3000/rpc/{function_or_sp_name}?_start={DATE variable}&_fscs_id={TEXT variable}&direction={BOOL varialbe}&_days={INT variable}“`
    - EXAMPLE:
      `curl "http://localhost:3000/rpc/bin_devices_over_time?_start=2022-05-10&_fscs_id=AA0003-001&_direction=true&_days=2"`
    - Returns an array of INTs (device counts per hour) starting at 12 AM EDT, length 24, for Date+1 Day
    - EXAMPLE:  
          [[12,13,13,13,13,13,17,17,15,16,21,22,20,23,20,16,18,21,21,20,20,26,21,21],[26,26,26,25,24,25,25,24,23,27,21,20,18,19,23,17,20,15,18,20,18,15,15,14]]

  - update_presence
    - `curl “http://localhost:3000/rpc/update_presence?_start={TIMESTAMPTZ variable}&_end={TIMESTAMPTZ variable}_fscs_id={CHAR(16) variable}&_sensor={INT varialbe}&_manufacture={INT variable}“`
    - EXAMPLE:
      `curl “http://localhost:3000/rpc/update_presence?_start=2022-09-12 02:21:50+00&_end=2022-09-12 04:21:50+00&_fscs_id=AA0003-001&_sensor=2&_manufacture=7”`
    - Returns sensor_id upon success requires valid JWT
    - `curl "http://localhost:3000/rpc/{function_or_sp_name}`

  - lib_search_fscs
    - `curl "http://localhost:3000/rpc/{function_or_sp_name}?_fscs_id={TEXT variable}"`
    - EXAMPLE:
        `curl "http://localhost:3000/rpc/lib_search_fscs?_fscs_id=AA0003-001"`
    - Returns a single row entry matching the fscs_id
    - EXAMPLE:
      [{"stabr":"AK","fscskey":"AK0001","fscs_seq":2,"c_fscs":"Y","libid":"AK0001-002"...]
  
  - lib_search_state
    - `curl "http://localhost:3000/rpc/{function_or_sp_name}?_state_code={TEXT variable}"`
    - EXAMPLE:
        `curl "http://localhost:3000/rpc/lib_search_state?_state_code=AK"`
    - Returns any row matching the _state_code, case insensative, will search
    on 1-2 chars.
    - EXAMPLE:
      [{"stabr":"AK","fscskey":"AK0001","fscs_seq":2,"c_fscs":"Y","libid":"AK0001-002"...]

  - lib_search_name
    - `curl "http://localhost:3000/rpc/{function_or_sp_name}?_name={TEXT variable}"`
    - EXAMPLE:
        `curl "http://localhost:3000/rpc/lib_search_name?_name=POINT"`
    - Returns any row matching the _point, case insensative, will search for %_name% meaning all chars before or after the _name param.
    - EXAMPLE:
      [{"stabr":"AK","fscskey":"AK0001","fscs_seq":2,"c_fscs":"Y","libid":"AK0001-002"...]


## Persisted Data

- Lives in /imls-backend/data folder

## Errors

If you get error messages from the API that say "could not find the function in the schema cache", `docker-compose kill -s SIGUSR1 server` will reload the PostgREST cache.
