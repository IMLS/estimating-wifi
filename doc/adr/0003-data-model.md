# 1. Record architecture decisions

Date: 2022-07-22

## Status

Accepted

## Context

In Phase 3, the developers used a singular table called "durations" to record information about the sensor and the duration of the detection between the sensor and patron's device. The "events" table was used for logging but will be deprecated in favor of using Sentry.

For Phase 3.5, the developers revisited the data model approach to clean it up before scaling in Georgia, and to ensure that it would dovetail with accessing it well for data visualization purposes.

The team explored different data models: a relational model with three tables (libraries, sensors, and durations/presences), tables per device, and tables per library.

- We decided the tables per library idea simply represented a grouping of devices, and it would better to be more explicit
- The tables per device idea helped with security: a sensor would write to their own individual table on Directus and have its own API key
  - If the API key were compromised, only the table for that sensor would be compromised, not the whole table
  - ![Data taxonomy for tables per device idea](/doc/images/data_heriachy.jpg)
- The relational model with three tables (libraries, sensors, and durations/presences) was the simplest option in terms of number of tables and leveraging primary key / foreign key relationships
  - This model poses a security risk if the API key for a sensor were compromised, since the whole "sensors" table would be at risk
  - ![ER diagram for sensors table](/doc/images/ER_Diagram_Sensors.png)

The team decided to lower the importance in decision-making of someone compromising the API key and sending false data to the database because it would be challenging to classify false data. Additionally, they lowered the importance of performance as a decision-making factor because at the scale of deployment (hundreds of sensors), the selected data model would not impact performance greatly.

## Decision

For phase 3.5, the team decided to use the relational model with three tables for the sake of simplicity and the level of scaling required within our timeframe for GA.

## Consequences

A device will need to use its API key to write to the durations/presences table. Another way of blocking library device MAC addresses will need to be devised, since the data model selected does not have a configurations table.


