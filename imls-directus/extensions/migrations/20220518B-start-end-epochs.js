module.exports = {
  async up(knex) {
    await knex.schema.alterTable("durations", (table) => {
      table.bigInteger("new_start", { useTz: true });
      table.bigInteger("new_end", { useTz: true });
      // equivalent to `UPDATE durations SET new_start = extract(epoch from "start"::timestamp);` etc
      knex("durations").update({
        new_start: 'extract(epoch from "start"::timestamptz)',
        new_end: 'extract(epoch from "end"::timestamptz)',
      });
      table.dropColumns("start", "end");
      table.renameColumn("new_start", "start");
      table.renameColumn("new_end", "end");
    });
  },

  async down(knex) {
    await knex.schema.alterTable("durations", (table) => {
      table.text("old_start");
      table.text("old_end");
      knex("durations").update({
        old_start: `to_char(to_timestamp("start"), 'YYYY-MM-DD"T"HH24:MI:SSTZH:TZM"')`,
        old_end: `to_char(to_timestamp("end"), 'YYYY-MM-DD"T"HH24:MI:SSTZH:TZM"')`,
      });
      table.dropColumns("start", "end");
      table.renameColumn("old_start", "start");
      table.renameColumn("old_end", "end");
    });
  },
};
