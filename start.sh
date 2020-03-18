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
CIRCUIT_BREAKER_YAML="circuit_breaker.yaml"
LOAD_BALANCING_YAML="load_balancing.yaml"
MONITORING_YAML="monitoring.yaml"
LAGER_YAML="lager.yaml"
RATE_LIMITING_YAML="rate_limiting.yaml"
TLS_YAML="tls.yaml"
AUTH_YAML="auth.yaml"
TRACING_YAML="tracing.yaml"
ROUTER_YAML="router.yaml"

TMP_DIR="/tmp"

check_config_files(){
    # configs can be mounted, maybe config map
    if [ -f "$TMP_DIR/$MESHER_YAML" ]; then
        echo "$MESHER_YAML is customed"
        cp -f $TMP_DIR/$MESHER_YAML $ETC_CONF_DIR/$MESHER_YAML
    fi
    copy_tmp2mesher $CHASSIS_YAML
    copy_tmp2mesher $MICROSERVICE_YAML
    copy_tmp2mesher $CIRCUIT_BREAKER_YAML
    copy_tmp2mesher $LOAD_BALANCING_YAML
    copy_tmp2mesher $MONITORING_YAML
    copy_tmp2mesher $LAGER_YAML
    copy_tmp2mesher $RATE_LIMITING_YAML
    copy_tmp2mesher $TLS_YAML
    copy_tmp2mesher $AUTH_YAML
    copy_tmp2mesher $TRACING_YAML
    copy_tmp2mesher $ROUTER_YAML
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

# configure handler chain
if [ ! -z "$OUTGOING_HANDLER" ]; then
    sed -i s/"outgoing:.*"/"outgoing: $OUTGOING_HANDLER"/g $MESHER_CONF_DIR/$CHASSIS_YAML
fi
if [ ! -z "$INCOMING_HANDLER" ]; then
    sed -i s/"incoming:.*"/"incoming: $INCOMING_HANDLER"/g $MESHER_CONF_DIR/$CHASSIS_YAML
fi

# configure refreshInterval
if [ ! -z "$REFRESH_INTERVAL" ]; then
    sed -i s/"refreshInterval:.*"/"refreshInterval: $REFRESH_INTERVAL"/g $MESHER_CONF_DIR/$CHASSIS_YAML
fi

#add instance_properties if any provided
if [ ! -z "$MD" ]; then
  echo "  instance_properties:" >> $MESHER_CONF_DIR/$MICROSERVICE_YAML
  while : 
  do
    KeyArray=$(echo $MD | cut -d '|' -f1)
    RestArray=$(echo $MD | cut -d '|' -f2-)
    Writer="-"
    Key=$(echo $KeyArray | cut -d '=' -f1)
    Value=$(echo $KeyArray | cut -d '=' -f2-)
    Writer="    $Writer$Key"
    Writer="$Writer: $Value"
    echo $Writer| sed 's/-/    /'
    echo $Writer | sed 's/-/    /'>> $MESHER_CONF_DIR/$MICROSERVICE_YAML
    if [ "$KeyArray" = "$RestArray" ]; then
      break
    fi
    MD=$RestArray
  done
fi

# ENABLE_PROXY_TLS decide whether mesher is https proxy or http proxy
if [[ $TLS_ENABLE && $TLS_ENABLE == true ]]; then
    sed -i '/ssl:/a \ \ mesher.Provider.cipherPlugin: default \n \ mesher.Provider.verifyPeer: false \n \ mesher.Provider.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 \n \ mesher.Provider.protocol: TLSv1.2 \n \ mesher.Provider.caFile: \n \ mesher.Provider.certFile: /etc/ssl/meshercert/kubecfg.crt \n \ mesher.Provider.keyFile: /etc/ssl/meshercert/kubecfg.key \n \ mesher.Provider.certPwdFile: \n' $MESHER_CONF_DIR/$TLS_YAML
fi

exec $MESHER_DIR/mesher --config $ETC_CONF_DIR/$MESHER_YAML
