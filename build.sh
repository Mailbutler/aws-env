#!/bin/bash

BUILD_DIR=bin
NAME=aws-env

mkdir -p $BUILD_DIR

go mod tidy

for GOOS in darwin; do
    for GOARCH in amd64 arm64; do
        rm -f $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "compiling $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "upx-ing $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        upx --best --lzma $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done

for GOOS in linux; do
    for GOARCH in 386 amd64 arm64; do
        rm -f $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "compiling $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "upx-ing $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        upx --best --lzma $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done

for GOOS in windows; do
    for GOARCH in 386 amd64; do
        rm -f $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "compiling $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
        echo "upx-ing $BUILD_DIR/$NAME-$GOOS-$GOARCH"
        upx --best --lzma $BUILD_DIR/$NAME-$GOOS-$GOARCH
    done
done
