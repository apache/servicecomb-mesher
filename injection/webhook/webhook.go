/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package webhook

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apache/servicecomb-mesher/injection/templates"
	"github.com/go-mesh/openlogging"
	"k8s.io/api/admission/v1beta1"
)

// Webhook struct
type Webhook struct {
	server    *http.Server
	config    *config
	templater templates.Templater
}

// NewWebhook new webhook server
func NewWebhook(ops ...ConfigOption) (*Webhook, error) {
	conf := toConfig(ops...)
	templater, err := templates.NewTemplater(conf.templateName, conf.sidecarConfig, conf.sidecarTemplate)
	if err != nil {
		openlogging.Error(fmt.Sprintf("new templater failed: error = %v, template = %s", err.Error(), conf.templateName))
		return nil, err
	}

	return &Webhook{
		server: &http.Server{
			Addr:      fmt.Sprintf(":%d", conf.port),
			TLSConfig: &tls.Config{GetCertificate: conf.getCertificate},
		},
		config:    conf,
		templater: templater,
	}, nil
}

// Run webhook server
func (wh *Webhook) Run() error {
	err := wh.config.loadTLSConfigFiles()
	if err != nil {
		openlogging.Error("load tls configure files failed: " + err.Error())
		return err
	}

	http.HandleFunc("/v1/mesher/inject", wh.injectHandler)

	if err := wh.server.ListenAndServeTLS("", ""); err != nil {
		openlogging.Error("Webhook ListenAndServeTLS failed: " + err.Error())
		return err
	}
	return nil
}

func (wh *Webhook) injectHandler(rw http.ResponseWriter, r *http.Request) {
	// Webhooks are sent a POST request, with Content-Type: application/json
	if r.Method != http.MethodPost {
		openlogging.Error("request method not allowed: " + r.Method)
		http.Error(rw, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "invalid Content-Type, want `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		openlogging.Error("request body is not found")
		http.Error(rw, "request body is not found", http.StatusBadRequest)
		return
	}

	var response *v1beta1.AdmissionResponse

	reviewReq := &v1beta1.AdmissionReview{}
	err = json.Unmarshal(body, reviewReq)
	if err != nil {
		openlogging.Error(fmt.Sprintf("json unmarshal review request: error = %v, bytes = %s", err, string(body)))
		http.Error(rw, fmt.Sprintf("json marshal review response failed: %s", err), http.StatusBadRequest)
		return
	} else {
		response = wh.inject(reviewReq)
	}

	reviewResp := v1beta1.AdmissionReview{}
	if response != nil {
		reviewResp.Response = response
		if reviewReq.Request != nil {
			reviewResp.Response.UID = reviewReq.Request.UID
		}
	}

	resp, err := json.Marshal(reviewResp)
	if err != nil {
		openlogging.Error("json marshal review response failed: " + err.Error())
		http.Error(rw, fmt.Sprintf("json marshal review response failed: %s", err), http.StatusBadRequest)
		return
	}

	_, err = rw.Write(resp)
	if err != nil {
		openlogging.Error("write response failed: " + err.Error())
		http.Error(rw, fmt.Sprintf("write response failed: %s", err), http.StatusInternalServerError)
	}
}
