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
	"bytes"
	"fmt"
	"github.com/acmestack/log4go/util"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Level = int32

const (
	FATAL Level = 0
	PANIC Level = 1
	ERROR Level = 2
	WARN  Level = 3
	INFO  Level = 4
	DEBUG Level = 5
)

const (
	// CallerNone none
	CallerNone = 0
	// CallerShortFile short file name
	CallerShortFile = 1
	// CallerLongFile long file name
	CallerLongFile = 1 << 1
	// CallerFileMask long file name
	CallerFileMask = 3
	// CallerShortFunc caller short func
	CallerShortFunc = 1 << 2
	// CallerLongFunc caller long func
	CallerLongFunc = 1 << 3
	// CallerSimpleFunc caller simple func
	CallerSimpleFunc = 1 << 4
	// CallerFuncMask func mask
	CallerFuncMask = 7 << 2
)

const (
	AutoColor = iota

	DisableColor
)

const (
	// TimestampKey LogTime
	TimestampKey = "LogTime"
	// LevelKey LogLevel
	LevelKey = "LogLevel"
	// CallerKey LogCaller
	CallerKey = "LogCaller"
	// ContentKey LogContent
	ContentKey = "LogContent"
	// NameKey LogName
	NameKey = "LogName"
)

var (
	//前景色
	ColorGreen   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	ColorWhite   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	ColorYellow  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	ColorRed     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	ColorBlue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	ColorMagenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	ColorCyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	ColorReset   = string([]byte{27, 91, 48, 109})

	ForeGreen   = "\033[97;32m"
	ForeWhite   = "\033[90;37m"
	ForeYellow  = "\033[90;33m"
	ForeRed     = "\033[97;31m"
	ForeBlue    = "\033[97;34m"
	ForeMagenta = "\033[97;35m"
	ForeCyan    = "\033[97;36m"

	//背景色
	BackGreen   = "\033[97;42m"
	BackWhite   = "\033[90;47m"
	BackYellow  = "\033[90;43m"
	BackRed     = "\033[97;41m"
	BackBlue    = "\033[97;44m"
	BackMagenta = "\033[97;45m"
	BackCyan    = "\033[97;46m"

	ResetColor = "\033[0m"
)

// 级别及名称映射
var LogTag = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	PANIC: "PANIC",
	FATAL: "FATAL",
}

// 默认值
var (
	DefaultColorFlag     = DisableColor
	DefaultPrintFileFlag = CallerShortFile
	DefaultFatalNoTrace  = false
	DefaultLevel         = INFO
	DefaultWriters       = map[Level]io.Writer{
		DEBUG: os.Stdout,
		INFO:  os.Stdout,
		WARN:  os.Stdout,
		ERROR: os.Stderr,
		PANIC: os.Stderr,
		FATAL: os.Stderr,
	}
)

type LoggingOpt func(l *logging)

// Logging base
type Logging interface {
	LogF(level Level, depth int, keyValues util.KeyValues, format string, args ...interface{})

	Log(level Level, depth int, keyValues util.KeyValues, args ...interface{})

	LogLn(level Level, depth int, keyValues util.KeyValues, args ...interface{})

	// SetFormatter setting Formatter
	SetFormatter(f util.Formatter)

	// SetLogLevel 设置日志严重级别，低于该级别的将不被输出（线程安全）
	SetLogLevel(severityLevel Level)

	// IsEnabled 判断参数级别是否会输出（线程安全）
	IsEnabled(severityLevel Level) bool

	// SetOutput 设置输出的Writer，注意该方法会将所有级别都配置为参数writer（线程安全）
	SetOutput(w io.Writer)

	// SetOutputBySeverity 设置对应日志级别的Writer（线程安全）
	SetOutputBySeverity(severityLevel Level, w io.Writer)

	// GetOutputBySeverity 获得对应日志级别的Writer（线程安全）
	GetOutputBySeverity(severityLevel Level) io.Writer

	// Clone 获得一个clone的对象（线程安全）
	Clone() Logging
}

type ExitFunc func(code int)
type PanicFunc func(interface{})

type logging struct {
	timeFormatter   func(t time.Time) string
	callerFormatter func(file string, line int, funcName string) string
	exitFunc        ExitFunc
	panicFunc       PanicFunc
	formatter       atomic.Value
	colorFlag       int
	fileFlag        int
	fatalNoTrace    bool

	level Level

	writers sync.Map

	bufPool sync.Pool
}

func SimplifyNameFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return s[:1]
}

var defaultLogging util.Value = util.NewSimpleValue(NewLogging())

// DefaultLogging 获得默认Logging
func DefaultLogging() Logging {
	return defaultLogging.Load().(Logging)
}

func GetLogging() util.Value {
	return defaultLogging
}

func NewLogging(opts ...LoggingOpt) Logging {
	ret := &logging{
		timeFormatter:   timeFormat,
		callerFormatter: callerFormat,
		exitFunc:        defaultExit,
		panicFunc:       defaultPanic,
		//formatter:     nil,
		colorFlag:    DefaultColorFlag,
		fileFlag:     DefaultPrintFileFlag,
		fatalNoTrace: DefaultFatalNoTrace,
		level:        DefaultLevel,

		bufPool: sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(nil)
		}},
	}

	for k, v := range DefaultWriters {
		ret.writers.Store(k, v)
	}

	for _, v := range opts {
		v(ret)
	}
	return ret
}

func (l *logging) getCaller(depth int) string {
	if l.fileFlag != CallerNone {
		pc, file, line, ok := Caller(depth+3, true)
		if !ok {
			return "???"
		}

		if (l.fileFlag & CallerShortFile) != 0 {
			file = shortFile(file)
		}

		if (l.fileFlag & CallerFileMask) == 0 {
			file = ""
			line = -1
		}
		var funcName string
		if (l.fileFlag & CallerFuncMask) != 0 {
			funcName = runtime.FuncForPC(pc).Name()
			if (l.fileFlag & CallerShortFunc) != 0 {
				idx := strings.LastIndex(funcName, ".")
				if idx != -1 && idx < (len(funcName)-1) {
					funcName = funcName[idx+1:]
				}
			} else if (l.fileFlag & CallerSimpleFunc) != 0 {
				funcName = simpleFuncName(funcName)
			}
		}
		return l.callerFormatter(file, line, funcName)
	}
	return ""
}

func (l *logging) format(writer io.Writer, level Level, depth int, keyValues util.KeyValues, log string) {
	caller := l.getCaller(depth)

	var (
		lvColor    string
		resetColor string
	)
	if l.colorFlag == AutoColor {
		lvColor = selectLevelColor(level)
		resetColor = ResetColor
	}

	formatter := l.formatter.Load()
	if formatter != nil {
		innerKvs := util.NewKeyValues()
		_ = innerKvs.Add(TimestampKey, time.Now(), LevelKey, LogTag[level], CallerKey, caller)
		_, _ = util.MergeKeyValues(innerKvs, keyValues)
		if log == "\n" {
			log = ""
		}
		_ = innerKvs.Add(ContentKey, log)
		_ = formatter.(util.Formatter).Format(writer, innerKvs)
	} else {
		_, _ = writer.Write([]byte(fmt.Sprintf("%s [%s%s%s] %s %s%s",
			l.timeFormatter(time.Now()), lvColor, LogTag[level], resetColor, caller, l.formatKeyValues(keyValues), log)))
	}
}

func (l *logging) formatKeyValues(keyValues util.KeyValues) string {
	if keyValues == nil || keyValues.Len() == 0 {
		return ""
	}

	buf := bytes.Buffer{}
	for _, k := range keyValues.Keys() {
		buf.WriteString(l.formatValue(keyValues.Get(k)))
		buf.WriteByte(' ')
	}
	return buf.String()
}

func (l *logging) formatValue(o interface{}) string {
	if o == nil {
		return ""
	}

	if t, ok := o.(time.Time); ok {
		if l.timeFormatter != nil {
			return l.timeFormatter(t)
		}
	}
	return util.FormatValue(o, false)
}

func (l *logging) LogF(level Level, depth int, keyValues util.KeyValues, format string, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	length := len(format)
	if length > 0 {
		if format[length-1] != '\n' {
			format = format + "\n"
		}
	}
	logInfo := fmt.Sprintf(format, args...)
	w := l.selectWriter(level)
	l.format(w, level, depth, keyValues, logInfo)

	if level == PANIC {
		l.panicFunc(util.NewKeyValues(ContentKey, logInfo))
	} else if level <= FATAL {
		l.processFatal(w)
	}

	//l.output(level, buf)
}

func (l *logging) Log(level Level, depth int, keyValues util.KeyValues, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	logInfo := fmt.Sprint(args...)
	w := l.selectWriter(level)
	l.format(w, level, depth, keyValues, logInfo)

	if level == PANIC {
		l.panicFunc(util.NewKeyValues(ContentKey, logInfo))
	} else if level <= FATAL {
		l.processFatal(w)
	}
}

func (l *logging) LogLn(level Level, depth int, keyValues util.KeyValues, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	logInfo := fmt.Sprintln(args...)
	w := l.selectWriter(level)
	l.format(w, level, depth, keyValues, logInfo)

	if level == PANIC {
		l.panicFunc(util.NewKeyValues(ContentKey, logInfo))
	} else if level <= FATAL {
		l.processFatal(w)
	}
}

func (l *logging) getBuffer() *bytes.Buffer {
	buf := l.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (l *logging) putBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	if buf.Len() > 256 {
		//let big buffers die a natural death.
		return
	}
	l.bufPool.Put(buf)
}

func (l *logging) processFatal(writer io.Writer) {
	if !l.fatalNoTrace {
		trace := stacks(true)
		writer.Write(trace)
	}
	l.exitFunc(-1)
}

func (l *logging) Clone() Logging {
	ret := &logging{
		timeFormatter:   l.timeFormatter,
		callerFormatter: l.callerFormatter,
		//formatter:     l.formatter,
		colorFlag:    l.colorFlag,
		fileFlag:     l.fileFlag,
		fatalNoTrace: l.fatalNoTrace,
		level:        l.level,
		//writers:       map[Level]io.Writer{},

		bufPool: sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(nil)
		}},
	}
	ret.formatter.Store(l.formatter.Load())
	l.writers.Range(func(key, value interface{}) bool {
		ret.writers.Store(key, value)
		return true
	})
	return ret
}

//func (l *logging) output(level Level) {
//	if level >= FATAL {
//		if !l.fatalNoTrace {
//			trace := stacks(true)
//			buf.WriteString("\n")
//			buf.Write(trace)
//		}
//		l.selectWriter(level).Write(buf.Bytes())
//		os.Exit(-1)
//	} else {
//		l.selectWriter(level).Write(buf.Bytes())
//	}
//	l.putBuffer(buf)
//}

func (l *logging) selectWriter(level Level) io.Writer {
	for i := level; i <= DEBUG; i++ {
		v, ok := l.writers.Load(i)
		if ok && v != nil {
			return v.(io.Writer)
		}
	}
	return os.Stdout
}

func (l *logging) SetFormatter(f util.Formatter) {
	l.formatter.Store(f)
}

func (l *logging) GetFormatter() util.Formatter {
	v := l.formatter.Load()
	if v == nil {
		return nil
	}
	return v.(util.Formatter)
}

func (l *logging) SetLogLevel(severity Level) {
	atomic.StoreInt32(&l.level, severity)
}

func (l *logging) IsEnabled(severityLevel Level) bool {
	return atomic.LoadInt32(&l.level) >= severityLevel
}

// Logging不会自动为输出的Writer加锁，如果需要加锁请使用LockedWriter：
// logging.SetOutPut(&writer.LockedWriter{w})
func (l *logging) SetOutput(w io.Writer) {
	for i := FATAL; i <= DEBUG; i++ {
		l.writers.Store(i, w)
	}
}

// Logging不会自动为输出的Writer加锁，如果需要加锁请使用LockedWriter：
// logging.SetOutputBySeverity(level, &writer.LockedWriter{w})
func (l *logging) SetOutputBySeverity(severityLevel Level, w io.Writer) {
	l.writers.Store(severityLevel, w)
}

func (l *logging) GetOutputBySeverity(severityLevel Level) io.Writer {
	v, ok := l.writers.Load(severityLevel)
	if !ok {
		return nil
	}
	return v.(io.Writer)
}

func selectLevelColor(level Level) string {
	if level == INFO {
		return ForeGreen
	} else if level == WARN {
		return ForeYellow
	} else if level == DEBUG {
		return ForeCyan
	} else if level < WARN {
		return ForeRed
	}
	return ""
}

func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func stacks(all bool) []byte {
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
	}
	return trace
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006-01-02 15:04:05")
	return timeString
}

func callerFormat(file string, line int, funcName string) string {
	//no file
	if file != "" {
		if funcName == "" {
			return fmt.Sprintf("%s:%d", file, line)
		} else {
			return fmt.Sprintf("%s:%d (%s)", file, line, funcName)
		}
	} else {
		if funcName == "" {
			return ""
		} else {
			return "(" + funcName + ")"
		}
	}
}

// SetTimeFormatter 配置内置Logging实现的时间格式化函数
func SetTimeFormatter(f func(t time.Time) string) func(*logging) {
	return func(logging *logging) {
		logging.timeFormatter = f
	}
}

// SetCallerFormatter 配置内置Logging实现的时间格式化函数
func SetCallerFormatter(f func(file string, line int, funcName string) string) func(*logging) {
	return func(logging *logging) {
		logging.callerFormatter = f
	}
}

// SetExitFunc 配置内置Logging Fatal退出处理函数
func SetExitFunc(f ExitFunc) func(*logging) {
	return func(logging *logging) {
		logging.exitFunc = f
	}
}

// SetPanicFunc 配置内置Logging Panic处理函数
func SetPanicFunc(f PanicFunc) func(*logging) {
	return func(logging *logging) {
		logging.panicFunc = f
	}
}

func defaultExit(code int) {
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		os.Exit(code)
	} else {
		if err = p.Signal(os.Interrupt); err != nil {
			// Per https://golang.org/pkg/os/#Signal, “Interrupt is not implemented on
			// Windows; using it with os.Process.Signal will return an error.”
			// Fall back to Kill instead.
			if err = p.Signal(syscall.SIGTERM); err != nil {
				if err = p.Kill(); err != nil {
					os.Exit(code)
				}
			}
		}
	}
}

func defaultPanic(v interface{}) {
	panic(v)
}

func simpleFuncName(funcName string) string {
	segs := strings.Split(funcName, "/")
	buf := strings.Builder{}
	buf.Grow(len(funcName) / 2)
	size := len(segs) - 1
	for i := 0; i < size; i++ {
		if len(segs[i]) > 0 {
			buf.WriteString(segs[i][:1])
			buf.WriteString(".")
		}
	}
	buf.WriteString(segs[size])
	return buf.String()
}

// SetColorFlag 配置内置Logging实现的颜色的标志，有AutoColor、DisableColor、ForceColor
func SetColorFlag(flag int) func(*logging) {
	return func(logging *logging) {
		logging.colorFlag = flag
	}
}

// SetLogLevel
func SetLogLevel(level Level) func(*logging) {
	return func(logging *logging) {
		logging.level = level
	}
}

// SetCallerFlag 配置内置Logging实现的文件输出标志，有ShortFile、LongFile
func SetCallerFlag(flag int) func(*logging) {
	return func(logging *logging) {
		logging.fileFlag = flag
	}
}

// SetFatalNoTrace 配置内置Logging实现是否在发生致命错误时打印堆栈，默认打印
func SetFatalNoTrace(noTrace bool) func(*logging) {
	return func(logging *logging) {
		logging.fatalNoTrace = noTrace
	}
}

const autogeneratedFrameName = "<autogenerated>"

func FramesToCaller() int {
	for i := 1; i < 3; i++ {
		_, file, _, _ := runtime.Caller(i + 1)
		if file != autogeneratedFrameName {
			return i
		}
	}
	return 1
}
