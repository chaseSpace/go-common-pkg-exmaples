#!/bin/bash

# protoc v3.5.1
# protoc-gen-go v1.5.3
# 注意这里不需要指定输出位置，因为全路径已经在所有proto文件中定义好
protoc -I=./pb -I=./pb/protoc-v3.5.1 --go_out=plugins=grpc:. $(find ./pb/protocol/ -name '*.proto')

<<comment
除此之外，官方支持如下方式为特定文件指定 import path,  --go_opt=M${PROTO_FILE}=${GO_IMPORT_PATH}
protoc --proto_path=src \
  --go_opt=Mprotos/buzz.proto=example.com/project/protos/fizz \
  --go_opt=Mprotos/bar.proto=example.com/project/protos/foo \
  protos/buzz.proto protos/bar.proto
comment
