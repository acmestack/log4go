// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/acmestack/log4go/log"
	"github.com/acmestack/log4go/logfactory"
	"github.com/acmestack/log4go/util"
	"testing"
)

func TestFactory(t *testing.T) {

	logfactory.ResetLogging(logfactory.NewLogging(
		logfactory.SetFatalNoTrace(true),
		logfactory.SetExitFunc(func(i int) {
			t.Log("exit: ", i)
		}),
		logfactory.SetPanicFunc(func(v interface{}) {
			if kvs, ok := v.(util.KeyValues); ok {
				t.Log("panic !", kvs.GetAll())
			}
		})),
	)

	type TestStructInTest struct{}
	logger := logfactory.GetLogger(TestStructInTest{})
	logger.InfoF("this is a %s test\n", "infof")
	logger.InfoLn("this is a infoln test")
	logger.Info("this is a info test\n")
}

func TestLogging(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// reset default init at first
		// no fatal trace, do not exit
		log.NewLogging(
			logfactory.SetLogLevel(logfactory.DEBUG),
			logfactory.SetColorFlag(logfactory.AutoColor))

		log.DebugLn("this is a Debugln test")
		log.InfoLn("this is a Infoln test")
		log.WarnLn("this ia a Warnln test")
		log.ErrorLn("this ia a Errorln test")
		//log.PanicLn("this ia a Panicln test")
	})
}
