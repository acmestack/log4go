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

package util

import (
	"sync/atomic"
)

// Value 存储值对象工具，interface不做类型检查，需用户自行确认存取类型
type Value interface {
	// Store 存储值
	Store(interface{})
	// Load 取出值
	Load() interface{}
}

type SimpleValue struct {
	o interface{}
}

func NewSimpleValue(o interface{}) *SimpleValue {
	return &SimpleValue{o: o}
}

func (l *SimpleValue) Load() interface{} {
	return l.o
}

func (l *SimpleValue) Store(o interface{}) {
	l.o = o
}

type AtomicValue atomic.Value

func NewAtomicValue(o interface{}) *AtomicValue {
	ret := &AtomicValue{}
	ret.Store(o)
	return ret
}
func (l *AtomicValue) Load() interface{} {
	return (*atomic.Value)(l).Load().(*SimpleValue).Load()
}

func (l *AtomicValue) Store(o interface{}) {
	(*atomic.Value)(l).Store(&SimpleValue{o: o})
}
