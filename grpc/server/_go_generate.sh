#!/usr/bin/env bash

echo -e "go generate...\c"

cd ../
protoc -I=. -I=inside_pkg -I="$GOOGLE_PROTOBUF_INCLUDE" req.proto inside_pkg/item.proto --go_out=plugins=grpc:pb_test

echo "OK"