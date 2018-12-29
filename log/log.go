package log

import "github.com/amoghe/distillog"

func NewLevelLogger(level string, toAttachLogger distillog.Logger) distillog.Logger {
	return LevelLogger{level, toAttachLogger}
}

type LevelLogger struct {
	lvl string
	log distillog.Logger
}

func (ll LevelLogger) levelIsEnabled(level string) bool {
	switch ll.lvl {
	case "DEBUG":
		return level == "DEBUG" || level == "WARN" || level == "INFO" || level == "ERROR"
	default:
		return false
	}
}

func (ll LevelLogger) Close() error {
	return ll.log.Close()
}

func (ll LevelLogger) Debugf(format string, v ...interface{}) {
	if ll.levelIsEnabled("DEBUG") {
		ll.log.Debugf(format, v...)
	}
}

func (ll LevelLogger) Debugln(v ...interface{}) {
	if ll.levelIsEnabled("DEBUG") {
		ll.log.Debugln(v...)
	}
}

func (ll LevelLogger) Errorf(format string, v ...interface{}) {
	if ll.levelIsEnabled("ERROR") {
		ll.log.Errorf(format, v...)
	}
}

func (ll LevelLogger) Errorln(v ...interface{}) {
	if ll.levelIsEnabled("ERROR") {
		ll.log.Errorln(v...)
	}
}

func (ll LevelLogger) Infof(format string, v ...interface{}) {
	if ll.levelIsEnabled("INFO") {
		ll.log.Infof(format, v...)
	}
}

func (ll LevelLogger) Infoln(v ...interface{}) {
	if ll.levelIsEnabled("INFO") {
		ll.log.Infoln(v...)
	}
}

func (ll LevelLogger) Warningf(format string, v ...interface{}) {
	if ll.levelIsEnabled("WARN") {
		ll.log.Warningf(format, v...)
	}
}

func (ll LevelLogger) Warningln(v ...interface{}) {
	if ll.levelIsEnabled("WARN") {
		ll.log.Warningln(v...)
	}
}
