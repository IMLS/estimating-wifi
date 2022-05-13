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

# SQL migrations

If using legacy databases, you may have to create your own table and convert the previous date columns to timedate. Example:

    CREATE TABLE durations2(pi_serial TEXT, fcfs_seq_id TEXT, device_tag TEXT, patron_index INTEGER, id INTEGER PRIMARY KEY AUTOINCREMENT, manufacturer_index INTEGER, session_id TEXT, end datetime, start datetime);

    insert into durations2 select d.pi_serial,d.fcfs_seq_id,d.device_tag,d.patron_index,d.id,d.manufacturer_index,d.session_id,strftime('%s', d.start) * 1000,strftime('%s', d.end) * 1000 from durations d;

