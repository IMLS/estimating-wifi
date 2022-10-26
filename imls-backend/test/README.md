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

## Running the tests

Now, we use Poetry to run the tests.

```
source .env ; poetry run pytest
```

is a one-liner you can use to make sure you always have the most recent env variables.

