# 1. Record architecture decisions

Date: 2022-09-07

## Status

Accepted

## Context

The authentication process is for preventing an unauthorized user or device from inserting data into the Database.  

Without authentication a requester sends a GET or POST to the API and then the API passes the request to the database to process.

While this is good for a private system, it is not safe for a public system so it requires the implementation of authentication using a method to authenticate the requestor.

## Decision

For the current phase of the project JWT has been chosen as the mechanism for sensor and user authentication. The sensor will pass the bearer token to the API for validation and the API will pass the valid role to the Database.

![AUTH pipeline diagram](/doc/images/auth_piplene_overlay.png)

## Consequences

The JWT will need to be stored on the device for a longer duration to lessen the operational effort to maintain the keys. A rotation time or expiration of the key should be determined.

