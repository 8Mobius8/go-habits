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

var logLevelRegexPairs map[LogLevel]string

func init() {
	logLevelRegexPairs = make(map[LogLevel]string)
	logLevelRegexPairs[Debug] = "DEBUG|debug"
	logLevelRegexPairs[Warning] = "WARN|warn"
	logLevelRegexPairs[Info] = "INFO|info"
	logLevelRegexPairs[Error] = "ERROR|error"
}

// ParseLevelFromString uses regexp package to determine what `LogLevel`
// you would like to use. Useful for making configuration
func ParseLevelFromString(s string) LogLevel {
	for ll, regex := range logLevelRegexPairs {
		matched, _ := regexp.MatchString(regex, s)
		if matched {
			return ll
		}
	}

	// Fallback to Error level if log level didn't match
	return Error
}

func matchLevel(s string, ll LogLevel, regex string) bool {
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
