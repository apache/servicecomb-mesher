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

package health

import (
	"context"
	"errors"
	"github.com/apache/servicecomb-mesher/proxy/config"
	"github.com/go-chassis/foundation/httpclient"
	"github.com/go-mesh/openlogging"
	"io/ioutil"
	"regexp"
	"strconv"
)

//HTTPCheck checks http service
func HTTPCheck(check *config.HealthCheck, address string) error {
	c, err := httpclient.New(&httpclient.DefaultOptions)
	if err != nil {
		openlogging.Error("can not get http client: " + err.Error())
		//must not return error, because it is mesher error
		return nil
	}
	var url = "http://" + address
	if check.URI != "" {
		url = url + check.URI
	}
	resp, err := c.Get(context.Background(), url, nil)
	if err != nil {
		openlogging.Error("server can not be connected: " + err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if check.Match != nil {
		if check.Match.Status != "" {
			n, _ := strconv.Atoi(check.Match.Status)
			if resp.StatusCode != n {
				return errors.New("status is not " + check.Match.Status)
			}
		}
		if check.Match.Body != "" {
			re := regexp.MustCompile(check.Match.Body)
			if !re.Match(body) {
				return errors.New("body does not match " + check.Match.Body)
			}
		}
	} else {
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return nil
}
