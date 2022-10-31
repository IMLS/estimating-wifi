# Testing the API

To test the API, there's some configuration that is needed first.

## Clean start

From a clean start

```
rm -rf data ; docker compose up
```

You should be able to get a complete, clean start of the server.

Because all of the code of the application is now encoded as part of the `init`, you will typically need to tear down the server, remove the DB, and restart clean to apply changes. (This is, at least, the most reliable way to get the most recent version of the backend application.) This was chosen as an approach for two reasons:

1. It is easier to track during `dev`, and 
2. We can easily include the app in the container for deployment.

## Setting up the .test_config

You'll first need to set some environment variables. In a `.test_config` file:

```
export POETRY_SCHEME="http"
export POETRY_HOSTNAME="localhost"
export POETRY_PORT=3000
DATABASE_URL="postgres://postgres:imlsimls@localhost:5432/imls?sslmode=disable"
```

These variables are used in the unit tests via `os.getenv()`. The JWT_SECRET needs to be set in accordance with the secret injected into the DB (see below).

## Loading test data

The first time you run the tests, you'll need to source in your test data. You will also need to do this every time you `rm` the `data` directory.

```
source .test_config ; source ./setup-for-tests.sh
```

This runs a script that will create a user/password for unit testing, as well as load a bunch of test data. This should **never** be run in `prod`.

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

## Running the tests

Now, we use Poetry to run the tests.

```
source .env ; poetry run pytest
```

is a one-liner you can use to make sure you always have the most recent env variables.
