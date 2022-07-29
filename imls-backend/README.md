# Backend

## Developer configuration

This is for local configuration only and should never be run in production.

### Run Postgrest

- Run `docker-compose up`

### Run migrations

- install [dbmate](https://github.com/amacneil/dbmate)
- configure your `.env` for dbmate:
  - `DATABASE_URL="postgres://postgres:imlsimls@localhost:5432/imls?sslmode=disable"`
- Open a separate terminal and run `dbmate up`

## Query the DB

- `curl -s http://127.0.0.1:3000/presences `

## Persisted Data

- Lives in /imls-backend/data folder