module.exports = {
  async up(knex) {
    await knex.schema.createTable('durations', (table) => {
      table.increments('id');
      table.string('pi_serial', 16);
      table.string('fcfs_seq_id', 16);
      table.string('device_tag', 32);
      table.string('session_id', 255);
      table.integer('patron_index');
      table.integer('manufacturer_index');
      table.text('start');
      table.text('end');
    });
    await knex.schema.createTable('events', (table) => {
      table.increments('id');
      table.string('pi_serial', 16);
      table.string('fcfs_seq_id', 16);
      table.string('device_tag', 32);
      table.string('session_id', 255);
      table.timestamp('localtime');
      table.timestamp('servertime').defaultTo(knex.fn.now());
      table.string('tag', 255);
      table.text('info');
    });
  },

  async down(knex) {
    await knex.schema.dropTable('durations');
    await knex.schema.dropTable('events');
  },
};
