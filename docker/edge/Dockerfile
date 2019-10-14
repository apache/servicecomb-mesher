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

FROM golang:1.12.10 as builder

COPY . /servicecomb-mesher/
WORKDIR /servicecomb-mesher/
ENV GOPROXY=https://goproxy.io
RUN go build -a github.com/apache/servicecomb-mesher/cmd/mesher


FROM frolvlad/alpine-glibc:latest
RUN mkdir -p /opt/mesher && \
    mkdir -p /etc/mesher/conf && \
    mkdir -p /etc/ssl/mesher/
# To upload schemas using env enable SCHEMA_ROOT as environment variable using dockerfile or pass while running container
#ENV SCHEMA_ROOT=/etc/chassis-go/schemas umcomment in future

ENV CHASSIS_HOME=/opt/mesher/

COPY --from=builder /servicecomb-mesher/mesher $CHASSIS_HOME
COPY docker/edge/microservice.yaml docker/edge/chassis.yaml docker/edge/lager.yaml $CHASSIS_HOME/conf/
COPY docker/edge/mesher.yaml /etc/mesher/conf/
COPY docker/edge/start.sh  $CHASSIS_HOME
WORKDIR $CHASSIS_HOME
ENTRYPOINT ["sh", "/opt/mesher/start.sh"]
