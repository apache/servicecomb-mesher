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
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/var/lib/jenkins/workspace/Mesher
rm -rf /var/lib/jenkins/workspace/Mesher/src
export BUILD_DIR=/var/lib/jenkins/workspace/Mesher
export PROJECT_DIR=$(dirname $BUILD_DIR)

mkdir -p /var/lib/jenkins/workspace/Mesher/src/github.com/apache/servicecomb-mesher

#To checkout to particular commit or tag
if [ $CHECKOUT_VERSION == "latest" ]; then
    echo "using latest code"
else
    git checkout $CHECKOUT_VERSION
fi

cp -r /var/lib/jenkins/workspace/mesher/* /var/lib/jenkins/workspace/Mesher/src/github.com/apache/servicecomb-mesher/

release_dir=$PROJECT_DIR/release
repo="github.com"
project="go-mesh"

if [ -d $release_dir ]; then
    rm -rf $release_dir
fi
mkdir -p $release_dir

cd $BUILD_DIR/src/$repo/$project/mesher

GO111MODULE=on go mod download
GO111MODULE=on go mod vendor
go build -o mesher -a

if [ $VERSION != "latest" ]; then
    cd $PROJECT_DIR/mesher
    git tag -a $TAG_VERSION -m "$TAG_MESSAGE"
    git push origin $TAG_VERSION
fi

export WORK_DIR=$BUILD_DIR/src/$repo/$project/mesher

cp -r $WORK_DIR/licenses $release_dir
#cp $WORK_DIR/NOTICE $release_dir
cp -r $WORK_DIR/conf $release_dir
cp $WORK_DIR/start.sh  $release_dir
cp $WORK_DIR/mesher  $release_dir
if [ ! "$GIT_COMMIT" ];then
   export GIT_COMMIT=`git rev-parse HEAD`
fi

export GIT_COMMIT=`echo $GIT_COMMIT | cut -b 1-7`
RELEASE_VERSION=${releaseVersion:-"latest"}
BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S +%z")

cat << EOF > $release_dir/VERSION
---
version:    $RELEASE_VERSION
commit:     $GIT_COMMIT
built:      $BUILD_TIME
Go-Chassis: $go_sdk_version
EOF


cd $release_dir

chmod +x start.sh mesher

pkg_name="mesher-$VERSION-linux-amd64.tar.gz"

tar zcvf $pkg_name licenses conf mesher VERSION
if [ $JOB_NAME != "" ]; then
    cp $release_dir/$pkg_name /var/lib/jenkins/mesher-release
fi

if [ $VERSION != "latest" ]; then
    date=$(date +%Y-%m-%d)
    DIR_NAME="mesher-release-$date"
    mkdir -p /var/lib/jenkins/userContent/mesher-release/$DIR_NAME
    cp $release_dir/$pkg_name /var/lib/jenkins/userContent/mesher-release/$DIR_NAME
fi
tar zcvf $WORK_DIR/mesher.tar.gz licenses conf start.sh mesher VERSION
exit 0
