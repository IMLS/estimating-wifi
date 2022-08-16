Date: 2022-08-12

# Update to Data Model ADR

The team decided to update the data model. The sensors and libraries tables will be loaded during the bulk install process, so all of the fields need to be accessible to the State IT Admin Director for the library systems they will set up. The serial and version of the sensor will not be available at bulk install time for the IT director, so the team decided to move these fields to the heartbeats table.

The sensors and libraries tables will be populated and updated during the bulk install or update processes. The heartbeats and presences tables will be populated after the bulk installation process.

The team also decided to remove the patron_index as a field from the presences table. The imls-raspberry-pi directory does currently does not track return patrons, so the team decided that this feature was out of scope for the time being.

![ER diagram for the IMLS Wifi project](/doc/images/ER_Diagram_IMLS_Wifi_v3.png)

Date: 2022-08-09

# Update to Data Model ADR

The team decided to update the data model to add a heartbeats table. This table will enable the developers and users to see how often the sensor is reporting to the database without issues. The sensor will write to the heartbeat table every hour with its status.

![ER diagram for the IMLS Wifi project](/doc/images/ER%20Diagram%20-%20IMLS%20Wifi_v2.png)

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


