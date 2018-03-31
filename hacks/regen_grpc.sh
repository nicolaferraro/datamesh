#!/bin/bash

PRJ_DIR=$(dirname $(readlink -f $(dirname $0)))
cd $PRJ_DIR

protoc --proto_path=. --proto_path=./hacks --go_out=plugins=grpc:. ./service/datamesh.proto

