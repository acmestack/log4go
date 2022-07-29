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
	"runtime"
	"sync"
)

// frameMap 用于缓存调用点，基准测试表明使用缓存大约有 50% 的性能提升。
var frameMap sync.Map

// Caller 获取调用点的文件及行号信息，fast 为 true 时使用缓存进行加速。
func Caller(skip int, fast bool) (pc uintptr, file string, line int, ok bool) {

	if !fast {
		pc, file, line, ok = runtime.Caller(skip + 1)
		return
	}

	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+2, rpc[:])
	if n < 1 {
		return
	}
	pc2 := rpc[0]
	if v, ok := frameMap.Load(pc2); ok {
		e := v.(*runtime.Frame)
		return pc2, e.File, e.Line, true
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	frameMap.Store(pc2, &frame)
	return pc2, frame.File, frame.Line, true
}
