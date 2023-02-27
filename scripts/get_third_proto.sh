#!/bin/bash

mkdir -p ./data/proto/third_party/google/protobuf

mkdir tmp && cd tmp
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf
#git checkout -b target 02f4f392c023d8845381d2db9fb6d7eb167908ff

cd src/google/protobuf
cp empty.proto ../../../../../data/proto/third_party/google/protobuf/
cp any.proto ../../../../../data/proto/third_party/google/protobuf/

cd ../../../../../
rm -rf tmp
