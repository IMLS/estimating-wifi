# Testing the API

To test the API, there's some configuration that is needed first.

## Setting up the .env

You'll first need to set some environment variables. In a `.env` file:

```
export POETRY_SCHEME="https"
export POETRY_HOSTNAME="localhost"
export POETRY_PORT=3000
```

## Sourcing those into the shell environment

You will then need to read those into your shell where you are running tests.

```
source .env
```

You will need to source this every time you open a new shell and want to run tests. `poetry` does not have a mechanism to automatically load environment variables/env files at this time. This might be a FIXME to make our lives easier. (Certainly, this is a CI/CD problem.)

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

Also, to run these tests, there are some things that need to happen to the database. These are not yet automated.

1. ALTER DATABASE needs to happen, inserting `app.jwt_secret` as a parameter of the DB.
1. The `users` table in the `basic_auth` namespace needs to be populated with a user. The test assumes `KY0069-002` and `hello-goodbye` as a user and password. The role should be `sensor`.

You can run the following in a DB shell (DBeaver, or directly via `psql`):

First, 

```
psql -h localhost -U postgres -W
```
and then

```
INSERT INTO basic_auth.users VALUES ('KY0069-002', 'hello-goodbye', 'sensor');
ALTER DATABASE imls SET "app.jwt_secret" TO 'SOMETHINGREALLYLONGLIKEREALLYREALLYLONG';
NOTIFY pgrst, 'reload schema'
```

Note that the secret will, ultimately, need to be automated into our setup/teardown process for production. Currently, there is no production environment, and this is only for local testing.



## Running the tests

Now, we use Poetry to run the tests.

```
source .env ; poetry run pytest
```

is a one-liner you can use to make sure you always have the most recent env variables.
