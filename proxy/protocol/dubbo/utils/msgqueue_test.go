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
	"container/list"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func NewMsgQueueForTest(maxMsgNum int) *MsgQueue {
	q := new(MsgQueue)
	q.msgList = list.New()
	q.mtx = new(sync.Mutex)
	q.msgCount = 0
	q.maxMsgNum = maxMsgNum
	q.state = Actived
	q.notEmptyCond = sync.NewCond(q.mtx)
	q.notFullCond = sync.NewCond(q.mtx)
	return q
}

func TestMsgQueue(t *testing.T) {
	t.Run("case empty", func(t *testing.T) {
		q := NewMsgQueueForTest(2)

		// append msg
		eMSG := "msg to send"
		err := q.Enqueue(eMSG)
		assert.NoError(t, err)

		assert.Equal(t, false, q.isEmpty())
		assert.Equal(t, false, q.isFull())

		// pop msg
		dMSG, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, eMSG, dMSG)

		// case empty
		done := make(chan struct{})
		go func(done chan struct{}) {
			dMSG, err = q.Dequeue()
			assert.NoError(t, err)

			// Deavtive
			q.Deavtive()

			// error
			_, err = q.Dequeue()
			assert.Error(t, err)

			if done != nil {
				close(done)
			}

		}(done)

		select {
		case <-time.After(time.Second * 10):
		case <-done:
		}

	})

	t.Run("case Full", func(t *testing.T) {
		eMSG := "msg to send"
		// case Full
		q1 := NewMsgQueueForTest(2)
		err := q1.Enqueue(eMSG)
		assert.NoError(t, err)

		err = q1.Enqueue(eMSG)
		assert.NoError(t, err)
		done1 := make(chan struct{})
		go func(c chan struct{}) {
			err = q1.Enqueue(eMSG)
			t.Log(err)
			assert.Error(t, err)

			if c != nil {
				close(c)
			}

		}(done1)

		select {
		case <-time.After(time.Second * 10):
		case <-done1:
		}

		// Deavtive
		q1.Deavtive()

		// error
		err = q1.Enqueue(eMSG)
		assert.Error(t, err)
	})
}
