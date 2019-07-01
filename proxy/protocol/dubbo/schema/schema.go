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

package schema

import (
	"fmt"
	"strings"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
)

const (
	InQuery = 0
	InBody  = 1
)

//InterfaceSchema is a struct
type InterfaceSchema struct {
	SvcName      string
	JavaClsName  string
	BasePath     string
	Version      string
	AppId        string
	ServiceId    string
	methordArray []DefMethod
}

//DefMethod is a struct
type DefMethod struct {
	ownerSvc string
	Path     string //  样例: /test/StringArray 需要是包含basepath的URL
	OperaID  string
	Paras    []MethParam             //key 参数名称
	Verb     string                  // get post ...
	Responds map[string]*MethRespond //key 为返回码 200 ,404...
}

//GetRspSchema is a method to get response schema
func (this *DefMethod) GetRspSchema(status int) *MethRespond {
	strStatus := fmt.Sprintf("%d", status)
	if _, ok := this.Responds[strStatus]; ok {
		return this.Responds[strStatus]
	} else {
		return nil
	}
}

//GetParamNameAndWhere is a method to get parameter name
func (this *DefMethod) GetParamNameAndWhere(indx int) (string, int) {
	for _, v := range this.Paras {
		if v.Indx == indx {
			//return v.name, v.where
			if strings.EqualFold(v.Where, "query") {
				return v.Name, InQuery
			} else {
				return v.Name, InBody
			}
		}
	}
	return "", InQuery
}

//GetParamSchema is a method to get parameter schema
func (this *DefMethod) GetParamSchema(indx int) *MethParam {
	for _, v := range this.Paras {
		if v.Indx == indx {
			return &v
		}
	}
	return nil
}

//MethRespond is a struct
type MethRespond struct {
	Status string // 200 404...
	DType  string
	ObjRef DefType
}

//MethParam is a struct
type MethParam struct {
	Name            string            //参数名称
	Dtype           string            //参数类型
	AdditionalProps map[string]string //附加参数，如果是Dtype是map，则在这里定义value类型
	Items           map[string]string //当Dtype为array时使用
	Required        bool              //是否必需
	Where           string            //存储位置  query or body
	Indx            int
	ObjRef          DefType
}

//DefType 契约definitions字段里定义的类型定义
type DefType struct {
	typeName   string
	dType      string
	JvmClsName string
	fileds     map[string]DefField
}

//DefField is a struct
type DefField struct {
	dType   string
	formate string
}

//GetDefTypeFromDef is a function to get defintion type
func GetDefTypeFromDef(defs map[string]registry.Definition, ref string) DefType {
	def := DefType{}
	names := strings.Split(ref, "/")
	def.typeName = names[len(names)-1]
	if sDef, ok := defs[def.typeName]; ok {
		def.dType = sDef.Types
		def.JvmClsName = sDef.XJavaClass
	}
	return def
}

//CovertSwaggerMethordToLocalMethord is a function to convert swagger method to local method
func CovertSwaggerMethordToLocalMethord(schema *registry.SchemaContent, src *registry.MethodInfo, dst *DefMethod) {
	dst.OperaID = src.OperationID
	tmpParas := make([]MethParam, len(src.Parameters))
	i := 0
	for _, para := range src.Parameters {
		var defPara MethParam
		defPara.Name = para.Name
		if para.Type == "" {
			if para.Schema.Type != "" {
				defPara.Dtype = para.Schema.Type
			} else {
				defPara.Dtype = "object"
				if para.Schema.Reference != "" {
					defPara.ObjRef = GetDefTypeFromDef(schema.Definition, para.Schema.Reference)
				}
			}
		} else {
			defPara.Dtype = para.Type
		}
		defPara.Required = para.Required
		defPara.Where = para.In
		defPara.Indx = i
		tmpParas[i] = defPara
		i++
	}
	dst.Paras = tmpParas
	tmpRsps := make(map[string]*MethRespond)
	for key, rsp := range src.Response {
		var defRsp MethRespond
		if dtype, ok := rsp.Schema["type"]; ok {
			defRsp.DType = dtype
		} else {
			defRsp.DType = "object"
			if dRef, ok := rsp.Schema["$ref"]; ok {
				defRsp.ObjRef = GetDefTypeFromDef(schema.Definition, dRef)
			}
		}
		defRsp.Status = key
		tmpRsps[key] = &defRsp
	}
	dst.Responds = tmpRsps
}

//GetSvcByInterface is a function to get service by interface name
func GetSvcByInterface(interfaceName string) *registry.MicroService {
	value, ok := svcToInterfaceCache.Get(interfaceName)
	if !ok || value == nil {
		lager.Logger.Infof("Get svc from remote, interface: %s", interfaceName)
		svc := registry.DefaultContractDiscoveryService.GetMicroServicesByInterface(interfaceName)
		if svc != nil {
			svcKey := strings.Join([]string{svc[0].ServiceName, svc[0].Version, svc[0].AppID}, "/")
			lager.Logger.Infof("Cached svc [%s] for interface %s", svcKey, interfaceName)
			svcToInterfaceCache.Set(interfaceName, svc[0], 0)
			refresher.Add(newInterfaceJob(interfaceName))
		} else {
			return nil
		}
		value, ok = svcToInterfaceCache.Get(interfaceName)
	}
	if value != nil {
		if service, ok2 := value.(*registry.MicroService); ok2 {
			return service
		}
	}
	return nil
}

//GetSupportProto is a function to get supported protocol
func GetSupportProto(svc *registry.MicroService) string {
	if svc == nil {
		return ""
	}
	proto := "dubbo"
	value, ok := protoCache.Get(svc.ServiceID)
	if !ok || value == nil {
		lager.Logger.Infof("Get proto from remote, serviceID: %s", svc.ServiceID)
		ins, err := registry.DefaultServiceDiscoveryService.GetMicroServiceInstances(runtime.ServiceID, svc.ServiceID)
		if err != nil {
			return proto
		}
		lager.Logger.Infof("Cached proto for serviceID %s", svc.ServiceID)
		protoCache.Set(svc.ServiceID, protoForService(ins), 0)
		refresher.Add(newProtoJob(svc.ServiceID))

		value, ok = protoCache.Get(svc.ServiceID)
	}

	if value != nil {
		if cached, ok2 := value.(string); ok2 {
			return cached
		}
	}
	return proto
}

//GetSvcNameByInterface is a function to get service name by interface
func GetSvcNameByInterface(interfaceName string) string {
	svc := registry.DefaultContractDiscoveryService.GetMicroServicesByInterface(interfaceName)
	for _, v := range svc {
		return v.ServiceName
	}
	return ""
}

//GetMethodByInterface is a function to get method from interface name
func GetMethodByInterface(interfaceName string, operateID string) *DefMethod {
	var meth DefMethod

	schema := registry.DefaultContractDiscoveryService.GetSchemaContentByInterface(interfaceName)
	for path, pathSchema := range schema.Paths {
		for verb, m := range pathSchema {
			if strings.EqualFold(m.OperationID, operateID) {
				meth.Verb = strings.ToUpper(verb)
				meth.Path = schema.BasePath + path
				meth.OperaID = operateID
				CovertSwaggerMethordToLocalMethord(&schema, &m, &meth)
				return &meth
			}
		}
	}
	return nil
}

//GetSchemaMethodBySvcURL is a function to get schema method from URl
func GetSchemaMethodBySvcURL(svcName string, env string, ver string, app string, verb string, url string) (*registry.SchemaContent, *DefMethod) {
	schemas := registry.DefaultContractDiscoveryService.GetSchemaContentByServiceName(svcName, ver, app, env)
	var curMethrod *DefMethod
	var curSchema *registry.SchemaContent
	for _, v := range schemas {
		curSchema = v
		methord := GetMethodInfoSchemaByURL(v, verb, url)
		if methord != nil {
			if curMethrod == nil {
				curMethrod = methord
			} else {
				if len(curMethrod.Path) < len(methord.Path) {
					curMethrod = methord
				}
			}
		}
	}
	return curSchema, curMethrod
}

//GetMethodInfoSchemaByURL is a function to get method info schema from URl
func GetMethodInfoSchemaByURL(schema *registry.SchemaContent, verb string, url string) *DefMethod {
	var curMax = 0
	basePath := schema.BasePath
	var method *registry.MethodInfo
	var path string
	for key, v := range schema.Paths {
		if strings.HasPrefix(url, basePath+key) {
			if tmp, ok := v[verb]; ok {
				if len(key) > curMax {
					curMax = len(key)
					method = &tmp
					path = key
				}
			}
		}
	}
	if method != nil {
		tmpMeth := &DefMethod{}
		tmpMeth.Path = basePath + path
		tmpMeth.Verb = verb
		CovertSwaggerMethordToLocalMethord(schema, method, tmpMeth)
		return tmpMeth
	} else {
		return nil
	}
}
