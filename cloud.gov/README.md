# Setup

    cf create-service aws-rds small-psql directus-database

    cf push directus-demo2 --docker-image directus/directus:v9.0.0-rc.51

    cd app; cf push --vars-file vars.yml

## Re-exporting environment variables

We want to retrieve system-generated environment variables that were automatically set as part of creating and pushing the above services.

    cf env directus-demo2

The following system-provided variables (from `VCAP_SERVICES` in the above command) will need to be retrieved and re-exported as user environment variables so that the directus container can pick them up:

- `cf set-env directus-demo2 DB_CLIENT "pg"`
- `cf set-env directus-demo2 DB_PORT "5432"`
- `cf set-env directus-demo2 DB_HOST "<your db host>"`
- `cf set-env directus-demo2 DB_DATABASE "<your db name>"`
- `cf set-env directus-demo2 DB_USER "<your db user>"`
- `cf set-env directus-demo2 DB_PASSWORD "<your db password>"`

Create your own key/secret and admin login and export them as well:

- `cf set-env directus-demo2 KEY "<your key here>"`
- `cf set-env directus-demo2 SECRET "<your secret here>"`
- `cf set-env directus-demo2 ADMIN_EMAIL "<your admin email>"`
- `cf set-env directus-demo2 ADMIN_PASSWORD "<your admin password"`

Disable the cache:

- `cf set-env directus-demo2 CACHE_ENABLED ""`

Finally, `cf restage directus-demo2`.
