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

package log

import (
	"github.com/acmestack/log4go/logfactory"
	"github.com/acmestack/log4go/util"
)

var (
	Depth util.Value = util.NewSimpleValue(int(1))

	Fields util.Value = util.NewSimpleValue(nil)
)

var defaultLogging util.Value = util.NewSimpleValue(logfactory.NewLogging())

// DefaultLogging 获得默认Logging
func DefaultLogging() logfactory.Logging {
	return defaultLogging.Load().(logfactory.Logging)
}

func NewLogging(opts ...logfactory.LoggingOpt) {
	defaultLogging.Store(logfactory.NewLogging(opts...))
}

// Debug 使用默认的Logging，输出Debug级别的日志
func Debug(args ...interface{}) {
	DefaultLogging().Log(logfactory.DEBUG, Depth.Load().(int), getLogField(), args...)
}

// Debugln 使用默认的Logging，输出Debug级别的日志
func DebugLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.DEBUG, Depth.Load().(int), getLogField(), args...)
}

// Debugf 使用默认的Logging，输出Debug级别的日志
func DebugF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.DEBUG, Depth.Load().(int), getLogField(), fmt, args...)
}

// Info 使用默认的Logging，输出Info级别的日志
func Info(args ...interface{}) {
	DefaultLogging().Log(logfactory.INFO, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Info级别的日志
func InfoLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.INFO, Depth.Load().(int), getLogField(), args...)
}

// Infof 使用默认的Logging，输出Info级别的日志
func InfoF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.INFO, Depth.Load().(int), getLogField(), fmt, args...)
}

// 使用默认的Logging，输出Warn级别的日志
func Warn(args ...interface{}) {
	DefaultLogging().Log(logfactory.WARN, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Warn级别的日志
func WarnLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.WARN, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Warn级别的日志
func WarnF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.WARN, Depth.Load().(int), getLogField(), fmt, args...)
}

// 使用默认的Logging，输出Error级别的日志
func Error(args ...interface{}) {
	DefaultLogging().Log(logfactory.ERROR, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Error级别的日志
func ErrorLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.ERROR, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Error级别的日志
func ErrorF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.ERROR, Depth.Load().(int), getLogField(), fmt, args...)
}

// 使用默认的Logging，输出Panic级别的日志，注意会触发panic
func Panic(args ...interface{}) {
	DefaultLogging().Log(logfactory.PANIC, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Panic级别的日志，注意会触发panic
func PanicLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.PANIC, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Panic级别的日志，注意会触发panic
func PanicF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.PANIC, Depth.Load().(int), getLogField(), fmt, args...)
}

// 使用默认的Logging，输出Fatal级别的日志，注意会触发程序退出
func Fatal(args ...interface{}) {
	DefaultLogging().Log(logfactory.FATAL, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Fatal级别的日志，注意会触发程序退出
func FatalLn(args ...interface{}) {
	DefaultLogging().LogLn(logfactory.FATAL, Depth.Load().(int), getLogField(), args...)
}

// 使用默认的Logging，输出Fatal级别的日志，注意会触发程序退出
func FatalF(fmt string, args ...interface{}) {
	DefaultLogging().LogF(logfactory.FATAL, Depth.Load().(int), getLogField(), fmt, args...)
}

func getLogField() util.KeyValues {
	ret := Fields.Load()
	if ret == nil {
		return nil
	}
	return ret.(util.KeyValues)
}

// 配置默认的调用深度
func WithDepth(depth int) {
	Depth.Store(depth)
}

// 配置默认附件信息
func WithFields(keyAndValues ...interface{}) {
	Fields.Store(util.NewKeyValues(keyAndValues...))
}
