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
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

type operation string

const (
	OperationAdd     operation = "add"
	OperationReplace operation = "replace"
	OperationRemove  operation = "remove"
)

type jsonPatch struct {
	Op    operation   `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func addPatch(src interface{}, basePath string) (patch []jsonPatch) {
	return append(patch, jsonPatch{Op: OperationAdd, Path: basePath, Value: src})
}

func removePatch(basePath string) (patch []jsonPatch) {
	return append(patch, jsonPatch{Op: OperationRemove, Path: basePath})
}

func createJSONPatch(dst *corev1.PodSpec, src *corev1.PodSpec) ([]byte, error) {
	patch := make([]jsonPatch, 0, 10)
	patch = append(patch, containersToPatch(dst.InitContainers, src.InitContainers, "Name", "/spec/initContainers")...)
	patch = append(patch, containersToPatch(dst.Containers, src.Containers, "Name", "/spec/containers")...)
	patch = append(patch, sliceToPatch(dst.Volumes, src.Volumes, "Name", "/spec/volumes")...)
	patch = append(patch, sliceToPatch(dst.ImagePullSecrets, src.ImagePullSecrets, "Name", "/spec/imagePullSecrets")...)

	if src.DNSConfig != nil {
		patch = append(patch, addPatch(src.DNSConfig, "/spec/dnsConfig")...)
	}

	if dst.SecurityContext != nil {
		patch = append(patch, addPatch(dst.SecurityContext, "/spec/securityContext")...)
	}

	return json.Marshal(patch)
}

func containersToPatch(dst []corev1.Container, src []corev1.Container, field string, basePath string) (patch []jsonPatch) {
	patch = append(patch, sliceToPatch(dst, src, field, basePath)...)
	for i, dic := range dst {
		for _, sic := range src {
			if dic.Name == sic.Name {
				patch = append(patch, sliceToPatch(dic.Env, sic.Env, "Name", fmt.Sprintf("%v/%v/%v", basePath, i, "env"))...)
				patch = append(patch, sliceToPatch(dic.VolumeMounts, sic.VolumeMounts, "Name", fmt.Sprintf("%v/%v/%v", basePath, i, "volumeMounts"))...)
			}
		}
	}
	return
}

func sliceToPatch(dst interface{}, src interface{}, field string, basePath string) (patch []jsonPatch) {
	dVal := reflect.ValueOf(dst)
	sVal := reflect.ValueOf(src)
	if dVal.Kind() != sVal.Kind() || dVal.IsNil() && sVal.IsNil() {
		return
	}
	return addSlicePatch(dVal, sVal, field, basePath)
}

func removeSlicePatch(dst reflect.Value, src reflect.Value, field string, basePath string) (patch []jsonPatch) {
	for i := dst.Len() - 1; i >= 0; i-- {
		dv := dst.Index(i)
		if matchSliceByField(src, dv, field) {
			continue
		}
		patch = append(patch, removePatch(fmt.Sprintf("%v/%v", basePath, i))...)
	}
	return patch
}

func addSlicePatch(dst reflect.Value, src reflect.Value, field string, basePath string) (patch []jsonPatch) {
	first := dst.Len() == 0
	for i := 0; i < src.Len(); i++ {
		sv := src.Index(i)
		if matchSliceByField(dst, sv, field) {
			continue
		}

		value := sv.Interface()
		path := basePath
		if first {
			first = false
			value = []interface{}{value}
		} else {
			path += "/-"
		}
		patch = append(patch, addPatch(value, path)...)
	}
	return patch
}

func matchSliceByField(set reflect.Value, match reflect.Value, field string) bool {
	for i := 0; i < set.Len(); i++ {
		item := set.Index(i)
		if item.FieldByName(field).Interface() == match.FieldByName(field).Interface() {
			return true
		}
	}
	return false
}
