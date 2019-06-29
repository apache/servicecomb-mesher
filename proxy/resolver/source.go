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

package resolver

import (
	"errors"
	"github.com/go-chassis/go-chassis/core/registry"
)

var (
	//ErrFoo is of type error
	ErrFoo = errors.New("resolved as a nil service")
)

//SourceResolver is a interface which has Resolve function
type SourceResolver interface {
	Resolve(source string) *registry.SourceInfo
}

var sr SourceResolver = &DefaultSourceResolver{}

//DefaultSourceResolver is a struct
type DefaultSourceResolver struct {
}

//Resolve is a method which resolves service endpoint
func (sr *DefaultSourceResolver) Resolve(source string) *registry.SourceInfo {
	if source == "127.0.0.1" {
		return nil
	}
	si := registry.GetIPIndex(source)

	return si
}

//GetSourceResolver returns interface object
func GetSourceResolver() SourceResolver {
	return sr
}
