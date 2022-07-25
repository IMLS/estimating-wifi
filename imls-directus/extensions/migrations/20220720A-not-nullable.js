module.exports = {
  async up(knex) {
    await knex.schema.alterTable("durations", (table) => {
      table.string("serial").notNullable();
      table.string("fcfs_seq_id").notNullable();
      table.string("device_tag").notNullable();
      table.string("session_id").notNullable();
      table.integer("patron_index").notNullable();
      table.integer("manufacturer_index").notNullable();
      table.timestamp("start").notNullable();
      table.timestamp("end").notNullable();
    });
  },

  async down(knex) {
    await knex.schema.alterTable("durations", (table) => {
      table.string("serial").nullable();
      table.string("fcfs_seq_id").nullable();
      table.string("device_tag").nullable();
      table.string("session_id").nullable();
      table.integer("patron_index").nullable();
      table.integer("manufacturer_index").nullable();
      table.timestamp("start").nullable();
      table.timestamp("end").nullable();
    });
  },
};
