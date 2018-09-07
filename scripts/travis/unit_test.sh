#!/bin/sh
set -e

echo "mode: atomic" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    echo $d
    cd $GOPATH/src/$d
    if [ $(ls | grep _test.go | wc -l) -gt 0 ]; then
        go test -cover -covermode atomic -coverprofile coverage.out
        if [ -f coverage.out ]; then
            sed '1d;$d' coverage.out >> $GOPATH/src/github.com/go-chassis/mesher/coverage.txt
        fi
    fi
done