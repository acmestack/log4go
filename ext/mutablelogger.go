/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ext

import (
	"github.com/acmestack/log4go/logfactory"
	"github.com/acmestack/log4go/util"
)

type mutableLog struct {
	logging util.Value
	depth   int
	fields  util.KeyValues
	name    string
}

type mutableLoggerFactory struct {
	logfactory.LoggerFactory
}

func NewMutableFactory(logging logfactory.Logging) *mutableLoggerFactory {
	return NewMutableFactoryWithValue(util.NewAtomicValue(logging))
}

func NewMutableFactoryWithValue(v util.Value) *mutableLoggerFactory {
	ret := &mutableLoggerFactory{}
	ret.Value = v
	return ret
}

func (fac *mutableLoggerFactory) GetLogger(o ...interface{}) logfactory.Logger {
	name := util.GetObjectName(fac.SimplifyNameFunc, o...)
	return newMutableLogger(fac.Value, nil, name)
}

func newMutableLogger(loggingValue util.Value, fields util.KeyValues, name ...string) *mutableLog {
	if fields == nil {
		fields = util.NewKeyValues()
	}
	var t string
	if len(name) > 0 {
		t = name[0]
		if t != "" {
			fields.Add(logfactory.NameKey, t)
		}
	}
	ret := &mutableLog{
		logging: loggingValue,
		depth:   1,
		name:    t,
		fields:  fields,
	}
	return ret
}

func (l *mutableLog) getLogging() logfactory.Logging {
	return l.logging.Load().(logfactory.Logging)
}

func (l *mutableLog) DebugEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.DEBUG)
}

func (l *mutableLog) Debug(args ...interface{}) {
	l.getLogging().Log(logfactory.DEBUG, l.depth, l.fields, args...)
}

func (l *mutableLog) DebugLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.DEBUG, l.depth, l.fields, args...)
}

func (l *mutableLog) DebugF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.DEBUG, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) InfoEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.INFO)
}

func (l *mutableLog) Info(args ...interface{}) {
	l.getLogging().Log(logfactory.INFO, l.depth, l.fields, args...)
}

func (l *mutableLog) InfoLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.INFO, l.depth, l.fields, args...)
}

func (l *mutableLog) InfoF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.INFO, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) WarnEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.WARN)
}

func (l *mutableLog) Warn(args ...interface{}) {
	l.getLogging().Log(logfactory.WARN, l.depth, l.fields, args...)
}

func (l *mutableLog) WarnLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.WARN, l.depth, l.fields, args...)
}

func (l *mutableLog) WarnF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.WARN, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) ErrorEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.ERROR)
}

func (l *mutableLog) Error(args ...interface{}) {
	l.getLogging().Log(logfactory.ERROR, l.depth, l.fields, args...)
}

func (l *mutableLog) ErrorLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.ERROR, l.depth, l.fields, args...)
}

func (l *mutableLog) ErrorF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.ERROR, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) PanicEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.PANIC)
}

func (l *mutableLog) Panic(args ...interface{}) {
	l.getLogging().Log(logfactory.PANIC, l.depth, l.fields, args...)
}

func (l *mutableLog) PanicLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.PANIC, l.depth, l.fields, args...)
}

func (l *mutableLog) PanicF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.PANIC, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) FatalEnabled() bool {
	return l.getLogging().IsEnabled(logfactory.FATAL)
}

func (l *mutableLog) Fatal(args ...interface{}) {
	l.getLogging().Log(logfactory.FATAL, l.depth, l.fields, args...)
}

func (l *mutableLog) FatalLn(args ...interface{}) {
	l.getLogging().LogLn(logfactory.FATAL, l.depth, l.fields, args...)
}

func (l *mutableLog) FatalF(fmt string, args ...interface{}) {
	l.getLogging().LogF(logfactory.FATAL, l.depth, l.fields, fmt, args...)
}

func (l *mutableLog) IsEnabled(severityLevel logfactory.Level) bool {
	return l.getLogging().IsEnabled(severityLevel)
}

func (l *mutableLog) WithName(name string) logfactory.Logger {
	if l == nil {
		return nil
	}

	if l.name != "" {
		name = l.name + "." + name
	}
	ret := newMutableLogger(l.logging, l.fields.Clone(), name)
	ret.fields.Add(logfactory.NameKey, ret.name)
	ret.depth = l.depth

	return ret
}

func (l *mutableLog) WithFields(keyAndValues ...interface{}) logfactory.Logger {
	if l == nil {
		return nil
	}
	ret := newMutableLogger(l.logging, l.fields.Clone(), l.name)
	ret.fields.Add(keyAndValues...)
	ret.depth = l.depth

	return ret
}

func (l *mutableLog) WithDepth(depth int) logfactory.Logger {
	if l == nil {
		return nil
	}
	ret := newMutableLogger(l.logging, l.fields.Clone(), l.name)
	ret.depth += depth

	return ret
}
