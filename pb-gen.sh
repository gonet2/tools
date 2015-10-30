#!/bin/sh
protoc  ./*.proto --go_out=plugins=grpc:src/proto
