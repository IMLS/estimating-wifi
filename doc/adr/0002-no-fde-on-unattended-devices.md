# 2. No FDE on Unattended Devices

Date: 2022-07011

## Status

Proposed

## Context

One of the core concepts of using Raspberry Pi devices as sensors is that the
sensors may be physically located in locations that inhibit, make inconvenient,
or even prevent direct access (e.g., mounted on top of lamp posts in library
parking lots).  As a result, these devices, for all intents and purposes, are
unattended.

Encryption -- including Full Disk Encryption (FDE) -- requires a key of some
sort to unencrypt the filesystem(s) stored on the device.  This key may be
stored on a removable storage device (e.g., USB drive) or it may be typed
by an operator at system boot.  This prevents unauthorized access should the
physical media backing the system become compromised.  If the key to unencrypt
the media is present on the device (e.g., stored as a file), the storage
may as well not be encrypted at all.

As a result, using Full Disk Encrpytion (FDE) on a device required an
operator to provide the key in order for the system to boot; without an
operator providing the key, the system will not boot.  Therefore, devices
must either be attended or they may not utilize Full Disk Encryption.  Since
the devices are installed in locations non-condusive to an being attended,
the remaining option is to not use Full Disk Encryption (FDE).

It may be technically possible to create multiple images (FDE and non-FDE);
however, this will complicate support options (e.g., doubling the number of
disk images that need to be built, tested, and maintained).

## Decision

Deployed devices will not use Full Disk Encryption (FDE).

## Consequences

Everything stored on the device will either be unencrypted (i.e., plain-
text) or incorporate a means to unencrypt data.  This includes source code
(likely a non-issue given the source code is publicly available) and
operating system code (similarly available to the public); however, it
will also include sensory data (even if stored ephemerally) as well as
credentials (e.g., API keys) used to interact with the system backend.
