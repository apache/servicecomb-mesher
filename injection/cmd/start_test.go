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

package cmd

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func getCliContext(args []string) (*cli.Context, error) {
	app := cli.NewApp()
	app.HideVersion = true
	app.Name = "injector"
	app.Usage = "Kubernetes webhook for automatic ServiceComb mesher injection."
	app.Commands = []cli.Command{
		GetCmdStart(),
		GetCmdVersion(),
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)
	err := flagSet.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	return cli.NewContext(app, flagSet, nil), nil
}

func TestStart(t *testing.T) {
	defer func() {
		if info := recover(); info != nil {
			t.Log("has panic")
		}
	}()
	ctx, err := getCliContext([]string{"test", "start"})
	assert.Nil(t, err)
	start(ctx)
}
