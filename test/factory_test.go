// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/acmestack/log4go/ext"
	"github.com/acmestack/log4go/logfactory"
	"testing"
)

func TestFactoryTest(t *testing.T) {
	logger := logfactory.GetLogger("test")
	logger.InfoF("this is a %s test\n", "infof")
	logger.InfoLn("this is a infoln test")
	logger.Info("this is a info test\n")
}

func TestFactoryTag(t *testing.T) {
	logger := logfactory.GetLogger()
	logger.WarnLn("test")
	logger = logfactory.GetLogger(nil)
	logger.WarnLn("test")
	logger = logger.WithName("test2")
	logger.WarnLn("test2")
	logger = logger.WithName("test3")
	logger.WarnLn("test3")
	logger = logger.WithFields("FieldKey", "FieldValue")
	logger.WarnLn("test4")
	logger = logger.WithFields("FieldKey", "FieldValue2")
	logger.WarnLn("test5")

	type TestStructInTest struct{}
	logger = logfactory.GetLogger(TestStructInTest{})
	logger.WarnLn("A")

	logger = logfactory.GetLogger(1)
	logger.WarnLn("int")
}

func TestFactorySimplifyName(t *testing.T) {
	fac := logfactory.NewFactory(logfactory.NewLogging())
	fac.SimplifyNameFunc = logfactory.SimplifyNameFirstLetter

	logfactory.ResetFactory(fac)

	type TestStructInTest struct{}
	logger := logfactory.GetLogger(TestStructInTest{})
	logger.WarnLn("A")

	logger = logfactory.GetLogger(1)
	logger.WarnLn("int")
}

func TestMutableFactoryTag(t *testing.T) {

	logfactory.ResetFactory(ext.NewMutableFactory(
		logfactory.NewLogging(
			logfactory.SetLogLevel(logfactory.DEBUG),
			logfactory.SetColorFlag(logfactory.AutoColor))))
	logger := logfactory.GetLogger()
	logger.WarnLn("test")
	logger = logfactory.GetLogger(nil)
	logger.WarnLn("test")
	logger = logger.WithName("test2")
	logger.WarnLn("test2")
	logger = logger.WithName("test3")
	logger.WarnLn("test3")
	logger = logger.WithFields("FieldKey", "FieldValue")
	logger.WarnLn("test4")
	logger = logger.WithFields("FieldKey", "FieldValue2")
	logger.WarnLn("test5")

	type TestStructInTest struct{}
	logger = logfactory.GetLogger(TestStructInTest{})
	logger.WarnLn("A")

	logger = logfactory.GetLogger(1)
	logger.WarnLn("int")
}
