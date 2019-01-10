package log

import (
	levels "github.com/8Mobius8/go-habits/log/levels"
	"github.com/amoghe/distillog"
)

func NewLevelLogger(level levels.LogLevel, toAttachLogger distillog.Logger) distillog.Logger {
	return LevelLogger{level, toAttachLogger}
}

type LevelLogger struct {
	lvl levels.LogLevel
	log distillog.Logger
}

func (ll LevelLogger) levelIsEnabled(level levels.LogLevel) bool {
	return level >= ll.lvl
}

func (ll LevelLogger) Close() error {
	return ll.log.Close()
}

func (ll LevelLogger) Debugf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Debug) {
		ll.log.Debugf(format, v...)
	}
}

func (ll LevelLogger) Debugln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Debug) {
		ll.log.Debugln(v...)
	}
}

func (ll LevelLogger) Warningf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Warning) {
		ll.log.Warningf(format, v...)
	}
}

func (ll LevelLogger) Warningln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Warning) {
		ll.log.Warningln(v...)
	}
}

func (ll LevelLogger) Infof(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Info) {
		ll.log.Infof(format, v...)
	}
}

func (ll LevelLogger) Infoln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Info) {
		ll.log.Infoln(v...)
	}
}

func (ll LevelLogger) Errorf(format string, v ...interface{}) {
	if ll.levelIsEnabled(levels.Error) {
		ll.log.Errorf(format, v...)
	}
}

func (ll LevelLogger) Errorln(v ...interface{}) {
	if ll.levelIsEnabled(levels.Error) {
		ll.log.Errorln(v...)
	}
}
