package levels

import "regexp"

// LogLevel represents a set log level, in the range 0-3 inclusive: Debug 0,
// Warning 1, Info 2, Error 3. Integers are used so that it makes including
// log levels simple using integer comparision.
type LogLevel int

// LogLevel constants
const (
	Debug   = 0
	Warning = 1
	Info    = 2
	Error   = 3
)

// ParseLevelFromString uses regexp package to determine what `LogLevel`
// you would like to use. Useful for making configuration
func ParseLevelFromString(s string) LogLevel {
	matched, _ := regexp.MatchString("DEBUG|debug", s)
	if matched {
		return Debug
	}
	matched, _ = regexp.MatchString("WARN|warn", s)
	if matched {
		return Warning
	}
	matched, _ = regexp.MatchString("INFO|info", s)
	if matched {
		return Info
	}
	matched, _ = regexp.MatchString("ERROR|error", s)
	if matched {
		return Error
	}
	return -1
}
