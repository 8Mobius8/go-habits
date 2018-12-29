package log_test

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/8Mobius8/go-habits/log"
)

var _ = Describe("Logger", func() {
	BeforeEach(func() {
		DebugfWasCalled, DebuglnWasCalled = false, false
	})
	It("exist", func() {
		logger := NewLevelLogger("", TestLogger{})
		Expect(logger).NotTo(BeNil())
	})
	It("with DEBUG log level will call all log functions", func() {
		testLogger := TestLogger{}
		testLogger.ExpectedFormat = "%s"
		testLogger.ExpectedParameters = []interface{}{"test"}
		testLogger.DebugEnabled = true
		testLogger.WarningEnabled = true
		testLogger.InfoEnabled = true
		testLogger.ErrorEnabled = true

		lvlLogger := NewLevelLogger("DEBUG", testLogger)
		Expect(lvlLogger).NotTo(BeNil())

		lvlLogger.Debugf("%s", "test")
		lvlLogger.Debugln("test")
		lvlLogger.Warningf("%s", "test")
		lvlLogger.Warningln("test")
		lvlLogger.Infof("%s", "test")
		lvlLogger.Infoln("test")
		lvlLogger.Errorf("%s", "test")
		lvlLogger.Errorln("test")

		Expect(DebugfWasCalled).Should(BeTrue(), "Debugf was not called")
		Expect(DebuglnWasCalled).Should(BeTrue(), "Debugln was not called")
		Expect(WarningfWasCalled).Should(BeTrue(), "Warningf was not called")
		Expect(WarninglnWasCalled).Should(BeTrue(), "Warningln was not called")
		Expect(InfofWasCalled).Should(BeTrue(), "Infof was not called")
		Expect(InfolnWasCalled).Should(BeTrue(), "Infoln was not called")
		Expect(ErrorfWasCalled).Should(BeTrue(), "Errorf was not called")
		Expect(ErrorlnWasCalled).Should(BeTrue(), "Errorln was not called")
	})
	It("with WARN log level will call WARN,INFO,ERROR log functions", func() {
		testLogger := TestLogger{}
		testLogger.ExpectedParameters = []interface{}{"test"}
		testLogger.ExpectedFormat = "%s"
		testLogger.DebugEnabled = false

		lvlLogger := NewLevelLogger("INFO", testLogger)
		Expect(lvlLogger).NotTo(BeNil())

		lvlLogger.Debugf("%s", "test")
		lvlLogger.Debugln("test")

		Expect(DebugfWasCalled).To(BeFalse())
		Expect(DebuglnWasCalled).To(BeFalse())
	})

	// lvlLogger.Errorf(logContents)
	// lvlLogger.Errorln(logContents)
	// lvlLogger.Infof(logContents)
	// lvlLogger.Infoln(logContents)
	// lvlLogger.Warningf(logContents)
	// lvlLogger.Warningln(logContents)
})

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

type TestLogger struct {
	DebugEnabled   bool
	WarningEnabled bool
	InfoEnabled    bool
	ErrorEnabled   bool

	ExpectedParameters []interface{}
	ExpectedFormat     string
}

func (ml TestLogger) Close() error { return nil }

func (ml TestLogger) Debugf(format string, v ...interface{}) {
	DebugfWasCalled = true
	if !ml.DebugEnabled {
		ginkgo.Fail("`Debugf` function should not be called")
	}
	ml.ExpectPrintf(format, v...)
}

func (ml TestLogger) Debugln(v ...interface{}) {
	DebuglnWasCalled = true
	if !ml.DebugEnabled {
		ginkgo.Fail("`Debugln` function should not be called")
	}
	ml.ExpectLogParameters(v...)
}

func (ml TestLogger) Errorf(format string, v ...interface{}) {
	ErrorfWasCalled = true
	if !ml.ErrorEnabled {
		ginkgo.Fail("`Errorf` function should not be called")
	}
	ml.ExpectPrintf(format, v...)
}

func (ml TestLogger) Errorln(v ...interface{}) {
	ErrorlnWasCalled = true
	if !ml.ErrorEnabled {
		ginkgo.Fail("`Errorln` function should not be called")
	}
	ml.ExpectLogParameters(v...)
}

func (ml TestLogger) Infof(format string, v ...interface{}) {
	InfofWasCalled = true
	if !ml.InfoEnabled {
		ginkgo.Fail("`Infof` function should not be called")
	}
	ml.ExpectPrintf(format, v...)
}

func (ml TestLogger) Infoln(v ...interface{}) {
	InfolnWasCalled = true
	if !ml.InfoEnabled {
		ginkgo.Fail("`Infoln` function should not be called")
	}
	ml.ExpectLogParameters(v...)
}

func (ml TestLogger) Warningf(format string, v ...interface{}) {
	WarningfWasCalled = true
	if !ml.WarningEnabled {
		ginkgo.Fail("`Warningf` function should not be called")
	}
	ml.ExpectPrintf(format, v...)
}

func (ml TestLogger) Warningln(v ...interface{}) {
	WarninglnWasCalled = true
	if !ml.WarningEnabled {
		ginkgo.Fail("`Warningln` function should not be called")
	}
	ml.ExpectLogParameters(v...)
}

func (ml TestLogger) ExpectPrintf(format string, v ...interface{}) {
	ml.ExpectLogFormat(format)
	ml.ExpectLogParameters(v...)
}

func (ml TestLogger) ExpectLogParameters(v ...interface{}) {
	Expect(v).To(BeEquivalentTo(ml.ExpectedParameters))
}

func (ml TestLogger) ExpectLogFormat(format string) {
	Expect(format).To(BeEquivalentTo(ml.ExpectedFormat))
}
