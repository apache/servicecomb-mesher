#!/bin/sh

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

gosec -fmt=json -out=results.json ./...
cat results.json
issueCount=$(tail -7 results.json | grep '"found":' | sed 's/[[:space:]]//g' | awk -F":" '{print $2}')
rm -rf results.json
if (($? == 0 && 35 > $issueCount)); then
	echo "No GoSecure warnings found"
	exit 0
else
	echo "GoSecure Warnings found"
	exit 1
fi
