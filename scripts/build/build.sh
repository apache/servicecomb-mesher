#!/usr/bin/env bash
set -e
set -x
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/var/lib/jenkins/workspace/Mesher
rm -rf /var/lib/jenkins/workspace/Mesher
export BUILD_DIR=/var/lib/jenkins/workspace/Mesher
export PROJECT_DIR=$(dirname $BUILD_DIR)

mkdir -p /var/lib/jenkins/workspace/Mesher/src/github.com/go-chassis/mesher

cp -r /var/lib/jenkins/workspace/mesher/* /var/lib/jenkins/workspace/Mesher/src/github.com/go-chassis/mesher/

release_dir=$PROJECT_DIR/release
repo="github.com"
project="go-chassis"

if [ -d $release_dir ]; then
    rm -rf $release_dir
fi
mkdir -p $release_dir

cd $BUILD_DIR/src/$repo/$project/mesher
#To checkout to particular commit or tag
if [ $CHECKOUT_VERSION == "latest" ]; then
    echo "using latest code"
else
    git checkout $CHECKOUT_VERSION
fi
GO111MODULE=on go mod download
GO111MODULE=on go mod vendor
go build -o mesher -a


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
tar zcvf $WORK_DIR/mesher.tar.gz licenses conf start.sh mesher VERSION
exit 0
