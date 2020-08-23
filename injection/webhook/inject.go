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
	"encoding/json"
	"fmt"

	"github.com/go-mesh/openlogging"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (wh *Webhook) inject(review *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	reqBytes := review.Request.Object.Raw
	pod := &corev1.Pod{}
	if err := json.Unmarshal(reqBytes, pod); err != nil {
		openlogging.Error(fmt.Sprintf("json marshal pod: err = %s, bytes = %s", err.Error(), string(reqBytes)))
		return responseFromErrorMessage(err)
	}

	tmplSpec, err := wh.templater.PodSpecFromTemplate(pod)
	if err != nil {
		openlogging.Error("get pod spec from template failed, error = " + err.Error())
		return responseFromErrorMessage(err)
	}

	tmplSpec.Containers = addK8sServiceAccount(pod.Spec.Containers, tmplSpec.Containers)

	pathBytes, err := createJSONPatch(&pod.Spec, tmplSpec)
	if err != nil {
		openlogging.Error("create patch failed: " + err.Error())
		return responseFromErrorMessage(err)
	}

	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   pathBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

func addK8sServiceAccount(dst, src []corev1.Container) []corev1.Container {
	var serviceAccount *corev1.VolumeMount
	for _, dc := range dst {
		serviceAccount = getServiceAccount(dc.VolumeMounts)
		if serviceAccount != nil {
			break
		}
	}

	for index, sc := range src {
		if matchVolumeMountByName(sc.VolumeMounts, serviceAccount.Name) {
			continue
		}
		src[index].VolumeMounts = append(src[index].VolumeMounts, *serviceAccount)
	}
	return src
}

func getServiceAccount(volumeMounts []corev1.VolumeMount) *corev1.VolumeMount {
	for _, vMount := range volumeMounts {
		if vMount.MountPath == "/var/run/secrets/kubernetes.io/serviceaccount" {
			return &vMount
		}
	}
	return nil
}

func matchVolumeMountByName(volumeMounts []corev1.VolumeMount, name string) bool {
	for _, vMount := range volumeMounts {
		if vMount.Name == name {
			return true
		}
	}
	return false
}

func responseFromErrorMessage(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{Result: &metav1.Status{Message: err.Error()}}
}
