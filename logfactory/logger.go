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

// LogDebug interface
type LogDebug interface {
	DebugEnabled() bool
	Debug(args ...interface{})
	DebugLn(args ...interface{})
	DebugF(fmt string, args ...interface{})
}

// LogInfo interface
type LogInfo interface {
	InfoEnabled() bool
	Info(args ...interface{})
	InfoLn(args ...interface{})
	InfoF(fmt string, args ...interface{})
}

// LogWarn interface
type LogWarn interface {
	WarnEnabled() bool
	Warn(args ...interface{})
	WarnLn(args ...interface{})
	WarnF(fmt string, args ...interface{})
}

// LogError interface
type LogError interface {
	ErrorEnabled() bool
	Error(args ...interface{})
	ErrorLn(args ...interface{})
	ErrorF(fmt string, args ...interface{})
}

// LogPanic interface
type LogPanic interface {
	PanicEnabled() bool
	Panic(args ...interface{})
	PanicLn(args ...interface{})
	PanicF(fmt string, args ...interface{})
}

// LogFatal interface
type LogFatal interface {
	FatalEnabled() bool
	Fatal(args ...interface{})
	FatalLn(args ...interface{})
	FatalF(fmt string, args ...interface{})
}

// Logger interface 实现了常用的日志方法
type Logger interface {
	LogDebug
	LogInfo
	LogWarn
	LogError
	// LogPanic Panic level log interface. Note that panic will be triggered
	LogPanic
	// LogFatal Fatal level log interface, please note that it will trigger the program exit
	LogFatal

	// WithName 附加日志名称，注意会附加父Logger的名称，格式为：父Logger名称 + '.' + name
	WithName(name string) Logger

	// WithFields 附加日志信息，注意会附加父Logger的附加信息，如果相同则会覆盖
	WithFields(keyAndValues ...interface{}) Logger

	// WithDepth 配置日志的调用深度，注意会在父Logger的基础上调整深度
	WithDepth(depth int) Logger
}
