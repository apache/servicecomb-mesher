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

set -e

MESHER_DIR=$(cd $(dirname $0); pwd)
MESHER_CONF_DIR="$MESHER_DIR/conf"

MESHER_YAML="mesher.yaml"
ETC_CONF_DIR="/etc/mesher/conf"

CHASSIS_YAML="chassis.yaml"
MICROSERVICE_YAML="microservice.yaml"
MONITORING_YAML="monitoring.yaml"
LAGER_YAML="lager.yaml"
TLS_YAML="tls.yaml"
AUTH_YAML="auth.yaml"
TRACING_YAML="tracing.yaml"

TMP_DIR="/tmp"

check_config_files(){
    # configs can be mounted, maybe config map
    if [ -f "$TMP_DIR/$MESHER_YAML" ]; then
        echo "$MESHER_YAML is customed"
        cp -f $TMP_DIR/$MESHER_YAML $ETC_CONF_DIR/$MESHER_YAML
    fi
    copy_tmp2mesher $CHASSIS_YAML
    copy_tmp2mesher $MONITORING_YAML
    copy_tmp2mesher $LAGER_YAML
    copy_tmp2mesher $TLS_YAML
    copy_tmp2mesher $AUTH_YAML
    copy_tmp2mesher $TRACING_YAML
}
copy_tmp2mesher(){
    if [ -f $TMP_DIR/$1 ]; then
        echo "$1 is customed"
        cp -f $TMP_DIR/$1 $MESHER_CONF_DIR/$1
    fi
}

check_config_files

net_name=$(ip -o -4 route show to default | awk '{print $5}')
listen_addr=$(ifconfig $net_name | grep -E 'inet\W' | grep -o -E [0-9]+.[0-9]+.[0-9]+.[0-9]+ | head -n 1)

# replace ip addr
sed -i s/"listenAddress:\s\{1,\}[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/"listenAddress: $listen_addr"/g $MESHER_CONF_DIR/$CHASSIS_YAML

exec $MESHER_DIR/mesher --config $ETC_CONF_DIR/$MESHER_YAML --mode edge
