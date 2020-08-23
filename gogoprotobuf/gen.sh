#!/usr/bin/env bash

gogo_google=/c/Users/LEI/go/pkg/mod/github.com/gogo/protobuf@v1.3.1/protobuf

protoc -I=. --gogo_out=plugins=grpc:. req.proto