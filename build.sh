#!/bin/bash

BUILD_DIR=bin
NAME=aws-env

mkdir -p $BUILD_DIR

for GOOS in darwin linux windows; do
    for GOARCH in amd64 arm64; do
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done
