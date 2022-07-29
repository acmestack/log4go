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
	"encoding/base64"
	"fmt"
	"github.com/acmestack/log4go/ext"
	"github.com/acmestack/log4go/log"
	"github.com/acmestack/log4go/logfactory"
	"github.com/acmestack/log4go/util"
	"github.com/acmestack/log4go/writer"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncWriter(t *testing.T) {
	w := writer.NewAsyncWriter(os.Stdout, nil, 10, true)
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				atomic.AddInt32(&count, 1)
				rand.Read(b)
				_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	wait.Wait()
	t.Log(count)
	w.Close()
}

func TestAsyncBufWriter(t *testing.T) {
	w := writer.NewAsyncBufferWriter(os.Stdout, nil, writer.Config{
		FlushSize:     100,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				atomic.AddInt32(&count, 1)
				rand.Read(b)
				_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestRotateFile(t *testing.T) {
	w := writer.NewRotateFileWriter(&writer.RotateFile{
		Path: "./target/test.log",
	}, writer.Config{
		FlushSize:     100,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				atomic.AddInt32(&count, 1)
				rand.Read(b)
				_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestRotateFilePart(t *testing.T) {
	w := writer.NewRotateFileWriter(&writer.RotateFile{
		Path:        "./target/test.log",
		MaxFileSize: 10,
	}, writer.Config{
		FlushSize:     100,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				atomic.AddInt32(&count, 1)
				rand.Read(b)
				_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestRotateFilePartAndTime(t *testing.T) {
	w := writer.NewRotateFileWriter(&writer.RotateFile{
		Path:            "./target/test.log",
		MaxFileSize:     60,
		RotateFrequency: 1 * writer.RotateEverySecond,
	}, writer.Config{
		FlushSize:     10,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})

	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for j := 0; j < 10; j++ {
				time.Sleep(300 * time.Millisecond)
				atomic.AddInt32(&count, 1)
				rand.Read(b)
				_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestRotateFilePartAndTimeWithZip(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		w := writer.NewRotateFileWriter(&writer.RotateFile{
			Path:            "./target/test.log",
			MaxFileSize:     1000,
			RotateFrequency: writer.RotateEveryMinute,
			RotateFunc:      writer.ZipLogsAsync,
		}, writer.Config{
			FlushSize:     1000,
			BufferSize:    1024,
			FlushInterval: 50 * time.Millisecond,
			Block:         false,
		})

		for i := 0; i < 300; i++ {
			time.Sleep(300 * time.Millisecond)
			_, err := w.Write([]byte(fmt.Sprintf("[%d][%s]\n", i, time.Now().Format("2006-01-02-15-04-05"))))
			if err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(time.Second)
		w.Close()
	})

	t.Run("fast rotate", func(t *testing.T) {
		w := writer.NewRotateFileWriter(&writer.RotateFile{
			Path:            "./target/test.log",
			MaxFileSize:     10,
			RotateFrequency: writer.RotateEverySecond,
			RotateFunc:      writer.ZipLogsAsync,
		}, writer.Config{
			FlushSize:     100,
			BufferSize:    10,
			FlushInterval: 1 * time.Millisecond,
			Block:         true,
		})

		var count int32 = 0
		wait := sync.WaitGroup{}
		wait.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				defer wait.Done()
				b := make([]byte, 10)
				for j := 0; j < 10; j++ {
					time.Sleep(300 * time.Millisecond)
					atomic.AddInt32(&count, 1)
					rand.Read(b)
					_, err := w.Write([]byte(strconv.Itoa(int(count)) + base64.StdEncoding.EncodeToString(b) + "\n"))
					if err != nil {
						t.Fatal(err)
					}
				}
			}()
		}

		wait.Wait()
		w.Close()
		t.Log(count)
	})
}

func TestBufferedRotateFileWriterePartAndTimeWithZip(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		w := writer.NewBufferedRotateFileWriter(&writer.BufferedRotateFile{
			Path:            "./target/test.log",
			MaxFileSize:     1000,
			RotateFrequency: writer.RotateEveryMinute,
			RotateFunc:      writer.ZipLogsAsync,
		}, writer.Config{
			FlushSize:     1000,
			BufferSize:    1024,
			FlushInterval: 50 * time.Millisecond,
			Block:         false,
		})

		for i := 0; i < 300; i++ {
			time.Sleep(300 * time.Millisecond)
			_, err := w.Write([]byte(fmt.Sprintf("[%d][%s]\n", i, time.Now().Format("2006-01-02-15-04-05"))))
			if err != nil {
				t.Fatal(err)
			}
		}
		w.Close()
	})

	t.Run("file test", func(t *testing.T) {
		w := writer.NewBufferedRotateFileWriter(&writer.BufferedRotateFile{
			Path:            "./target/test.log",
			MaxFileSize:     10,
			RotateFrequency: writer.RotateEverySecond,
			RotateFunc:      writer.ZipLogsAsync,
		}, writer.Config{
			FlushSize:     100,
			BufferSize:    0,
			FlushInterval: 1 * time.Second,
			Block:         true,
		})

		var count int32 = 0
		wait := sync.WaitGroup{}
		wait.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				defer wait.Done()
				b := make([]byte, 10)
				for j := 0; j < 10; j++ {
					time.Sleep(300 * time.Millisecond)
					v := atomic.AddInt32(&count, 1)
					rand.Read(b)
					_, err := w.Write([]byte(fmt.Sprintf("[%d][%s][%s]\n", v, time.Now().Format("2006-01-02-15-04-05"), base64.StdEncoding.EncodeToString(b))))
					if err != nil {
						t.Fatal(err)
					}
				}
			}()
		}

		wait.Wait()
		w.Close()
		t.Log(count)
	})
}

func TestMultiWriterLog(t *testing.T) {
	w := writer.NewRotateFileWriter(&writer.RotateFile{
		Path: "./target/test.log",
	}, writer.Config{
		FlushSize:     100,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})
	log.DefaultLogging().SetOutput(io.MultiWriter(os.Stdout, w))
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				x := atomic.AddInt32(&count, 1)
				rand.Read(b)
				log.InfoLn(strconv.Itoa(int(x)) + "-" + base64.StdEncoding.EncodeToString(b))
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestMultiWriterMutableLogger(t *testing.T) {
	w := writer.NewRotateFileWriter(&writer.RotateFile{
		Path: "./target/test.log",
	}, writer.Config{
		FlushSize:     100,
		BufferSize:    10,
		FlushInterval: 1 * time.Millisecond,
		Block:         true,
	})

	logging := logfactory.NewLogging(logfactory.SetColorFlag(logfactory.DisableColor))
	logging.SetOutput(io.MultiWriter(os.Stdout, w))
	logfactory.ResetFactory(ext.NewMutableFactory(logging))
	logger := logfactory.GetLogger()
	var count int32 = 0
	wait := sync.WaitGroup{}
	wait.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			b := make([]byte, 10)
			for i := 0; i < 10; i++ {
				x := atomic.AddInt32(&count, 1)
				rand.Read(b)
				logger.InfoLn(strconv.Itoa(int(x)) + "-" + base64.StdEncoding.EncodeToString(b))
			}
		}()
	}

	wait.Wait()
	w.Close()
	t.Log(count)
}

func TestValue(t *testing.T) {
	v := util.NewAtomicValue(os.Stdout)
	t.Log(v.Load())
	v.Store(writer.NewAsyncBufferWriter(os.Stdout, nil, writer.Config{}))
	t.Log(v.Load())
}
