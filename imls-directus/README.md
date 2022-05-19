# Directus extensions

Two prototype extensions for wifi sensor work are provided:

- Display: `unix-timestamp`
  - Converts a unix timestamp to a readable format
- Panel: `wifi`
  - Shows daily "wifi minutes served" statistics.

## Building

First, add or edit your `.env` file if you wish to override any environment variables set in `docker-compose.yml`. I typically override `ADMIN_EMAIL` and `ADMIN_PASSWORD`:

    ADMIN_EMAIL="me@gsa.gov"
    ADMIN_PASSWORD="changeme"

This step is not necessary if you wish to retain the `docker-compose.yml` defaults for administrator access.

Second, run `docker-compose up`.

This command spins up three services: `redis`, `postgres`, and `directus`. Once the Directus extensions are built and migrations applied (which might take a while), you will be able to login at `http://localhost:8055` via the administrator credentials above.

## Hydration

Currently the Directus instance has no data, and is only lightly [configured](./extensions/migrations/snapshot.yaml). You will need to either seed with your own data or run a configured `session-counter` to write to this Directus instance.

We plan to allow users to add randomized data: stay tuned.

## Reloading

You can also watch and automatically rebuild extension changes by going to the extension top directory and doing:

    npx directus-extension build -w

## Running locally

You will need to generate your own `.env` file; the easiest way to do so is by running `npm init directus-project` in a new directory, following the prompts (at the moment the migrations only support `postgres`), and copying over the resulting `.env` to this directory.

Then you can run:

    EXTENSIONS_AUTO_RELOAD=true npx directus start
