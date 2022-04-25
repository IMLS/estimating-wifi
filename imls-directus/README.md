# Directus extensions

Two prototype extensions for wifi sensor work are provided:

- Display: `unix-timestamp`
  - Converts a unix timestamp to a readable format
- Panel: `wifi`
  - Shows daily "wifi minutes served" statistics.

## Building

Run the following to enable automatic reloading of extensions:

    EXTENSIONS_AUTO_RELOAD=true npx directus start

You can also watch and automatically rebuild extension changes by going to the extension top directory and doing:

    npx directus-extension build -w

Running both of these commands at the same time will provide a live reloading experience.
