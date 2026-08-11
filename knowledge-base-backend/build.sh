#!/bin/sh

mkdir -p build
cp -R ./migration ./build
cp -R ./script ./build
cp -R ./config ./build
cp cmd/knowledge_base_backend_config.ini ./build

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/knowledge_base_backend ./cmd
