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
	"sync"
)

const (
	MaxBufferMsg = 65535
	Actived      = 0
	Deactived    = 1
)

//MsgQueue thread safe queue
type MsgQueue struct {
	msgList      *list.List
	mtx          *sync.Mutex
	msgCount     int
	maxMsgNum    int
	state        int
	notEmptyCond *sync.Cond
	notFullCond  *sync.Cond
}

//NewMsgQueue is a function which initializes msgqueue value
func NewMsgQueue() *MsgQueue {
	tmp := new(MsgQueue)
	tmp.msgList = list.New()
	tmp.mtx = new(sync.Mutex)
	tmp.msgCount = 0
	tmp.maxMsgNum = MaxBufferMsg
	tmp.state = Actived
	tmp.notEmptyCond = sync.NewCond(tmp.mtx)
	tmp.notFullCond = sync.NewCond(tmp.mtx)
	return tmp
}

//Enqueue is method which enqueues message in queue
func (this *MsgQueue) Enqueue(msg interface{}) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.state == Deactived {
		return &BaseError{"Queue is deactive"}
	}
	if this.waitNotFullCond() == -1 {
		return &BaseError{"Enqueue time out"}
	}
	this.msgList.PushFront(msg)
	this.msgCount++
	this.notEmptyCond.Signal()
	return nil
}

//Dequeue is a method which dequeues message from queue
func (this *MsgQueue) Dequeue() (interface{}, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.waitNotEmptyCond() == -1 {
		return nil, &BaseError{"Queue is deactive"}
	}

	iter := this.msgList.Back()
	v := iter.Value
	this.msgList.Remove(iter)
	this.msgCount--
	this.notFullCond.Signal()
	return v, nil
}

//isEmpty is a method which checks whether queue is empty
func (this *MsgQueue) isEmpty() bool {
	if this.msgCount == 0 {
		return true
	} else {
		return false
	}
}

//isFull is a method which checks whether queue is full
func (this *MsgQueue) isFull() bool {
	if this.msgCount >= this.maxMsgNum {
		return true
	} else {
		return false
	}
}

//waitNotFullCond is a method which waits if queue is full
func (this *MsgQueue) waitNotFullCond() int {
	var result = 0

	if this.isFull() {
		this.notFullCond.Wait()
		if this.state != Actived {
			result = -1
			return result
		}
	}
	return result
}

//Deavtive is a method
func (this *MsgQueue) Deavtive() {
	this.state = Deactived
	this.notEmptyCond.Broadcast()
	this.notFullCond.Broadcast()
}

//waitNotEmptyCond is a method
func (this *MsgQueue) waitNotEmptyCond() int {
	var result = 0

	if this.isEmpty() {
		this.notEmptyCond.Wait()
		if this.state != Actived {
			result = -1
			return result
		}
	}
	return result
}
