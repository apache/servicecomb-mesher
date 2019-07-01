#!/usr/bin/env bash
set -e
set -x
export BUILD_DIR=$(cd "$(dirname "$0")"; pwd)
export PROJECT_DIR=$(dirname ${BUILD_DIR})

if [ -z "${GOPATH}" ]; then
 echo "missing GOPATH env, can not build"
 exit 1
fi
echo "GOPATH is "${GOPATH}


#To checkout to particular commit or tag
if [ "$VERSION" == "" ]; then
    echo "using latest code"
    VERSION="latest"
else
    git checkout $VERSION
fi

release_dir=$PROJECT_DIR/release
mkdir -p $release_dir
cd $PROJECT_DIR
GO111MODULE=on go mod download
GO111MODULE=on go mod vendor
go build -o mesher -a

cp -r $PROJECT_DIR/licenses $release_dir
cp -r $PROJECT_DIR/conf $release_dir
cp $PROJECT_DIR/start.sh  $release_dir
cp $PROJECT_DIR/mesher  $release_dir
if [ ! "$GIT_COMMIT" ];then
   export GIT_COMMIT=`git rev-parse HEAD`
fi

export GIT_COMMIT=`echo $GIT_COMMIT | cut -b 1-7`
BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S +%z")

cat << EOF > $release_dir/VERSION
---
version:    $VERSION
commit:     $GIT_COMMIT
built:      $BUILD_TIME
EOF


cd $release_dir

chmod +x start.sh mesher

pkg_name="mesher-$VERSION-linux-amd64.tar.gz"

tar zcvf $pkg_name licenses conf mesher VERSION
tar zcvf mesher.tar.gz licenses conf mesher VERSION start.sh





echo "building docker..."
cd ${release_dir}
cp ${PROJECT_DIR}/build/docker/proxy/Dockerfile ./
sudo docker build -t servicecomb/mesher:${VERSION} .