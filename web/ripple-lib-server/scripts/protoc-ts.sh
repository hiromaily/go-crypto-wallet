#!/bin/bash
# https://caddi.tech/archives/455

PLUGIN_TS=./node_modules/.bin/protoc-gen-ts
PLUGIN_GRPC=./node_modules/.bin/grpc_tools_node_protoc_plugin
DIST_DIR=./src/pb
PROTO_PATH=../../data/proto

rm -rf ${DIST_DIR}
mkdir ${DIST_DIR}

protoc \
--js_out=import_style=commonjs,binary:"${DIST_DIR}"/ \
--ts_out=import_style=commonjs,binary:"${DIST_DIR}"/ \
--grpc_out="${DIST_DIR}"/ \
--plugin=protoc-gen-grpc="${PLUGIN_GRPC}" \
--plugin=protoc-gen-ts="${PLUGIN_TS}" \
--proto_path=${PROTO_PATH}/rippleapi \
--proto_path=${PROTO_PATH}/third_party \
${PROTO_PATH}/rippleapi/*.proto \
${PROTO_PATH}/third_party/gogo/protobuf/gogoproto/*.proto


# "`yarn bin`"/grpc_tools_node_protoc \
#   --plugin=protoc-gen-grpc="`yarn bin`"/grpc_tools_node_protoc_plugin \
#   --js_out=import_style=commonjs,binary:./generated \
#   --grpc_out=./generated \
#   -I./ ./proto/proto/rippleapi/rippleapi.proto

# "`yarn bin`"/grpc_tools_node_protoc \
#   --plugin=protoc-gen-ts="`yarn bin`"/protoc-gen-ts \
#   --ts_out=./generated \
#   -I./ ./proto/proto/rippleapi/rippleapi.proto
