#!/usr/bin/env bash

echo -e "go generate...\c"

GOOGLE_PROTOBUF_INCLUDE=../_pb_include/v3.5.1
protoc -I=../ -I=../inside_pkg -I="$GOOGLE_PROTOBUF_INCLUDE" ../req.proto inside_pkg/item.proto --go_out=plugins=grpc:pb_test

echo "OK"