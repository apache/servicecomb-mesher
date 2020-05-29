/*
* Licensed to the Apache Software Foundation (ASF) under one or more
* contributor license agreements.  See the NOTICE file distributed with
* this work for additional information regarding copyright ownership.
* The ASF licenses this file to You under the Apache License, Version 2.0
* (the "License"); you may not use this file except in compliance with
* the License.  You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewThreadGroupWait(t *testing.T) {
	count := 0
	var done chan struct{}
	go func(done chan struct{}) {
		var tgw *ThreadGroupWait
		tgw = NewThreadGroupWait()
		tgw.Add(1)
		go func(tgw *ThreadGroupWait) {
			defer tgw.Done()
			count++
		}(tgw)

		tgw.Wait()
		close(done)
	}(done)

	select {
	case <-time.After(time.Second * 5):
	case <-done:
	}
	assert.Equal(t, 1, count)

	// case done count < 0
	tgw := NewThreadGroupWait()
	tgw.Done()
	tgw.Done()
}

type Task struct {
}

func (t *Task) Svc(interface{}) interface{} {
	return nil
}

func TestThrmgr(t *testing.T) {
	nr := NewRoutineManager()
	var done chan struct{}
	go func(done chan struct{}) {
		nr.Wait()
		close(done)
	}(done)

	time.AfterFunc(time.Second*2, func() {
		nr.Done()
	})

	select {
	case <-time.After(time.Second * 5):
	case <-done:
	}

	nr.Spawn(new(Task), "swap", "routinename")

}
