# Deploy Directus

We suggest deploying a [Directus docker image to cloud.gov](https://cloud.gov/docs/deployment/docker/).

Edit [vars.yml](./app/vars.yml) first. You will want to change, at a minimum, the app name, database name, and route. Then, run:

    cf create-service aws-rds small-psql <directus database name here>

    cf push <directus name here> --docker-image directus/directus:v9.0.0-rc.51

    cd app; cf push --vars-file vars.yml

## Re-exporting environment variables

We want to retrieve system-generated environment variables that were automatically set as part of creating and pushing the above services.

    cf env <directus name here>

The following system-provided variables (from `VCAP_SERVICES` in the above command) will need to be retrieved and re-exported as user environment variables so that the directus container can pick them up:

- `cf set-env <directus name here> DB_CLIENT "pg"`
- `cf set-env <directus name here> DB_PORT "5432"`
- `cf set-env <directus name here> DB_HOST "<your db host>"`
- `cf set-env <directus name here> DB_DATABASE "<your db name>"`
- `cf set-env <directus name here> DB_USER "<your db user>"`
- `cf set-env <directus name here> DB_PASSWORD "<your db password>"`

Create your own key/secret and admin login and export them as well:

- `cf set-env <directus name here> KEY "<your key here>"`
- `cf set-env <directus name here> SECRET "<your secret here>"`
- `cf set-env <directus name here> ADMIN_EMAIL "<your admin email>"`
- `cf set-env <directus name here> ADMIN_PASSWORD "<your admin password"`

Disable the cache:

- `cf set-env <directus name here> CACHE_ENABLED ""`

Finally, `cf restage <directus name here>`.

# Deploy Rabbit

A [manifest.yml](../manifest.yml) file is provided for cloud.gov deployment. You will want to set two environment variables when pushing:

    cf push -f manifest.yml --var secret_key="<your secret key here>" --var
rabbit_magic_header="<rabbit magic header>"

## Miscellany

Even though we use a pipenv environment (`Pipfile`) for local development, we still need to generate a `requirements.txt` file so that cloud.gov can pick the required ReVal dependency up properly. So, if you update any dependencies, re-generate `requirements.txt`:

    pipenv lock -r > requirements.txt
