FROM frolvlad/alpine-glibc:latest
RUN mkdir -p /opt/mesher && \
    mkdir -p /etc/mesher/conf && \
    mkdir -p /etc/ssl/meshercert/ && \
    mkdir -p /etc/chassis-go/schemas/
# To upload schemas using env enable SCHEMA_ROOT as environment variable using dockerfile or pass while running container
#ENV SCHEMA_ROOT=/etc/chassis-go/schemas umcomment in future
ADD mesher.tar.gz /opt/mesher
COPY conf/mesher.yaml /etc/mesher/conf
ENV CHASSIS_HOME=/opt/mesher/
WORKDIR $CHASSIS_HOME
ENTRYPOINT ["sh", "/opt/mesher/start.sh"]
