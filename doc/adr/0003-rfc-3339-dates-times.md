# 3. RFC-3339 UTC Dates and Times

Date: 2022-07-11

## Status

Proposed

## Context

We need to store dates and times that we can use for comparison purposes.  The
times, specifically, may be in different time zones (e.g., Eastern, Pacific),
Daylight Saving Time, etc..  Therefore, we need a common, standard means for
representing dates and times.

Several options have been considered, including epoch seconds (the number of
seconds elapsed since Midnight on January 1st, 1970 GMT), Microseconds since
January 1st, 2000 GMT (how PostgreSQL internally represents timestampz
values), ISO-8601, RFC-3339, and simply ignoring time zones completely.

While convenient for computational purposes, epoch seconds (and microseconds
since January) are inconveient for humans to quickly interpret).  The
two standards (ISO-8601 and RFC-3339) are very similar in most cases.  For
example, `2022-07-11T14:28:99+00:00` is compliant with both standards.  By
using a `T` between the date and time, avoiding negative zero in the
time offset, and including the colons and hyphens, the vast majority of
dates and times are valid in both formats.

Universal Coordinated Time (UTC) is a common, international standard used
internally by the Network Time Protocol (NTP) upon which Raspberry Pis are
dependent for setting their system clocks (they lack Real Time Clocks (RTCs).
Therefore, converting from UTC to localtime involves an extra conversion
(UTC -> localtime for storage, then localtime -> output) as compared to
storing all dates and times in UTC and converting them once when they are
displayed.

Moreover, UTC does not change with respect to Daylight Saving Time while
localtime calculations may.

Therefore, using UTC involves fewer conversions (so, fewer chances to
introduce conversion errors) and future-proofs storage against Daylight
Saving Time changes (periodic or legislative).

## Decision

RFC-3339 has been chosen to represent dates and times which are to be
represented according to Coordinated Univeral Time (UTC).

## Consequences

Earlier phases of the project used epoch seconds to represent date time
values; these will need to be converted to RFC-3339.

Visualization may need to accept UTC date time values and perform
conversion locally.  If custom middleware is used between the datastore
and the display, conversion may be performed there; however, this is less
likely than having the visualization layer interact directly with the API
(over which we have no control).
