package constants

const MACLENGTH = 17

// For how long do we recognize a device?
// 2 hours. This is 2 * 60 minutes * 60 seconds.
// If we see a MAC within this window, we "remember" it.
// If we see a MAC, 2h go by, and we see it again, we're going
// to "forget" the original sighting, and pretend the device is new.
// const MAC_MEMORY_DURATION_SEC = 2 * 60 * 60
