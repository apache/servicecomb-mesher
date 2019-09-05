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
export BUILD_DIR=$(cd "$(dirname "$0")"; pwd)
export PROJECT_DIR=$(dirname ${BUILD_DIR})

if [ -z "${GOPATH}" ]; then
 echo "missing GOPATH env, can not build"
 exit 1
fi
echo "GOPATH is "${GOPATH}


#To checkout to particular commit or tag
if [ "$VERSION" == "" ]; then
    echo "using latest code"
    VERSION="latest"
fi

release_dir=$PROJECT_DIR/release
mkdir -p $release_dir
cd $PROJECT_DIR
GO111MODULE=on go mod download
GO111MODULE=on go mod vendor
go build -a github.com/apache/servicecomb-mesher/cmd/mesher

cp -r $PROJECT_DIR/licenses $release_dir
cp -r $PROJECT_DIR/licenses/LICENSE $release_dir
cp -r $PROJECT_DIR/licenses/NOTICE $release_dir
cp -r $PROJECT_DIR/conf $release_dir
cp $PROJECT_DIR/start.sh  $release_dir
cp $PROJECT_DIR/mesher  $release_dir
if [ ! "$GIT_COMMIT" ];then
   export GIT_COMMIT=`git rev-parse HEAD`
fi

export GIT_COMMIT=`echo $GIT_COMMIT | cut -b 1-7`
BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S +%z")

cat << EOF > $release_dir/VERSION
---
version:    $VERSION
commit:     $GIT_COMMIT
built:      $BUILD_TIME
EOF


cd $release_dir

chmod +x start.sh mesher

component="apache-servicecomb-mesher"
x86_pkg_name="$component-$VERSION-linux-amd64.tar.gz"
arm_pkg_name="$component-$VERSION-linux-arm64.tar.gz"

#x86 release
tar zcvf $x86_pkg_name licenses conf mesher VERSION LICENSE NOTICE
tar zcvf mesher.tar.gz licenses conf mesher VERSION LICENSE NOTICE start.sh # for docker image


echo "building docker..."
cd ${release_dir}
cp ${PROJECT_DIR}/build/docker/proxy/Dockerfile ./
sudo docker build -t servicecomb/mesher-sidecar:${VERSION} .

# arm release
GOARCH=arm64 go build -a github.com/apache/servicecomb-mesher/cmd/mesher
tar zcvf $arm_pkg_name licenses conf mesher VERSION LICENSE NOTICE