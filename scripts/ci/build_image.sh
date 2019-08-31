#!/usr/bin/env bash

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

set -e
set -x

cd /var/lib/jenkins/workspace/Mesher/src/github.com/apache/servicecomb-mesher/

repo="github.com"
project="go-mesh"
export BUILD_DIR=/var/lib/jenkins/workspace/Mesher
export WORK_DIR=$BUILD_DIR/src/$repo/$project/mesher
cd $WORK_DIR

docker build -t gomesh/mesher:$VERSION .

cp /var/lib/jenkins/workspace/docker_login.sh .
bash docker_login.sh &> /dev/null

if [ $PUSH_WITH_LATEST_TAG == "YES" ]; then
    docker build -t gomesh/mesher:latest .
    docker push gomesh/mesher:latest
fi

docker push gomesh/mesher:$VERSION

exit 0
