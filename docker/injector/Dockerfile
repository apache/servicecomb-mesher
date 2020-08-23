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

COPY . /go/src/github.com/apache/servicecomb-mesher/ 
WORKDIR /go/src/github.com/apache/servicecomb-mesher/
ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on
RUN go build -a github.com/apache/servicecomb-mesher/cmd/injector

FROM frolvlad/alpine-glibc:latest
WORKDIR /home
COPY --from=builder /go/src/github.com/apache/servicecomb-mesher/injector /home
ENTRYPOINT ["./injector", "start"]
