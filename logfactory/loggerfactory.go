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

import "github.com/acmestack/log4go/util"

type LoggerFactoryI interface {
	// GetLogger 根据参数获得Logger
	// Param：根据默认实现，o可不填，直接返回一个没有名称的Logger。
	// 如果o有值，则只取第一个值，且当：
	// 		o为string时，使用string值作为Logger名称
	//		o为其他类型时，取package path + type name作为Logger名称，以"."分隔，如g.x.x.t.TestStructInTest
	GetLogger(o ...interface{}) Logger

	// Reset 重置Factory的Logging（线程安全）
	Reset(logging Logging) LoggerFactoryI

	// GetLogging 获得Factory的Logging（线程安全），可用来配置Logging
	// 也可以通过wrap Logging达到控制日志级别、日志输出格式的目的
	GetLogging() Logging
}

type LoggerFactory struct {
	Value            util.Value
	SimplifyNameFunc func(string) string
}

var defaultFactory util.Value = util.NewSimpleValue(NewFactory(DefaultLogging()))

func NewDefaultFactory(opts ...LoggingOpt) *LoggerFactory {
	return NewFactory(NewLogging(opts...))
}

func NewFactory(logging Logging) *LoggerFactory {
	return NewFactoryWithValue(util.NewAtomicValue(logging))
}

func NewFactoryWithValue(v util.Value) *LoggerFactory {
	ret := &LoggerFactory{
		Value: v,
	}
	return ret
}

func (fac *LoggerFactory) GetLogging() Logging {
	return fac.Value.Load().(Logging)
}

func (fac *LoggerFactory) GetLogger(o ...interface{}) Logger {
	name := util.GetObjectName(fac.SimplifyNameFunc, o...)
	return defaultLogger(fac.Value.Load().(Logging), nil, name)
}

func (fac *LoggerFactory) Reset(logging Logging) LoggerFactoryI {
	fac.Value.Store(logging)
	return fac
}

// ResetFactory 重新配置全局的默认LoggerFactory，该方法同时会重置全局的默认Logging
// 由于线程安全性受defaultLogging、defaultFactory初始化（调用InitOnce）的Value决定，
// 所以需要确定是否确实需要调用该方法重置Logging，并保证Value线程安全
func ResetFactory(fac LoggerFactoryI) {
	defaultFactory.Store(fac)
	ResetLogging(fac.GetLogging())
}

//// ResetLogging 重新配置全局的默认Logging，该方法同时会重置全局的默认LoggerFactory的Logging
//// 由于线程安全性受defaultLogging、defaultFactory初始化（调用InitOnce）的Value决定，
//// 所以需要确定是否确实需要调用该方法重置Logging，并保证Value线程安全
func ResetLogging(logging Logging) {
	GetLogging().Store(logging)
	defaultFactory.Load().(LoggerFactoryI).Reset(GetLogging().Load().(Logging))
}

// GetLogger 通过全局默认LoggerFactory获取Logger
// Param：根据默认实现，o可不填，直接返回一个没有名称的Logger。
// 如果o有值，则只取第一个值，且当：
// 		o为string时，使用string值作为Logger名称
//		o为其他类型时，取package path + type name作为Logger名称，以"."分隔，如g.x.x.t.TestStructInTest
func GetLogger(o ...interface{}) Logger {
	return defaultFactory.Load().(LoggerFactoryI).GetLogger(o...)
}
