# Testing the API

To test the API, there's some configuration that is needed first.

## Clean start

From a clean start

```
rm -rf data ; docker compose up
```

You should be able to get a complete, clean start of the server. It will also load test data.

## Setting up the .env

You'll first need to set some environment variables. In a `.env` file:

```
export POETRY_SCHEME="https"
export POETRY_HOSTNAME="localhost"
export POETRY_PORT=3000
export POETRY_JWT_SECRET="..."
```

These variables are used in the unit tests via `os.getenv()`. The JWT_SECRET needs to be set in accordance with the secret injected into the DB (see below).

## Sourcing those into the shell environment

You will then need to read those into your shell where you are running tests.

```
source .env
```

You will need to source this every time you open a new shell and want to run tests. `poetry` does not have a mechanism to automatically load environment variables/env files at this time. This might be a FIXME to make our lives easier. (Certainly, this is a CI/CD problem.)

Ultimately,

```
source .env && poetry run pytest
```

would be a robust run pattern. Redundant, but robust.

Our `.gitignore` will ignore this file.

## Setting up a venv

Next, set up a virtual enviornment for installing the Python dependencies for the tests.

```
python -m venv venv
```

and then

```
source venv/bin/activate
```

You will need to source this every time you open a new shell and want to run tests.

(Our `.gitignore` will ignore this directory.)

## Setting up for the tests

```
poetry install
```

Also, to run these tests, there are some things that need to happen to the database. 

1. ALTER DATABASE needs to happen, inserting `app.jwt_secret` as a parameter of the DB and as part of the Postgrest container.
1. The `users` table in the `basic_auth` namespace needs to be populated with a user. The test assumes `KY0069-002` and `hello-goodbye` as a user and password. The role should be `sensor`.

This is automated, at this moment, as part of the `500-ephemeral.sh` script in `init`. However, restarting the containers will not rerun this `init` process, and therefore this is a fragile/poor solution.

Your choice of JWT secret here needs to match your choice in the .env file.

## Running the tests

Now, we use Poetry to run the tests.

```
source .env ; source setup-for-tests.sh ; poetry run pytest
```

is a one-liner you can use to make sure you always have the most recent env variables.
