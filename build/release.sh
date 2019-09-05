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

export BUILD_DIR=$(cd "$(dirname "$0")"; pwd)
export PROJECT_DIR=$(dirname ${BUILD_DIR})

component="apache-servicecomb-mesher"
x86_pkg_name="$component-$VERSION-linux-amd64.tar.gz"
arm_pkg_name="$component-$VERSION-linux-arm64.tar.gz"
cd $PROJECT_DIR/release
#asc
gpg --armor --output "${x86_pkg_name}".asc --detach-sig "${x86_pkg_name}"
gpg --armor --output "${arm_pkg_name}".asc --detach-sig "${arm_pkg_name}"
#512
sha512sum "${x86_pkg_name}" > "${x86_pkg_name}".sha512
sha512sum "${arm_pkg_name}" > "${arm_pkg_name}".sha512
#src
wget "https://github.com/apache/servicecomb-mesher/archive/v${VERSION}.tar.gz"

src_name="${component}-${VERSION}-src.tar.gz"
mv "v${VERSION}.tar.gz" "${src_name}"

gpg --armor --output "$src_name.asc" --detach-sig "${src_name}"

sha512sum "${src_name}" > "${src_name}".sha512