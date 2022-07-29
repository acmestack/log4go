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
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotateFrequency = time.Duration

const (
	// 仅使用一个文件记录日志
	RotateNone RotateFrequency = 0
	// 每天凌晨滚动文件
	RotateEveryDay RotateFrequency = time.Hour * 24
	// 每小时整点时滚动文件
	RotateEveryHour RotateFrequency = time.Hour
	// WARNING: for test only!
	RotateEveryMinute RotateFrequency = time.Minute
	RotateEverySecond RotateFrequency = time.Second
)

func NewRotateFileWriter(f *RotateFile, conf ...Config) io.WriteCloser {
	if f == nil {
		return nil
	}

	err := f.Open()
	if err != nil {
		return nil
	}

	return NewAsyncBufferWriter(f, f.Close, conf...)
}

type RotateFile struct {
	//文件路径
	Path string
	// 文件的大小阈值
	MaxFileSize int64
	// 滚动频率
	RotateFrequency RotateFrequency
	// 滚动文件处理
	RotateFunc func(dir string, name string, files ...string) error

	// 滚动的时间格式
	timeFormat string
	// 滚动的时间间隔
	rotateDuration time.Duration

	timer      *time.Timer
	fileName   string
	dir        string
	file       *os.File
	curSize    int64
	part       int
	curTimeStr string
}

func (f *RotateFile) Open() error {
	if f.MaxFileSize == 0 {
		// no limit
		f.MaxFileSize = math.MaxInt64
	}
	if f.timeFormat == "" {
		f.timeFormat = "2006-01-02"
	}
	dir := filepath.Dir(f.Path)
	_, err := os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f.dir = dir
	f.fileName = filepath.Base(f.Path)

	f.file, err = os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := f.file.Stat()
	if err != nil {
		return err
	}
	f.curSize = info.Size()
	if f.RotateFrequency != RotateNone {
		f.setFrequency(f.RotateFrequency)
		f.setTimer()
	}
	return f.calcPart()
}

func (f *RotateFile) setTimer() {
	f.curTimeStr = time.Now().Format(f.timeFormat)
	t := f.nextTime()
	duration := t.Sub(time.Now())
	if duration < 0 {
		duration = 1
	}
	f.timer = time.NewTimer(duration)
}

func (f *RotateFile) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if f.timer != nil {
		select {
		case <-f.timer.C:
			err := f.rotateByTime()
			f.setTimer()
			if err != nil {
				return 0, err
			}
		default:
		}
	}
	n, err := f.file.Write(data)
	f.curSize += int64(n)
	if err != nil {
		return n, err
	}
	if f.curSize >= f.MaxFileSize {
		err := f.rotatePart()
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func (f *RotateFile) rotateByTime() error {
	//err := f.changeFile(fmt.Sprintf("%s-%s", f.curTimeStr, f.fileName))
	err := f.rotatePart()
	if err != nil {
		return err
	}
	if f.RotateFunc != nil {
		oldTimeStr := f.curTimeStr
		files, err := ioutil.ReadDir(f.dir)
		if err != nil {
			return err
		}
		var partsFiles []string
		for _, v := range files {
			i := strings.Index(v.Name(), oldTimeStr)
			if i != -1 {
				partsFiles = append(partsFiles, filepath.Join(f.dir, v.Name()))
			}
		}
		if len(partsFiles) > 0 {
			err = f.RotateFunc(f.dir, oldTimeStr+"-"+f.fileName, partsFiles...)
			if err != nil {
				return err
			}
		}
	}

	f.curTimeStr = time.Now().Format(f.timeFormat)
	f.part = 0
	return nil
}

func (f *RotateFile) calcPart() error {
	files, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return err
	}
	part := 0
	prefix := ""
	if f.curTimeStr == "" {
		prefix = "part"
	} else {
		prefix = f.curTimeStr + "-part"
	}
	for _, v := range files {
		i := strings.Index(v.Name(), prefix)
		if i != -1 {
			part++
		}
	}
	f.part = part
	return nil
}

func (f *RotateFile) rotatePart() error {
	if f.curTimeStr == "" {
		err := f.changeFile(fmt.Sprintf("part%d-%s", f.part, f.fileName))
		f.part++
		return err
	} else {
		err := f.changeFile(fmt.Sprintf("%s-part%d-%s", f.curTimeStr, f.part, f.fileName))
		f.part++
		return err
	}
}

func (f *RotateFile) changeFile(filename string) error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(f.dir, f.fileName), filepath.Join(f.dir, filename))
	if err != nil {
		return err
	}
	f.file, err = os.OpenFile(filepath.Join(f.dir, f.fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	f.curSize = 0
	return nil
}

func (f *RotateFile) nextTime() time.Time {
	timeStr := time.Now().Format(f.timeFormat)
	t, _ := time.ParseInLocation(f.timeFormat, timeStr, time.Local)
	return t.Add(f.rotateDuration)
}

func (f *RotateFile) Close() error {
	if f.timer != nil {
		f.timer.Stop()
	}
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *RotateFile) setFrequency(frequency RotateFrequency) {
	interval := frequency / RotateEveryDay
	if interval > 0 {
		f.rotateDuration = interval * RotateEveryDay
		f.timeFormat = "2006-01-02"
		return
	}

	interval = frequency / RotateEveryHour
	if interval > 0 {
		f.rotateDuration = interval * RotateEveryHour
		f.timeFormat = "2006-01-02-15"
		return
	}

	interval = frequency / RotateEveryMinute
	if interval > 0 {
		f.rotateDuration = interval * RotateEveryMinute
		f.timeFormat = "2006-01-02-15-04"
		return
	}

	interval = frequency / RotateEverySecond
	if interval > 0 {
		f.rotateDuration = interval * RotateEverySecond
		f.timeFormat = "2006-01-02-15-04-05"
		return
	}
}

func ZipLogsAsync(dir string, name string, files ...string) error {
	go ZipLogs(dir, name, files...)
	return nil
}

func ZipLogs(dir string, name string, files ...string) error {
	if len(files) == 0 {
		return nil
	}
	zipFile := filepath.Join(dir, name+".zip")
	f, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	for _, v := range files {
		err := compress(v, w)
		if err != nil {
			return err
		} else {
			err = os.Remove(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func compress(file string, w *zip.Writer) error {
	of, err := os.Open(file)
	if err != nil {
		return err
	}
	defer of.Close()

	info, err := of.Stat()
	if err != nil {
		return err
	}
	if !info.IsDir() {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		wh, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(wh, of)
		if err != nil {
			return err
		}
	}
	return nil
}
