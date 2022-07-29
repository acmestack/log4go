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
	"testing"
	"time"
)

func TestRotateTime(t *testing.T) {
	f := RotateFile{}
	f.setFrequency(40 * RotateEveryDay)
	t.Log("40 day")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(RotateEveryDay)
	t.Log("1 day")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(25 * RotateEveryHour)
	t.Log("25 hour")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(RotateEveryHour)
	t.Log("1 hour")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(70 * RotateEveryMinute)
	t.Log("70 minute")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(RotateEveryMinute)
	t.Log("1 minute")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(70 * RotateEverySecond)
	t.Log("70 second")
	t.Log(time.Now())
	t.Log(f.nextTime())

	f.setFrequency(RotateEverySecond)
	t.Log("1 second")
	t.Log(time.Now())
	t.Log(f.nextTime())
}
