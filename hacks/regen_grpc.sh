#!/bin/bash

PRJ_DIR=$(dirname $(readlink -f $(dirname $0)))
cd $PRJ_DIR

protoc --go_out=plugins=grpc:. ./service/datamesh.proto

