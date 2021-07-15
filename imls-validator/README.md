# Rabbit

An experimental validation server for the 10x shared components phase 3 project.

## Usage

Rabbit serves as a stand-alone API proxy for [Directus](https://directus.io/).

There is only one endpoint provided: `/validate/<collection>/`. The only action is `POST`. This endpoint takes an arbitrary array of JSON data, grabs the corresponding validation schema for that collection from directus, validates given data against the schema, and returns the result of validation, successful or otherwise.

Three HTTP headers are required to make this call:

- `X-Magic-Header`: secret key for the rabbit instance
- `X-Directus-Host`: directus host
- `X-Directus-Token`: directus token

Errors from the Directus instance (if any) will be returned verbatim. Otherwise, the endpoint returns a standard [ReVal](https://github.com/18F/ReVAL) validation object in JSON.

## Example

To validate a collection of `baz` objects:

- POST data to `/validate/baz/` with the appropriate headers

- Rabbit will:
  - Write the incoming request to `rabbit_raw` for debugging purposes
  - Retrieve the validation object from the `validators` table, given `baz` as a collection name.
     - Currently this validation object must be a ReVal [GoodtablesValidator](https://specs.frictionlessdata.io/table-schema/#field-descriptors) object.
  - Validate the data
  - If validation passes:
    - Write validated objects to the `baz` table
  - If validation does _not_ pass:
    - Write validation errors to the `rabbit_review` table for manual inspection
  - Return validation result.

## Versioning

We version Directus tables with a suffix of `_v{number}` to preserve backwards compatibility with older clients when we change the schema.

To set the version, pass in the following header:

- `X-Directus-Schema-Version`: directus schema version (defaults to 1)

Assuming the above header is set to 4, in the example above, the `validators_v4`, `rabbit_raw_v4`, `baz_v4`, and `rabbit_review_v4` tables will be used.

# Directus schema

- [directus_tables.sql](./directus_tables.sql) holds the raw SQL schema.

# Deployment

Please see [the cloud.gov README](./cloud.gov/README).
