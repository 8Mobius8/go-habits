package log_test

import (
	. "github.com/8Mobius8/go-habits/log"

	"github.com/8Mobius8/go-habits/log/levels"
	"github.com/amoghe/distillog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type logLevelTest struct {
	LogLevel levels.LogLevel
	Format   string
	Contents string
	Debug    bool
	Warn     bool
	Info     bool
	Error    bool
}

var _ = Describe("Logger", func() {
	var testLogger TestLogger
	BeforeEach(func() {
		resetCalled()
		testLogger = TestLogger{}
	})
	It("exists", func() {
		logger := NewLevelLogger(0, testLogger)
		Expect(logger).NotTo(BeNil())
	})

	DescribeTable("with configured log level should call log levels and above",
		func(tcfg logLevelTest) {
			lvlLogger := NewLevelLogger(tcfg.LogLevel, testLogger)
			Expect(lvlLogger).NotTo(BeNil())

			callLogFunctions(lvlLogger, tcfg.Format, tcfg.Contents)
			ExpectLogLevelsToBeCalled(tcfg.Debug, tcfg.Warn, tcfg.Info, tcfg.Error)
		},
		Entry("DEBUG level calls all",
			logLevelTest{levels.Debug, "%s", "test", true, true, true, true}),
		Entry("WARN level calls WARN,INFO,ERROR",
			logLevelTest{levels.Warning, "%s", "test", false, true, true, true}),
		Entry("INFO level calls INFO,ERROR",
			logLevelTest{levels.Info, "%s", "test", false, false, true, true}),
		Entry("ERROR level calls ERROR",
			logLevelTest{levels.Error, "%s", "test", false, false, false, true}),
	)
})

func callLogFunctions(l distillog.Logger, format string, v ...interface{}) {
	l.Debugf(format, v...)
	l.Debugln(v...)
	l.Warningf(format, v...)
	l.Warningln(v...)
	l.Infof(format, v...)
	l.Infoln(v...)
	l.Errorf(format, v...)
	l.Errorln(v...)
}

func ExpectLogLevelsToBeCalled(debug, warn, info, err bool) {
	Expect(DebugfWasCalled).Should(Equal(debug), "Debugf "+wasCalledOrNotCalled(!debug)+" called")
	Expect(DebuglnWasCalled).Should(Equal(debug), "Debugln "+wasCalledOrNotCalled(!debug)+" called")
	Expect(WarningfWasCalled).Should(Equal(warn), "Warningf "+wasCalledOrNotCalled(!warn)+" called")
	Expect(WarninglnWasCalled).Should(Equal(warn), "Warningln "+wasCalledOrNotCalled(!warn)+" called")
	Expect(InfofWasCalled).Should(Equal(info), "Infof "+wasCalledOrNotCalled(!info)+" called")
	Expect(InfolnWasCalled).Should(Equal(info), "Infoln "+wasCalledOrNotCalled(!info)+" called")
	Expect(ErrorfWasCalled).Should(Equal(err), "Errorf "+wasCalledOrNotCalled(!err)+" called")
	Expect(ErrorlnWasCalled).Should(Equal(err), "Errorln "+wasCalledOrNotCalled(!err)+" called")
}

func wasCalledOrNotCalled(wasCalled bool) string {
	if wasCalled {
		return "was"
	} else {
		return "was not"
	}
}

var (
	DebugfWasCalled    bool
	DebuglnWasCalled   bool
	WarningfWasCalled  bool
	WarninglnWasCalled bool
	InfofWasCalled     bool
	InfolnWasCalled    bool
	ErrorfWasCalled    bool
	ErrorlnWasCalled   bool
)

func resetCalled() {
	DebugfWasCalled = false
	DebuglnWasCalled = false
	WarningfWasCalled = false
	WarninglnWasCalled = false
	InfofWasCalled = false
	InfolnWasCalled = false
	ErrorfWasCalled = false
	ErrorlnWasCalled = false
}

type TestLogger struct{}

func (ml TestLogger) Close() error { return nil }

func (ml TestLogger) Debugf(format string, v ...interface{}) {
	DebugfWasCalled = true
}

func (ml TestLogger) Debugln(v ...interface{}) {
	DebuglnWasCalled = true
}

func (ml TestLogger) Errorf(format string, v ...interface{}) {
	ErrorfWasCalled = true
}

func (ml TestLogger) Errorln(v ...interface{}) {
	ErrorlnWasCalled = true
}

func (ml TestLogger) Infof(format string, v ...interface{}) {
	InfofWasCalled = true
}

func (ml TestLogger) Infoln(v ...interface{}) {
	InfolnWasCalled = true
}

func (ml TestLogger) Warningf(format string, v ...interface{}) {
	WarningfWasCalled = true
}

func (ml TestLogger) Warningln(v ...interface{}) {
	WarninglnWasCalled = true
}
