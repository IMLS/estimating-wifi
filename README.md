# Rabbit

An experimental validation server for the 10x shared components phase 3 project.

## Usage

Rabbit serves as a stand-alone API proxy for [Directus](https://directus.io/).

There is only one endpoint provided: `/validate/<collection>/`. The only action is `POST`. This endpoint takes an arbitrary array of JSON data, grabs the corresponding validation schema for that collection from directus, and returns the result.

Four HTTP headers are required to make this call:

- `X-Magic-Header`: secret key for the rabbit instance
- `X-Directus-Host`: directus host
- `X-Directus-Email`: directus user
- `X-Directus-Password`: directus password

Errors from the Directus instance (if any) will be returned verbatim. Otherwise, the endpoint returns a standard [ReVal](https://github.com/18F/ReVAL) validation object in JSON.

## Example

To validate a collection of `baz` objects:

- POST data to `/validate/baz/` with the appropriate headers

- Rabbit will:
  - Authenticate against the Directus instance given the host and credentials
  - Retrieve the validation object from its `validation` table, given `baz` as a collection name.
     - Currently this validation object must be a ReVal GoodtablesValidator object.
  - Write the incoming request to `rabbit_raw` for debugging purposes
  - Validate the data
  - If validation passes:
    - Write validated objects to the `baz` table
  - If validation does _not_ pass:
    - Write validation errors to the `rabbit_review` table for manual inspection
  - Return validation result.

# Directus SQL

Forthcoming.
