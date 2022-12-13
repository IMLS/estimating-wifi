# Building and Testing imls-backend

The imls-backend component of this repository runs a RESTful API server built with [Postgrest](https://postgrest.org/en/stable/). Using Docker and Python-Poetry, it is possible to build, run, and test the backend in a local development environment.

## Prerequisites

* git clone this repository locally
* A functioning Docker daemon running locally
* for unit tests, [install Poetry](https://python-poetry.org/docs/#installing-with-the-official-installer).

## Building container images

There are multiple steps here, but the build process should not have to happen often.

Navigate to the `imls-backend` directory.
```
cd imls-backend
```

The backend containers run on two images. Our `postgres` image is extended with security extensions, and must be built before you can proceed. 

```
docker build -t imls/postgres:latest -f Dockerfile.pgjwt .
```

Once you have built this container, you are ready to bring up the backend.

### Clearing the database (not necessary on first run)

During development it's often useful to completely refresh the database backend, including repopulating the database. Removing the `data/` subdirectory accomplishes this.

Make sure you are in the `imls-backend` directory so you don't blow away the wrong `data/` directory.
```
# confirm you are in the imls-backend dir, then delete
pwd
rm -rf data
```
### Setting up the Docker environment

To run the backend, you need a `.env` file in the same directory as the `docker-compose.yml` file. The `.env` file is configured with values specific for these containers, and needs to be created locally.

Create `imls-backend/.env` using a text editor of your choice, and set values for each of the following variables. You may copy-paste the values below directly into the file and things should "just work."

```
PGRST_JWT_SECRET="EmpowerMuseumsLibrariesGrantmakingResearchPolicyDevelopment"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="imlsimlsimls"
POSTGRES_DB="imls"
DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable"
```

**The above values are for development on LOCAL MACHINES ONLY**. Environment variables displayed here are inappropriate to deploy to a production environment.

### Running the backend 

Once you have your environment variables in place, you should be able to use the compose file to run the backend for development purposes.

```
docker compose up
```

If this is the first time you are running our stack, a Postgrest (API) and Postgres (DB) base image will be downloaded before launch. The launch sequence involves the API container waiting for the DB container. When the DB container starts for the first time, it will load a substantial amount of SQL from the `init` directory. This *only* happens the first time the containers start up; on all subsequent restarts of the system, the `init` directory is not read. You can learn more about this [from the Postgres container documentation](https://hub.docker.com/_/postgres). (This behavior is built into the base image, and is not the result of our work.)

### Optional: Loading test data

At this point, you have a an empty system. It has no users and no data but is ready for use. If you're new to the stack, or want to do anything interesting at all, you're going to need some data and at least one user.

`cd` to `imls-backend/test`.

This directory has some test data to populate the database. **The data are not actually from the libraries named in the file**. Using invalid library IDs breaks many things, so there are *real* library IDs attached to *fake* data.

We have a small script you can run to load that data.

We recommend running this script as follows:

```
./setup-for-tests.sh
```

This will read in the needed environment variables for your containers, and then run a sequence of SQL commands to insert test data as well as create a user. That user is library `KY0069-002` and the API key for that user is `hello-goodbye`.

## Testing everything

You should now be able to run the unit tests. This is a good test of whether your stack is functional. The tests in this file will only pass if you have loaded the test data in the previous step.

Poetry will read the `pyproject.toml` file and install all dependencies into a virtual environment with the command (issued from the imls-backend/ folder)

```
poetry install
```

After this, you can drop into the virtual environment with the command

```
poetry shell
```

Now run the tests
```
source .env ; cd test && poetry run pytest
```

If all goes well, you'll see all the tests pass.

### Optional: Using a DB browser

There are many options for DB browsers. We recommend using DBeaver. The community edition will work just fine.

[Instructions for installing DBeaver](https://dbeaver.io/download/).

You will need to create a new database connection. The connection parameters will be the same as in the `.env` file in the imls-backend/ folder. If you did not modify the environment variables, you should be able to use the username `postgres`, database `imls`, and password `imlsimlsimls`.

## What now?

You have now stood up the API and data storage backend. This does not include either the web-based frontend for browsing the data or the code that estimates wifi device presence. So, you're only part way there if you're looking to stand up the entire system on a dev machine.
