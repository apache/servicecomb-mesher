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
COPY docker/edge/microservice.yaml $CHASSIS_HOME/conf/
COPY docker/edge/chassis.yaml $CHASSIS_HOME/conf/
COPY docker/edge/mesher.yaml /etc/mesher/conf/
COPY docker/edge/lager.yaml $CHASSIS_HOME/conf/
COPY docker/edge/start.sh  $CHASSIS_HOME
WORKDIR $CHASSIS_HOME
ENTRYPOINT ["sh", "/opt/mesher/start.sh"]
