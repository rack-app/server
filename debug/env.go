package debug

import "os"

var inDebug bool

func init() {
	value, isSet := os.LookupEnv("RACK_APP_DEBUG")
	inDebug = isSet && value == "TRUE"
}
