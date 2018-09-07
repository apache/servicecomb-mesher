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

package version

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"gopkg.in/yaml.v2"
)

//Version is a struct which has attributes for version
type Version struct {
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit" yaml:"commit"`
	Built     string `json:"built" yaml:"built"`
	GoChassis string `json:"Go-Chassis" yaml:"Go-Chassis"`
}

//Constants
const (
	VersionFile    = "VERSION"
	DefaultVersion = "latest"
)

var version *Version

func setVersion() {
	v, err := getVersionSet()
	if err != nil {
		log.Printf("Get version failed, err: %s", err)
		version = &Version{}
		return
	}
	version = v
}

func getVersionSet() (*Version, error) {
	workDir, err := fileutil.GetWorkDir()
	if err != nil {
		return nil, err
	}
	p := filepath.Join(workDir, VersionFile)
	content, err := ioutil.ReadFile(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		log.Printf("%s not found, mesher version unknown", p)
		return &Version{}, nil
	}
	v := &Version{}
	err = yaml.Unmarshal(content, v)
	if err != nil {
		return nil, &os.PathError{Path: p, Err: err}
	}
	return v, nil
}

//Ver returns version
func Ver() *Version {
	return version
}

func init() {
	setVersion()
}
