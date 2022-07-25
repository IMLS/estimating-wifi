module.exports = {
  async up(knex) {
    await knex.schema.alterTable('durations', (table) => {
      // account for UUIDs on windows.
      table.string('pi_serial', 36).alter();
      table.renameColumn('pi_serial', 'serial');
      // convert from seconds since 01-01-1970 to microseconds since 01-01-2000.
      table.timestamp('tz_start', { useTz: true });
      table.timestamp('tz_end', { useTz: true });
      knex('durations').update({
        tz_start: '(timestamp "epoch" + start * interval "1 second")::timestamptz',
        tz_end: '(timestamp "epoch" + end * interval "1 second")::timestamptz',
      })
      table.dropColumns('start', 'end');
      table.renameColumn('tz_start', 'start');
      table.renameColumn('tz_end', 'end');
    });
  },

  async down(knex) {
    await knex.schema.alterTable('durations', (table) => {
      table.string('serial', 16).alter();
      table.renameColumn('serial', 'pi_serial');
      table.bigInteger('nontz_start', { useTz: true });
      table.bigInteger('nontz_end', { useTz: true });
      knex('durations').update({
        nontz_start: 'extract(epoch from "start")',
        nontz_end: 'extract(epoch from "end")',
      })
      table.dropColumns('start', 'end');
      table.renameColumn('nontz_start', 'start');
      table.renameColumn('nontz_end', 'end');
    });
  },
};
