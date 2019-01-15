package log

import (
	levels "github.com/8Mobius8/go-habits/log/levels"
	"github.com/amoghe/distillog"
)

// LevelLogger wrapper for distillog.Log that will print messages at configured level.
// Maybe be created with any logger that supports the distillog.Logger interface
type LevelLogger struct {
	lvl levels.LogLevel
	log distillog.Logger
}

// NewLevelLogger creates a new logger from the defined level and distillog.Logger
func NewLevelLogger(level levels.LogLevel, toAttachLogger distillog.Logger) distillog.Logger {
	return LevelLogger{level, toAttachLogger}
}

func (ll LevelLogger) levelIsEnabled(level levels.LogLevel) bool {
	return level >= ll.lvl
}

// Close calls internal distillog.Logger.Close
func (ll LevelLogger) Close() error {
	return ll.log.Close()
}

// Debugf calls internal distillog.Logger.Debugf if enabled
func (ll LevelLogger) Debugf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Debug) {
		ll.log.Debugf(format, v...)
	}
}

// Debugln calls internal distillog.Logger.Debugln if enabled
func (ll LevelLogger) Debugln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Debug) {
		ll.log.Debugln(v...)
	}
}

// Warningf calls internal distillog.Logger.Warningf if enabled
func (ll LevelLogger) Warningf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Warning) {
		ll.log.Warningf(format, v...)
	}
}

// Warningln calls internal distillog.Logger.Warningln if enabled
func (ll LevelLogger) Warningln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Warning) {
		ll.log.Warningln(v...)
	}
}

// Infof calls internal distillog.Logger.Infof if enabled
func (ll LevelLogger) Infof(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Info) {
		ll.log.Infof(format, v...)
	}
}

// Infoln calls internal distillog.Logger.Infoln if enabled
func (ll LevelLogger) Infoln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Info) {
		ll.log.Infoln(v...)
	}
}

// Errorf calls internal distillog.Logger.Errorf if enabled
func (ll LevelLogger) Errorf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Error) {
		ll.log.Errorf(format, v...)
	}
}

// Errorln calls internal distillog.Logger.Errorln if enabled
func (ll LevelLogger) Errorln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Error) {
		ll.log.Errorln(v...)
	}
}
