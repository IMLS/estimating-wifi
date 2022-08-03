# Backend

## Developer configuration

This is for local configuration only and should never be run in production.

### Run Postgrest

- Run `docker-compose up`

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

## Query the DB

- `curl -s http://127.0.0.1:3000/presences `
- Fields from the `presences` table
  - presence_id SERIAL PRIMARY KEY,
  - start_time TIMESTAMPTZ NOT NULL,
  - end_time TIMESTAMPTZ NOT NULL,
  - fscs_id VARCHAR(16) NOT NULL [FOREIGN KEY],
  - sensor_id SERIAL [FOREIGN KEY],
  - manufacturer_index INTEGER,
  - patron_index INTEGER,

## Persisted Data

- Lives in /imls-backend/data folder