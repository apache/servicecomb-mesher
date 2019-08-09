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

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func handlerGateway(w http.ResponseWriter, r *http.Request) {

	queryInfo, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("parse queryInfo wrong: ", queryInfo)
		fmt.Fprintln(w, "parse queryInfo wrong")
		return
	}
	fmt.Println("height:%s,weight:%s", queryInfo.Get("height"), queryInfo.Get("weight"))
	strReqUrl := "http://mersher-provider:4540/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	resp, err := http.Get(strReqUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("body err: ", body, err)
		return
	}
	fmt.Println("body: " + string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)

}

func handlerGateway2(w http.ResponseWriter, r *http.Request) {
	queryInfo, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("parse queryInfo wrong: ", queryInfo)
		fmt.Fprintln(w, "parse queryInfo wrong")
		return
	}
	fmt.Println("height:%s,weight:%s", queryInfo.Get("height"), queryInfo.Get("weight"))
	strReqUrl := "http://mersher-provider/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	//strReqUrl := "http://mersher-ht-provider/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	//strReqUrl := "http://mersher-ht-provider:4555/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	proxy, _ := url.Parse("http://127.0.0.1:30101") //将mesher设置为http代理
	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	req, err := http.NewRequest(http.MethodGet, strReqUrl, nil)
	if err != nil {
		fmt.Println("make req err: ", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Do err: ", resp, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("body err: ", body, err)
		return
	}
	fmt.Println("body: " + string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

func main() {
	http.HandleFunc("/bmi", handlerGateway2)
	http.ListenAndServe("192.168.88.64:4538", nil)
}
