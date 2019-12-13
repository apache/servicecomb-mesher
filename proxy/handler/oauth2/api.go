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

package oauth2

import (
	"golang.org/x/oauth2"
	"net/http"
)

var auth *OAuth2

// OAuth2 should implement oauth2 server side logic
// it is singleton
type OAuth2 struct {
	GrantType    string                                            // required
	UseConfig    *oauth2.Config                                    // required
	Authenticate func(accessToken string, req *http.Request) error // optional
}

// Use put a custom oauth2 logic
// then register handler to chassis
func Use(middleware *OAuth2) {
	auth = middleware
}
