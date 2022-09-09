# Data at import

For local development, we have multiple data files.

`test-data.sql.gz` contains *real* data that is anonymized. It is pulled from the Phase 3 data collection, and then updated to match the current data schema.

`imls-data-2020.csv` is data from the Institute of Museum and Library Services (IMLS). Specifically, we are using outlet data from 2020. It is public domain data, and we sourced it from https://www.imls.gov/research-evaluation/data-collection/public-libraries-survey.

The IMLS data file is loaded via `\copy`, and therefore we need both an `.sql` file and the `.csv`; the `.sql` file defines the table, and then uses the Postgres-specific `\copy` command to load the CSV file into that table. 

# To use this data

You will need to:

1. Delete your `data` directory.
2. `docker compose up`
3. Run your migrations (`dbmate up`)

The test data and IMLS data are loaded during step two. However, they are *only* loaded the *first time* the database is created (when it is initialized). Hence, the `data` directory must be cleared/removed and the containers spun up again "for the first time."