#!/bin/bash

# protoc v3.5.1
# protoc-gen-go v1.5.3
# 注意这里不需要指定输出位置，因为全路径已经在所有proto文件中定义好
protoc -I=./pb -I=./pb/protoc-v3.5.1 --go_out=plugins=grpc:. $(find ./pb/protocol/ -name '*.proto')