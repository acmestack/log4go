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

package writer

import (
	"io"
	"sync"
)

// Logging不会自动为输出的Writer加锁，如果需要加锁请使用这个封装工具：
// logging.SetOutPut(&writer.LockedWriter{w})
type LockedWriter struct {
	lock sync.Mutex
	W    io.Writer
}

func (lw *LockedWriter) Write(d []byte) (int, error) {
	lw.lock.Lock()
	defer lw.lock.Unlock()

	return lw.W.Write(d)
}

type LockedWriteCloser struct {
	lock sync.Mutex
	W    io.WriteCloser
}

func (lw *LockedWriteCloser) Write(d []byte) (int, error) {
	lw.lock.Lock()
	defer lw.lock.Unlock()

	return lw.W.Write(d)
}

func (lw *LockedWriteCloser) Close() error {
	lw.lock.Lock()
	defer lw.lock.Unlock()

	return lw.W.Close()
}
