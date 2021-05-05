package constants

const MACLENGTH = 17

const ExitNoUsername = -100
const ExitNoPassword = -101
const ExitProcessTimeout = -200

// 20210505 MCJ
// Making these variables so that tests can change
// the location...
var ConfigPath = "/opt/imls/config.yaml"
var AuthPath = "/opt/imls/auth.yaml"

const AuthTokenKey = "AUTHTOKEN"
const AuthEmailKey = "AUTHEMAIL"
