#!/usr/bin/env bash
set -e
set -x

cd /var/lib/jenkins/workspace/Mesher/src/github.com/go-chassis/mesher/

repo="github.com"
project="go-chassis"
export BUILD_DIR=/var/lib/jenkins/workspace/Mesher
export WORK_DIR=$BUILD_DIR/src/$repo/$project/mesher
cd $WORK_DIR

docker build -t gochassis/mesher:$VERSION .

cp /var/lib/jenkins/workspace/docker_login.sh .
bash docker_login.sh &> /dev/null

if [ $PUSH_WITH_LATEST_TAG == "YES" ]; then
    docker build -t gochassis/mesher:latest .
    docker push gochassis/mesher:latest
fi

docker push gochassis/mesher:$VERSION

exit 0
