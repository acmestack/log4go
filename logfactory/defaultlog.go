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

package logfactory

import (
	"github.com/acmestack/log4go/util"
)

type defaultlog struct {
	logging Logging
	depth   int
	fields  util.KeyValues
	name    string
}

//// Deprecated: use logfactory.GetLogger instead
//func New(name ...string) Logger {
//	return newLogger(defaultLogging, nil, name...)
//}
//
////Deprecated: use logfactory.GetLogger instead
//func NewLogger(logging Logging, name ...string) Logger {
//	return newLogger(logging, nil, name...)
//}
func defaultLogger(loggingValue Logging, fields util.KeyValues, name ...string) *defaultlog {
	if fields == nil {
		fields = util.NewKeyValues()
	}
	var t string
	if len(name) > 0 {
		t = name[0]
		if t != "" {
			fields.Add(NameKey, t)
		}
	}
	return &defaultlog{
		logging: loggingValue,
		depth:   1,
		name:    t,
		fields:  fields,
	}
}

func (l *defaultlog) DebugEnabled() bool {
	return l.IsEnabled(DEBUG)
}

func (l *defaultlog) Debug(args ...interface{}) {
	l.logging.Log(DEBUG, l.depth, l.fields, args...)
}

func (l *defaultlog) DebugLn(args ...interface{}) {
	l.logging.LogLn(DEBUG, l.depth, l.fields, args...)
}

func (l *defaultlog) DebugF(fmt string, args ...interface{}) {
	l.logging.LogF(DEBUG, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) InfoEnabled() bool {
	return l.IsEnabled(INFO)
}

func (l *defaultlog) Info(args ...interface{}) {
	l.logging.Log(INFO, l.depth, l.fields, args...)
}

func (l *defaultlog) InfoLn(args ...interface{}) {
	l.logging.LogLn(INFO, l.depth, l.fields, args...)
}

func (l *defaultlog) InfoF(fmt string, args ...interface{}) {
	l.logging.LogF(INFO, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) WarnEnabled() bool {
	return l.IsEnabled(WARN)
}

func (l *defaultlog) Warn(args ...interface{}) {
	l.logging.Log(WARN, l.depth, l.fields, args...)
}

func (l *defaultlog) WarnLn(args ...interface{}) {
	l.logging.LogLn(WARN, l.depth, l.fields, args...)
}

func (l *defaultlog) WarnF(fmt string, args ...interface{}) {
	l.logging.LogF(WARN, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) ErrorEnabled() bool {
	return l.IsEnabled(ERROR)
}

func (l *defaultlog) Error(args ...interface{}) {
	l.logging.Log(ERROR, l.depth, l.fields, args...)
}

func (l *defaultlog) ErrorLn(args ...interface{}) {
	l.logging.LogLn(ERROR, l.depth, l.fields, args...)
}

func (l *defaultlog) ErrorF(fmt string, args ...interface{}) {
	l.logging.LogF(ERROR, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) PanicEnabled() bool {
	return l.IsEnabled(PANIC)
}

func (l *defaultlog) Panic(args ...interface{}) {
	l.logging.Log(PANIC, l.depth, l.fields, args...)
}

func (l *defaultlog) PanicLn(args ...interface{}) {
	l.logging.LogLn(PANIC, l.depth, l.fields, args...)
}

func (l *defaultlog) PanicF(fmt string, args ...interface{}) {
	l.logging.LogF(PANIC, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) FatalEnabled() bool {
	return l.IsEnabled(FATAL)
}

func (l *defaultlog) Fatal(args ...interface{}) {
	l.logging.Log(FATAL, l.depth, l.fields, args...)
}

func (l *defaultlog) FatalLn(args ...interface{}) {
	l.logging.LogLn(FATAL, l.depth, l.fields, args...)
}

func (l *defaultlog) FatalF(fmt string, args ...interface{}) {
	l.logging.LogF(FATAL, l.depth, l.fields, fmt, args...)
}

func (l *defaultlog) IsEnabled(severityLevel Level) bool {
	return l.logging.IsEnabled(severityLevel)
}

func (l *defaultlog) WithName(name string) Logger {
	if l == nil {
		return nil
	}

	if l.name != "" {
		name = l.name + "." + name
	}
	ret := defaultLogger(l.logging, l.fields.Clone(), name)
	ret.fields.Add(NameKey, ret.name)
	ret.depth = l.depth

	return ret
}

func (l *defaultlog) WithFields(keyAndValues ...interface{}) Logger {
	if l == nil {
		return nil
	}
	ret := defaultLogger(l.logging, l.fields.Clone(), l.name)
	ret.fields.Add(keyAndValues...)
	ret.depth = l.depth

	return ret
}

func (l *defaultlog) WithDepth(depth int) Logger {
	if l == nil {
		return nil
	}
	ret := defaultLogger(l.logging, l.fields.Clone(), l.name)
	ret.depth += depth

	return ret
}
