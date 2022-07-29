# Backend

## Developer configuration

This is for local configuration only and should never be run in production.

### Run Postgrest

- `docker-compose up`

### Run migrations

- install [dbmate](https://github.com/amacneil/dbmate)
- configure your `.env` for dbmate:
  - `DATABASE_URL="postgres://postgres:imlsimls@localhost:5432/imls?sslmode=disable"`
- `dbmate up`
